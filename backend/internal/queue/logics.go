package queue

import (
	"context"
	"encoding/json"
	"fmt"
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
	ctx := context.Background()

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
			l.logger.Debug("Received message")
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
			l.logger.Debug("Worker received task", zap.Int("worker_id", id))
			go l.processTask(msg)
		case <-l.ctx.Done():
			l.logger.Info("Worker stopped", zap.Int("worker_id", id))
			return
		}
	}
}
func (l *QueueLogics) processTask(msg types.InputForm) {
	l.logger.Debug("Processing task", zap.String("config_id", msg.ConfigId))

	config, err := l.fetchConfig(msg.ConfigId)
	if err != nil {
		return
	}

	if err := l.updateRefreshTokens(config); err != nil {
		return
	}

	items, err := l.parseItems(msg.Data)
	if err != nil {
		return
	}

	for idx, item := range items {
		for {
			if err := l.processItem(idx, item, config); err != nil {
				l.logger.Error("Failed to process item", zap.Int("item_index", idx), zap.Error(err))
				time.Sleep(10 * time.Minute)
				continue
			}

			break
		}

		if config.Delay > 0 && !l.sleepWithContext(int32(config.Delay)) {
			return
		}
	}
}
func (l *QueueLogics) fetchConfig(configID string) (*configPb.Config, error) {
	config, err := l.configClient.Service.GetConfig(l.ctx, &configPb.GetConfigRequest{ConfigId: configID})
	if err != nil {
		l.logger.Error("Failed to get config", zap.Error(err))
		return nil, err
	}
	l.logger.Debug("Fetched config", zap.String("config_id", config.Id))
	return config, nil
}
func (l *QueueLogics) updateRefreshTokens(config *configPb.Config) error {
	refreshDTF, err := l.requestClient.Service.GetRefreshTokenDTF(l.ctx, &request.GetRefreshTokenRequest{
		Email: config.Email, Password: config.Password,
	})
	if err != nil {
		l.logger.Error("Failed to get refresh token DTF", zap.Error(err))
		return err
	}
	refreshVC, err := l.requestClient.Service.GetRefreshTokenVC(l.ctx, &request.GetRefreshTokenRequest{
		Email: config.Email, Password: config.Password,
	})
	if err != nil {
		l.logger.Error("Failed to get refresh token VC", zap.Error(err))
		return err
	}

	if _, err := l.configClient.Service.UpdateRefreshDTFToken(l.ctx, &configPb.UpdateRefreshTokenDTFRequest{
		RefreshTokenDTF: refreshDTF.RefreshToken, Id: config.Id,
	}); err != nil {
		l.logger.Error("Failed to update DTF refresh token", zap.Error(err))
		return err
	}

	if _, err := l.configClient.Service.UpdateRefreshVCToken(l.ctx, &configPb.UpdateRefreshTokenVCRequest{
		RefreshTokenVC: refreshVC.RefreshToken, Id: config.Id,
	}); err != nil {
		l.logger.Error("Failed to update VC refresh token", zap.Error(err))
		return err
	}

	return nil
}
func (l *QueueLogics) parseItems(data string) ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		l.logger.Error("Failed to parse JSON array", zap.Error(err))
		return nil, err
	}
	l.logger.Debug("Parsed JSON items count", zap.Int("count", len(items)))
	return items, nil
}
func (l *QueueLogics) processItem(idx int, data map[string]interface{}, config *configPb.Config) error {
	prompt := l.buildPrompt(config.Prompt, data)

	link, err := l.getItemLink(data["product_url"].(string))
	if err != nil {
		return fmt.Errorf("get link: %w", err)
	}

	urlImage, ok := data["images"].(string)
	if !ok {
		return fmt.Errorf("image url not found")
	}

	generated, err := l.requestClient.Service.GenerateText(l.ctx, &request.GenerateTextRequest{
		Prompt: prompt, Link: link.ShortURI,
	})
	if err != nil {
		return fmt.Errorf("generate text: %w", err)
	}

	switch config.WebSite {
	case configPb.WebSite_DTF:
		fileDTF, err := l.uploadImage(urlImage, request.WebSite_DTF)
		if err != nil {
			return fmt.Errorf("upload image DTF: %w", err)
		}

		if err := l.sendArticle(generated, fileDTF, request.WebSite_DTF, config.Id); err != nil {
			return err
		}
	case configPb.WebSite_VC_RU:
		fileVC, err := l.uploadImage(urlImage, request.WebSite_VC_RU)
		if err != nil {
			return fmt.Errorf("upload image VC: %w", err)
		}

		if err := l.sendArticle(generated, fileVC, request.WebSite_VC_RU, config.Id); err != nil {
			return err
		}
	}

	l.logger.Info("Article sent", zap.Int("item_index", idx))
	return nil
}
func (l *QueueLogics) buildPrompt(template string, data map[string]interface{}) string {
	for key, value := range data {
		placeholder := "{" + key + "}"
		template = strings.ReplaceAll(template, placeholder, fmt.Sprintf("%v", value))
	}
	return template
}
func (l *QueueLogics) uploadImage(url string, site request.WebSite) (*request.UploadMediaResponse, error) {
	l.logger.Debug("Upload image", zap.String("image_url", url))
	file, err := l.requestClient.Service.UploadMedia(l.ctx, &request.UploadMediaRequest{
		WebSite: site, UrlFile: url,
	})
	if err != nil {
		l.logger.Error("Failed to upload media", zap.String("site", site.String()), zap.Error(err))
		return nil, err
	}
	return file, nil
}
func (l *QueueLogics) sendArticle(text *request.Text, file *request.UploadMediaResponse, site request.WebSite, configId string) error {
	mediaBlock := &request.Block{
		Type: "media",
		Data: &request.BlockData{
			Items: []*request.MediaItem{{Image: file.Info}},
		},
	}
	var blocks []*request.Block

	blocks = append(blocks, mediaBlock)

	for _, block := range text.Blocks {
		blocks = append(blocks, block)
	}

	_, err := l.requestClient.Service.SendArticle(l.ctx, &request.SendArticleRequest{
		WebSite: site,
		Entry: &request.EntryRequest{
			Title: text.Title, Type: 1, Entry: &request.Entry{Blocks: blocks},
		},
		ConfigId: configId,
	})
	return err
}
func (l *QueueLogics) sleepWithContext(minute int32) bool {
	delay := time.Duration(minute) * time.Minute
	l.logger.Debug("Sleeping before next item", zap.Duration("delay", delay))
	select {
	case <-time.After(delay):
		return true
	case <-l.ctx.Done():
		l.logger.Info("Context cancelled, stopping worker")
		return false
	}
}
func (l *QueueLogics) getItemLink(url string) (*request.ItemLink, error) {
	l.logger.Debug("Getting item link", zap.String("item_url", url))
	itemLink, err := l.requestClient.Service.GetItemLink(l.ctx, &request.GetItemLinkRequest{
		URI: url,
	})
	if err != nil {
		l.logger.Error("Failed to get Item link", zap.Error(err))
		return nil, err
	}
	l.logger.Debug("Item link have got", zap.String("link", itemLink.ShortURI))
	return itemLink, nil
}
