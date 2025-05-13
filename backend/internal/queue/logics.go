package queue

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go.uber.org/zap"

	configPb "backend/generated/proto/config"
	"backend/generated/proto/request"
	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/types"
)

type QueueLogics struct {
	logger        *zap.Logger
	cfg           *config.Config
	ctx           context.Context
	configClient  *grpcfactory.Client[configPb.ConfigServiceClient]
	requestClient *grpcfactory.Client[request.RequestServiceClient]
	inputChannel  chan types.InputForm
	workerPool    chan types.InputForm
}

func NewLogics(
	logger *zap.Logger,
	cfg *config.Config,
	configClient *grpcfactory.Client[configPb.ConfigServiceClient],
	requestClient *grpcfactory.Client[request.RequestServiceClient],
) *QueueLogics {
	logger.Info("Creating QueueLogics instance")
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)

	q := &QueueLogics{
		logger:        logger,
		cfg:           cfg,
		ctx:           ctx,
		configClient:  configClient,
		requestClient: requestClient,
		inputChannel:  make(chan types.InputForm),
		workerPool:    make(chan types.InputForm, 100),
	}

	// Start the loop and workers
	go q.loop()
	for i := 0; i < 5; i++ { // 5 workers
		go q.worker(i)
	}

	logger.Info("QueueLogics instance created and workers started")
	return q
}

func (l *QueueLogics) loop() {
	l.logger.Info("Starting loop for message processing")
	for {
		select {
		case msg := <-l.inputChannel:
			l.logger.Debug("Received message", zap.Any("message", msg))
			l.workerPool <- msg // Send task to the worker pool
		case <-l.ctx.Done():
			l.logger.Info("Loop context cancelled, stopping...")
			return
		}
	}
}

func (l *QueueLogics) worker(id int) {
	l.logger.Info("Worker started", zap.Int("worker_id", id))
	for {
		select {
		case msg := <-l.workerPool:
			l.logger.Debug("Worker received task", zap.Int("worker_id", id), zap.Any("message", msg))
			go l.processTask(msg)
		case <-l.ctx.Done():
			l.logger.Info("Worker stopped", zap.Int("worker_id", id))
			return
		}
	}
}

func (l *QueueLogics) processTask(msg types.InputForm) {
	l.logger.Debug("Processing task", zap.Any("message", msg))
	// Fetch config from the config service
	config, err := l.configClient.Service.GetConfig(l.ctx, &configPb.GetConfigRequest{
		ConfigId: msg.ConfigId,
	})
	if err != nil {
		l.logger.Error("Failed to get config", zap.Error(err))
		return
	}
	l.logger.Debug("Fetched config", zap.Int64("config_id", config.Id))

	// Get refresh tokens
	refreshTokenDTF, err := l.requestClient.Service.GetRefreshTokenDTF(l.ctx, &request.GetRefreshTokenRequest{
		Email:    config.Email,
		Password: config.Password,
	})
	if err != nil {
		l.logger.Error("Failed to get refresh token DTF", zap.Error(err))
		return
	}
	l.logger.Debug("Fetched refresh token DTF", zap.String("refresh_token", refreshTokenDTF.RefreshToken))

	refreshTokenVC, err := l.requestClient.Service.GetRefreshTokenVC(l.ctx, &request.GetRefreshTokenRequest{
		Email:    config.Email,
		Password: config.Password,
	})
	if err != nil {
		l.logger.Error("Failed to get refresh token VC", zap.Error(err))
		return
	}
	l.logger.Debug("Fetched refresh token VC", zap.String("refresh_token", refreshTokenVC.RefreshToken))

	// Update refresh tokens
	_, err = l.configClient.Service.UpdateRefreshDTFToken(l.ctx, &configPb.UpdateRefreshTokenDTFRequest{
		RefreshTokenDTF: refreshTokenDTF.RefreshToken,
		Id:              config.Id,
	})
	if err != nil {
		l.logger.Error("Failed to update DTF refresh token", zap.Error(err))
		return
	}
	l.logger.Debug("DTF refresh token updated", zap.Int64("config_id", config.Id))

	_, err = l.configClient.Service.UpdateRefreshVCToken(l.ctx, &configPb.UpdateRefreshTokenVCRequest{
		RefreshTokenVC: refreshTokenVC.RefreshToken,
		Id:             config.Id,
	})
	if err != nil {
		l.logger.Error("Failed to update VC refresh token", zap.Error(err))
		return
	}
	l.logger.Debug("VC refresh token updated", zap.Int64("config_id", config.Id))

	// Process data for generation
	var data map[string][]string
	if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
		l.logger.Error("Failed to parse JSON", zap.Error(err))
		return
	}
	l.logger.Debug("Parsed JSON data", zap.Any("data", data))

	rowCount := len(data)
	for i := 0; i < rowCount; i++ {
		result := config.Prompt
		for key, values := range data {
			if i < len(values) {
				placeholder := "{" + key + "}"
				result = strings.ReplaceAll(result, placeholder, values[i])
			}
		}

		// Generate text
		generatedText, err := l.requestClient.Service.GenerateText(l.ctx, &request.GenerateTextRequest{
			Prompt: result,
		})
		if err != nil {
			l.logger.Error("Failed to generate text", zap.Error(err))
			continue
		}
		l.logger.Debug("Generated text", zap.String("title", generatedText.Title))

		// Send the generated article
		_, err = l.requestClient.Service.SendArticle(l.ctx, &request.SendArticleRequest{
			WebSite: request.WebSite_DTF,
			Entry: &request.EntryRequest{
				Id: 0,
				SubsiteId: 0,
				UserId: 0,
				Type: 1,
				Title: generatedText.Title,
				Entry: &request.Entry{
					Blocks: generatedText.Blocks,
				},
			},
			ConfigId: config.Id,
		})
		if err != nil {
			l.logger.Error("Failed to send article", zap.Error(err))
			continue
		}
		l.logger.Debug("Article sent successfully")

		// Delay if needed
		if config.Delay > 0 {
			delayDuration := time.Duration(config.Delay) * time.Hour
			l.logger.Debug("Delaying next row", zap.Duration("delay", delayDuration))
			select {
			case <-time.After(delayDuration):
			case <-l.ctx.Done():
				l.logger.Info("Worker context cancelled, stopping...")
				return
			}
		}
	}
}
