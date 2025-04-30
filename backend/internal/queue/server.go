package queue

import (
	queuepb "backend/generated/proto/queue"
	"backend/pkg/config"
	"bytes"
	"context"

	"github.com/rabbitmq/amqp091-go"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.uber.org/zap"
)



type QueueServer struct {
	queuepb.UnimplementedQueueServiceServer
	
	requestClient *
	rabbitClient *RabbitClient
	logger *zap.Logger
	cfg *config.Config
	ctx context.Context
}


func NewServer(logger *zap.Logger, cfg *config.Config, rabbitClient *RabbitClient) *QueueServer{
	return &QueueServer{
		rabbitClient: rabbitClient,
		logger: logger,
		cfg: cfg,
		ctx: context.Background(),
	}
}

func ConvertMarkdownToBlocks(md string) []queuepb.Block {
	var blocks []queuepb.Block
	parser := goldmark.New()
	root := parser.Parser().Parse(text.NewReader([]byte(md)))

	for node := root.FirstChild(); node != nil; node = node.NextSibling() {
		switch node.Kind() {
		case ast.KindParagraph:
			textContent := extractText(node, md)
			blocks = append(blocks, queuepb.Block{
				Type:   "text",
				Data: &queuepb.BlockData{
					Text: "<p>" + textContent + "</p>",
				},
			})

		case ast.KindFencedCodeBlock:
			code := string(node.Text([]byte(md)))
			lang := string(node.(*ast.FencedCodeBlock).Language([]byte(md)))
			blocks = append(blocks, queuepb.Block{
				Type:   "code",
				Data: &queuepb.BlockData{
					Text: code,
					Lang: lang,
				},
			})

		case ast.KindList:
			list := node.(*ast.List)
			var items []string
			for item := list.FirstChild(); item != nil; item = item.NextSibling() {
				text := extractText(item, md)
				items = append(items, text)
			}
			blocks = append(blocks, queuepb.Block{
				Type:   "list",
				Data: &queuepb.BlockData{
					Items: items,
					Type:  "OL", // or "OL" based on list.IsOrdered()
				},
			})
		}
	}

	return blocks
}

func extractText(n ast.Node, source string) string {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		switch t := c.(type) {
		case *ast.Text:
			buf.Write(t.Segment.Value([]byte(source)))
		default:
			buf.WriteString(extractText(c, source))
		}
	}
	return buf.String()
}

func StartController(rabbitClient *RabbitClient, requestClient *) {
	
}


func handler() {
	
	rabbitClient.Channel.PublishWithContext(
		rabbitClient.ctx,
		"", 
		rabbitClient.InputQueue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "applicatio",
		}

	)
}
