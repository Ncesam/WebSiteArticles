package queue

import (
	queuepb "backend/generated/proto/queue"
	requestPb "backend/generated/proto/request"
	"backend/internal/helpers"
	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/types"
	"context"
	"fmt"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)


type RabbitClient struct {
	Connection *amqp091.Connection
	Channel *amqp091.Channel
	InputQueue amqp091.Queue
	logger *zap.Logger
	cfg *config.Config
	ctx context.Context
}

func NewQueueClient(logger *zap.Logger, cfg *config.Config) (*RabbitClient, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.RABBIT.USERNAME, cfg.RABBIT.PASSWORD, cfg.RABBIT.HOST, cfg.RABBIT.PORT)
	connection, err := amqp091.Dial(dsn)
	if err != nil {
		logger.Error("Queue not connected", zap.Error(err))
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		logger.Error("Channel not initialize", zap.Error(err))
		return nil, err
	}
	client := &RabbitClient{
		logger: logger,
		cfg: cfg,
		Connection: connection,
		Channel: channel,
		ctx: context.Background(),
	}
	client.Migrate()
	return client, nil
}


func (client *RabbitClient) Migrate() error {
	inputQueue, err := client.Channel.QueueDeclare(
		"Articles",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		client.logger.Error("Migrate not initialize", zap.Error(err))
		return err
	}

	client.InputQueue = inputQueue
	return nil;
}

func (client *RabbitClient) StartConsume() (<-chan amqp091.Delivery, error) {
	return client.Channel.ConsumeWithContext(client.ctx, client.InputQueue.Name, "", true, false, false, false, nil)
}

func (client *RabbitClient) SendArticleToRequest(article types.Article) {
	
	var blocks []*queuepb.Block
	for _, block := range article.Blocks {
		var pbBlock *queuepb.Block

		switch block.Type {
			case "text":
				pbBlock = &queuepb.Block{
					Type: "text",
					Data: &queuepb.BlockData{Text: block.Data["text"].(string)},
				}

			case "code":
				pbBlock = &queuepb.Block{
					Type: "code",
					Data: &queuepb.BlockData{
						Text: block.Data["text"].(string),
						Lang: block.Data["lang"].(string),
					},
				}

			case "list":
				// Преобразуем []interface{} в []string
				itemsIface := block.Data["items"].([]interface{})
				items := make([]string, len(itemsIface))
				for i, v := range itemsIface {
					items[i] = v.(string)
				}
				pbBlock = &queuepb.Block{
					Type: "list",
					Data: &queuepb.BlockData{
						Items: items,
						Type:  block.Data["type"].(string),
					},
				}

			case "quote":
				pbBlock = &queuepb.Block{
					Type: "quote",
					Data: &queuepb.BlockData{
							Text:     block.Data["text"].(string),
						},
					},
		}

		if pbBlock != nil {
			blocks = append(blocks, pbBlock)
		}
	}
	_, err := request.Service.SendArticle(client.ctx, &requestPb.SendArticleRequest{
		WebSite: requestPb.WebSite_DTF,
		Title: article.Title,
		Blocks: article.Blocks,
	})
	if err != nil {
		helpers.HandleGrpcError(logger, c, err, "Failed to send article")
		return
	}
}

func (client *RabbitClient) ConsumeHandler(channel <-chan amqp091.Delivery) {
	for msg := range channel {
		client.logger.Debug("Recieved msg", zap.String("data", string(msg.Body)))
		body := string(msg.Body)
		
		client.SendArticleToRequest(&queuepb.Article{})
		msg.Ack(false)
	}
}

