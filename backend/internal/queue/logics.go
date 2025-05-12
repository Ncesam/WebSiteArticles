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
	logger.Info("Create a QueueLogics")
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

	go q.loop()
	for i := 0; i < 5; i++ { // 5 воркеров
		go q.worker(i)
	}

	return q
}

func (l *QueueLogics) loop() {
	for {
		select {
		case msg := <-l.inputChannel:
			l.logger.Debug("Received message")
			l.workerPool <- msg // отправляем задачу в пул на обработку
		case <-l.ctx.Done():
			l.logger.Info("Loop context cancelled, stopping...")
			return
		}
	}
}

func (l *QueueLogics) worker(id int) {
	for {
		select {
		case msg := <-l.workerPool:
			l.logger.Debug("Worker received task", zap.Int("worker_id", id))
			go l.processTask(msg)
		case <-l.ctx.Done():
			l.logger.Info("Worker stopped", zap.Int("worker_id", id))
			return
		}
	}
}

func (l *QueueLogics) processTask(msg types.InputForm) {
	config, err := l.configClient.Service.GetConfig(l.ctx, &configPb.GetConfigRequest{
		ConfigId: msg.ConfigId,
	})
	if err != nil {
		l.logger.Error("Failed to get config", zap.Error(err))
		return
	}

	refreshTokenDTF, err := l.requestClient.Service.GetRefreshTokenDTF(l.ctx, &request.GetRefreshTokenRequest{
		Email:    config.Email,
		Password: config.Password,
	})
	if err != nil {
		l.logger.Error("Failed to get refresh token", zap.Error(err))
		return
	}
	refreshTokenVC, err := l.requestClient.Service.GetRefreshTokenVC(l.ctx, &request.GetRefreshTokenRequest{
		Email:    config.Email,
		Password: config.Password,
	})
	if err != nil {
		l.logger.Error("Failed to get refresh token", zap.Error(err))
		return
	}
	_, err = l.configClient.Service.UpdateRefreshDTFToken(l.ctx, &configPb.UpdateRefreshTokenDTFRequest{
		RefreshTokenDTF: refreshTokenDTF.RefreshToken,
		Id: config.Id,
	})
	if err != nil {
		l.logger.Error("Failed to update refresh token", zap.Error(err))
		return
	}
	config, err = l.configClient.Service.UpdateRefreshVCToken(l.ctx, &configPb.UpdateRefreshTokenVCRequest{
		RefreshTokenVC: refreshTokenVC.RefreshToken,
		Id: config.Id,
	})
	if err != nil {
		l.logger.Error("Failed to update refresh token", zap.Error(err))
		return
	}


	var data map[string][]string
	if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
		l.logger.Error("Can't parse json", zap.Error(err))
		return
	}

	rowCount := len(data)
	for i := range rowCount {
		result := config.Prompt
		for key, values := range data {
			if i < len(values) {
				placeholder := "{" + key + "}"
				result = strings.ReplaceAll(result, placeholder, values[i])
			}
		}

		generatedText, err := l.requestClient.Service.GenerateText(l.ctx, &request.GenerateTextRequest{
			Prompt: result,
		})
		if err != nil {
			l.logger.Error("Can't generate text", zap.Error(err))
			continue
		}
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
			l.logger.Error("Can't send article", zap.Error(err))
			continue
		}
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
			l.logger.Error("Can't send article", zap.Error(err))
			continue
		}

		if config.Delay > 0 {
			delayDuration := time.Duration(config.Delay) * time.Hour
			l.logger.Debug("Delaying next row", zap.Duration("delay", delayDuration))
			select {
			case <-time.After(delayDuration):
			case <-l.ctx.Done():
				return
			}
		}
	}
}
