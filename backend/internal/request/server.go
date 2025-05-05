package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	configPb "backend/generated/proto/config"
	requestPb "backend/generated/proto/request"
	"backend/pkg/config"
	"backend/pkg/errors"
	grpcfactory "backend/pkg/grpcFactory"
)

type RequestServer struct {
	requestPb.UnimplementedRequestServiceServer
	logger        *zap.Logger
	api           *http.Client
	cfg			  *config.Config
	openaiClient *openai.Client
	configClient  *grpcfactory.Client[configPb.ConfigServiceClient]
	dtfAPIBaseURL string
	vcAPIBaseURL  string
}


func NewRequestServer(logger *zap.Logger, cfg *config.Config, configClient *grpcfactory.Client[configPb.ConfigServiceClient]) *RequestServer {
	api := http.Client{}
	client := openai.NewClient(cfg.OPEN_AI_KEY)
	return &RequestServer{
		logger:        logger,
		cfg:           cfg,
		api:           &api,
		openaiClient: client,
		configClient:  configClient,
		dtfAPIBaseURL: "https://api.dtf.ru",
		vcAPIBaseURL:  "https://api.vc.ru",
	}
}

func (s *RequestServer) SendArticle(ctx context.Context, req *requestPb.SendArticleRequest) (*requestPb.Empty, error) {
	s.logger.Info("SendArticle", zap.Any("website", req.WebSite), zap.String("title", req.Entry.Title))
	switch req.WebSite {
	case requestPb.WebSite_DTF:
		return s.sendToSite(ctx, req, s.dtfAPIBaseURL)
	case requestPb.WebSite_VC_RU:
		return s.sendToSite(ctx, req, s.vcAPIBaseURL)
	default:
		return nil, fmt.Errorf("unsupported website: %v", req.WebSite)
	}
}
func (s *RequestServer) sendToSite(ctx context.Context, req *requestPb.SendArticleRequest, baseURL string) (*requestPb.Empty, error) {
	token, err := s.fetchToken(ctx, req.ConfigId, baseURL)
	if err != nil {
		return nil, err
	}

	subsiteID, err := s.fetchSubsite(ctx, token, baseURL)
	if err != nil {
		return nil, err
	}

	return s.postEntry(ctx, token, subsiteID, req.Entry, baseURL)
}

func (s *RequestServer) fetchToken(ctx context.Context, configID int64, baseURL string) (string, error) {
	// 1) Получаем refreshToken из конфига
	cfg, err := s.configClient.Service.GetConfig(ctx, &configPb.GetConfigRequest{ConfigId: configID})
	if err != nil {
		s.logger.Error("GetConfig failed", zap.Error(err))
		return "", err
	}

	// 2) Multipart-запрос
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	var refreshToken string
	switch baseURL {
	case "https://api.dtf.ru":
		refreshToken = cfg.RefreshTokenDTF
	case "https://api.vc.ru":
		refreshToken = cfg.RefreshTokenVC
	}
	_ = w.WriteField("token", refreshToken)
	w.Close()

	endpoint := fmt.Sprintf("%s/v3.4/auth/refresh", baseURL)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", endpoint, buf)
	httpReq.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := s.api.Do(httpReq)
	if err != nil {
		s.logger.Error("Token request failed", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		Data struct {
			AccessToken string `json:"accessToken"`
		} `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &out); err != nil {
		s.logger.Error("Failed to parse token response", zap.Error(err))
		return "", errors.ErrInternalServer
	}
	return out.Data.AccessToken, nil
}

func (s *RequestServer) fetchSubsite(ctx context.Context, accessToken, baseURL string) (int64, error) {
	endpoint := fmt.Sprintf("%s/v2.1/subsite/me", baseURL)
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	httpReq.Header.Set("jwtauthorization", "Bearer "+accessToken)

	resp, err := s.api.Do(httpReq)
	if err != nil {
		s.logger.Error("Subsite request failed", zap.Error(err))
		return 0, err
	}
	defer resp.Body.Close()

	var out struct {
		Result struct {
			ID int64 `json:"id"`
		} `json:"result"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &out); err != nil {
		s.logger.Error("Failed to parse subsite response", zap.Error(err))
		return 0, errors.ErrInternalServer
	}
	return out.Result.ID, nil
}

func (s *RequestServer) postEntry(
	ctx context.Context,
	accessToken string,
	subsiteID int64,
	entry *requestPb.EntryRequest, // замените на реальный тип вашего entry
	baseURL string,
) (*requestPb.Empty, error) {
	// 1) Встраиваем subsiteID в тело
	//    предполагаем, что entry содержит поле SubsiteId int32
	entry.SubsiteId = int32(subsiteID)

	// 2) Упаковываем в JSON
	bodyMap := map[string]any{"entry": entry}
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	// 3) Multipart-запрос
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	_ = w.WriteField("entry", string(jsonBytes))
	w.Close()

	endpoint := fmt.Sprintf("%s/v2.1/editor", baseURL)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", endpoint, buf)
	httpReq.Header.Set("jwtauthorization", "Bearer "+accessToken)
	httpReq.Header.Set("Content-Type", w.FormDataContentType())

	if _, err := s.api.Do(httpReq); err != nil {
		s.logger.Error("Post entry failed", zap.Error(err))
		return nil, err
	}
	return &requestPb.Empty{}, nil
}

// UploadMedia — заглушка загрузки медиа
func (s *RequestServer) UploadMedia(ctx context.Context, req *requestPb.UploadMediaRequest) (*requestPb.UploadMediaResponse, error) {
	s.logger.Info("Uploading media", zap.String("type", req.Type))
	// Можешь здесь сохранить файл и вернуть UID
	return &requestPb.UploadMediaResponse{UID: "mocked-uid"}, nil
}

// GenerateText — генерация текста на основе промпта (можно подключить GPT)
func (s *RequestServer) GenerateText(ctx context.Context, req *requestPb.GenerateTextRequest) (*requestPb.Text, error) {
	s.logger.Info("Generating text from prompt", zap.String("prompt", req.Prompt))

    // 1) Собираем сообщение для Chat API
    systemMsg := openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleSystem,
        Content: "Ты — помощник, генерирующий структуру статьи для публикации.",
    }
    userMsg := openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: req.Prompt,
    }

    // 2) Отправляем запрос
    resp, err := s.openaiClient.CreateChatCompletion(
        ctx,
        openai.ChatCompletionRequest{
            Model:     openai.GPT3Dot5Turbo, // или "gpt-4"
            Messages:  []openai.ChatCompletionMessage{systemMsg, userMsg},
            MaxTokens: 1024,
            Temperature: 0.7,
        },
    )
    if err != nil {
        s.logger.Error("OpenAI request failed", zap.Error(err))
        return nil, err
    }

    // 3) Парсим ответ в блоки — например, ожидая, что ассистент вернёт JSON-массив блоков
    //    или просто разбиваем на параграфы
    text := resp.Choices[0].Message.Content

    // Простая логика: каждый абзац — отдельный блок
    paragraphs := strings.Split(text, "\n\n")
    blocks := make([]*requestPb.Block, 0, len(paragraphs))
    for _, p := range paragraphs {
        if trimmed := strings.TrimSpace(p); trimmed != "" {
            blocks = append(blocks, &requestPb.Block{
                Type: "paragraph",
                Data: &requestPb.BlockData{
                    Text: trimmed,
                },
            })
        }
    }

    return &requestPb.Text{
        Title:  paragraphs[0],
        Blocks: blocks,
    }, nil
}

func (s *RequestServer) fetchRefreshToken(ctx context.Context, req *requestPb.GetRefreshTokenRequest, baseURL string) (string, error) {
    s.logger.Info("Fetching refresh token", zap.String("email", req.Email), zap.String("baseURL", baseURL))

    if req.Email == "" || req.Password == "" {
        return "", errors.ErrInvalidCredentials
    }

    // Собираем multipart payload
    buf := &bytes.Buffer{}
    w := multipart.NewWriter(buf)
    _ = w.WriteField("email", req.Email)
    _ = w.WriteField("password", req.Password)
    if err := w.Close(); err != nil {
        s.logger.Error("Failed to write form fields", zap.Error(err))
        return "", errors.ErrInvalidCredentials
    }

    // Формируем запрос
    endpoint := fmt.Sprintf("%s/v3.4/auth/email/login", baseURL)
    httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, buf)
    if err != nil {
        s.logger.Error("Failed to create http request", zap.Error(err))
        return "", errors.ErrInternalServer
    }
    httpReq.Header.Set("Accept", "*/*")
    httpReq.Header.Set("User-Agent", "Mozilla/5.0")
    httpReq.Header.Set("Content-Type", w.FormDataContentType())

    // Выполняем
    res, err := s.api.Do(httpReq)
    if err != nil {
        s.logger.Error("Request error", zap.Error(err))
        return "", errors.ErrBadRequest
    }
    defer res.Body.Close()

    // Читаем тело
    body, err := io.ReadAll(res.Body)
    if err != nil {
        s.logger.Error("Error reading body", zap.Error(err))
        return "", errors.ErrInternalServer
    }

    // Парсим JSON
    var parsed struct {
        Data struct {
            RefreshToken string `json:"refreshToken"`
        } `json:"data"`
    }
    if err := json.Unmarshal(body, &parsed); err != nil {
        s.logger.Error("Error parsing JSON", zap.Error(err))
        return "", errors.ErrInternalServer
    }

    if parsed.Data.RefreshToken == "" {
        s.logger.Error("Refresh token not found in response")
        return "", errors.ErrInternalServer
    }
    return parsed.Data.RefreshToken, nil
}

func (s *RequestServer) GetRefreshTokenDTF(ctx context.Context, req *requestPb.GetRefreshTokenRequest) (*requestPb.GetRefreshTokenResponse, error) {
    token, err := s.fetchRefreshToken(ctx, req, "https://api.dtf.ru")
    if err != nil {
        return nil, err
    }
    return &requestPb.GetRefreshTokenResponse{RefreshToken: token}, nil
}

func (s *RequestServer) GetRefreshTokenVC(ctx context.Context, req *requestPb.GetRefreshTokenRequest) (*requestPb.GetRefreshTokenResponse, error) {
    token, err := s.fetchRefreshToken(ctx, req, "https://api.vc.ru")
    if err != nil {
        return nil, err
    }
    return &requestPb.GetRefreshTokenResponse{RefreshToken: token}, nil
}

