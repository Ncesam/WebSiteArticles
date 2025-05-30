package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	configPb "backend/generated/proto/config"
	requestPb "backend/generated/proto/request"
	"backend/pkg/config"
	"backend/pkg/errors"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/types"
)

type RequestServer struct {
	requestPb.UnimplementedRequestServiceServer
	logger        *zap.Logger
	api           *http.Client
	cfg           *config.Config
	openaiClient  *openai.Client
	configClient  *grpcfactory.Client[configPb.ConfigServiceClient]
	dtfAPIBaseURL string
	vcAPIBaseURL  string
}

func NewRequestServer(logger *zap.Logger, cfg *config.Config, configClient *grpcfactory.Client[configPb.ConfigServiceClient]) *RequestServer {
	api := http.Client{}
	client := openai.NewClient(cfg.OPEN_AI_KEY)
	logger.Info("Created Request Server")
	return &RequestServer{
		logger:        logger,
		cfg:           cfg,
		api:           &api,
		openaiClient:  client,
		configClient:  configClient,
		dtfAPIBaseURL: "https://api.dtf.ru",
		vcAPIBaseURL:  "https://api.vc.ru",
	}
}
func (s *RequestServer) SendArticle(ctx context.Context, req *requestPb.SendArticleRequest) (*requestPb.Empty, error) {
	s.logger.Debug("SendArticle request received",
		zap.String("title", req.Entry.Title),
		zap.String("config_id", req.ConfigId),
		zap.String("website", req.WebSite.String()),
	)

	switch req.WebSite {
	case requestPb.WebSite_DTF:
		s.logger.Debug("Routing to DTF")
		return s.sendToSite(ctx, req, s.dtfAPIBaseURL)
	case requestPb.WebSite_VC_RU:
		s.logger.Debug("Routing to VC.ru")
		return s.sendToSite(ctx, req, s.vcAPIBaseURL)

	default:
		s.logger.Warn("Unsupported website value", zap.Int32("website_enum_value", int32(req.WebSite)))
		return nil, fmt.Errorf("unsupported website: %v", req.WebSite)
	}
}
func (s *RequestServer) sendToSite(ctx context.Context, req *requestPb.SendArticleRequest, baseURL string) (*requestPb.Empty, error) {
	s.logger.Info("Sending article", zap.Any("website", req.WebSite), zap.String("title", req.Entry.Title))

	token, err := s.fetchToken(ctx, req.ConfigId, baseURL)
	if err != nil {
		s.logger.Error("Failed to fetch token", zap.Error(err), zap.String("config_id", req.ConfigId), zap.Any("website", req.WebSite))
		return nil, err
	}

	subsiteID, err := s.fetchSubsite(ctx, token, baseURL)
	if err != nil {
		s.logger.Error("Failed to fetch subsite", zap.Error(err), zap.Any("website", req.WebSite))
		return nil, err
	}

	s.logger.Info("Successfully fetched subsite ID", zap.Int64("subsite_id", subsiteID))
	return s.postEntry(ctx, token, subsiteID, req.Entry, baseURL)
}
func (s *RequestServer) fetchToken(ctx context.Context, configID string, baseURL string) (string, error) {
	s.logger.Info("Fetching token", zap.String("config_id", configID), zap.String("base_url", baseURL))

	// 1) Получаем refreshToken из конфига
	cfg, err := s.configClient.Service.GetConfig(ctx, &configPb.GetConfigRequest{ConfigId: configID})
	if err != nil {
		s.logger.Error("GetConfig failed", zap.Error(err), zap.String("config_id", configID))
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
		s.logger.Error("Token request failed", zap.Error(err), zap.String("config_id", configID), zap.String("base_url", baseURL))
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
		s.logger.Error("Failed to parse token response", zap.Error(err), zap.String("config_id", configID), zap.String("base_url", baseURL))
		return "", errors.ErrInternalServer
	}
	s.logger.Info("Successfully fetched token", zap.String("access_token", out.Data.AccessToken))
	return out.Data.AccessToken, nil
}
func (s *RequestServer) fetchSubsite(ctx context.Context, accessToken, baseURL string) (int64, error) {
	endpoint := fmt.Sprintf("%s/v2.1/subsite/me", baseURL)
	s.logger.Info("Fetching subsite ID", zap.String("endpoint", endpoint))

	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		s.logger.Error("Failed to create HTTP request for subsite", zap.Error(err), zap.String("endpoint", endpoint))
		return 0, errors.ErrInternalServer
	}
	httpReq.Header.Set("jwtauthorization", "Bearer "+accessToken)

	resp, err := s.api.Do(httpReq)
	if err != nil {
		s.logger.Error("Subsite request failed", zap.Error(err), zap.String("endpoint", endpoint))
		return 0, err
	}
	defer resp.Body.Close()

	s.logger.Info("Subsite response received", zap.Int("status_code", resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read subsite response body", zap.Error(err))
		return 0, errors.ErrInternalServer
	}

	var out struct {
		Result struct {
			ID int64 `json:"id"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		s.logger.Error("Failed to parse subsite response JSON", zap.Error(err), zap.ByteString("body", body))
		return 0, errors.ErrInternalServer
	}

	if out.Result.ID == 0 {
		s.logger.Error("Subsite ID is missing or zero in response", zap.ByteString("body", body))
		return 0, errors.ErrInternalServer
	}

	s.logger.Info("Successfully fetched subsite ID", zap.Int64("subsite_id", out.Result.ID))
	return out.Result.ID, nil
}
func (s *RequestServer) createEntry(ctx context.Context, subsiteID int64, entry *requestPb.EntryRequest) (*types.EntryRequest, error) {
	s.logger.Debug("Creating entry", zap.Int64("subsite_id", subsiteID))
	protoMarshal := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		AllowPartial:    true,
	}
	entry.UserId = subsiteID
	entry.SubsiteId = subsiteID
	entryJson, err := protoMarshal.Marshal(entry)
	if err != nil {
		s.logger.Error("Failed to marshal entry", zap.Error(err))
		return nil, err
	}
	s.logger.Debug(string(entryJson))
	var entryMap types.EntryRequest
	err = json.Unmarshal(entryJson, &entryMap)
	if err != nil {
		s.logger.Error("Failed to unmarshal entry", zap.Error(err))
		return nil, err
	}

	entryMap.IsEnabledComments = true
	entryMap.IsEnabledLikes = true
	entryMap.IsEnabledAd = true

	return &entryMap, nil
}
func (s *RequestServer) postEntry(ctx context.Context, accessToken string, subsiteID int64, entry *requestPb.EntryRequest, baseURL string) (*requestPb.Empty, error) {
	s.logger.Debug("Posting entry", zap.Int64("subsite_id", subsiteID), zap.String("base_url", baseURL), zap.String("entry_title", entry.Title))
	entryMap, err := s.createEntry(ctx, subsiteID, entry)
	if err != nil {
		s.logger.Error("Error to create entry", zap.Error(err))
		return nil, err
	}
	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(entryMap)
	if err != nil {
		s.logger.Error("Failed to marshal entry", zap.Error(err))
		return nil, err
	}

	reqJsonBytes := bytes.TrimSpace(buf.Bytes())
	s.logger.Debug(string(reqJsonBytes))

	w := multipart.NewWriter(buf)
	_ = w.WriteField("entry", string(reqJsonBytes))
	w.Close()

	endpoint := fmt.Sprintf("%s/v2.1/editor", baseURL)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", endpoint, buf)
	httpReq.Header.Set("jwtauthorization", "Bearer "+accessToken)
	httpReq.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := s.api.Do(httpReq)
	if err != nil {
		s.logger.Error("Post entry failed", zap.Error(err), zap.String("base_url", baseURL), zap.Int64("subsite_id", subsiteID))
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Status code isn't OK", zap.Int("status_code", resp.StatusCode))
		return nil, errors.ErrInternalServer
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read subsite response body", zap.Error(err))
		return nil, errors.ErrInternalServer
	}
	s.logger.Debug(string(body))

	s.logger.Debug("Successfully posted entry", zap.String("base_url", baseURL), zap.Int64("subsite_id", subsiteID))
	return &requestPb.Empty{}, nil
}
func (s *RequestServer) UploadMedia(ctx context.Context, req *requestPb.UploadMediaRequest) (*requestPb.UploadMediaResponse, error) {
	var endpoint string
	s.logger.Debug("Uploading media")
	switch req.WebSite {
	case requestPb.WebSite_DTF:
		endpoint = "https://upload.dtf.ru/v2.8/uploader/upload"
	case requestPb.WebSite_VC_RU:
		endpoint = "https://upload.vc.ru/v2.8/uploader/upload"
	default:
		endpoint = ""
	}
	resp, err := s.api.Get(req.UrlFile)
	if err != nil {
		s.logger.Error("Failed to get File", zap.Error(err))
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Status code isn't OK", zap.Int("status_code", resp.StatusCode))
		return nil, errors.ErrInternalServer
	}

	defer resp.Body.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part1, errFile1 := writer.CreateFormFile("files_0", "")
	_, errFile1 = io.Copy(part1, resp.Body)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return nil, errFile1
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, payload)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	respUpload, err := s.api.Do(httpReq)
	if err != nil {
		s.logger.Error("Upload failed", zap.Error(err))
		return nil, err
	}
	defer respUpload.Body.Close()

	if respUpload.StatusCode != http.StatusOK {
		s.logger.Error("Upload failed", zap.Int("status_code", respUpload.StatusCode))
		return nil, errors.ErrInternalServer
	}

	body, err := io.ReadAll(respUpload.Body)
	if err != nil {
		s.logger.Error("Failed to read upload response", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("Upload response: " + string(body))

	var uploadResponse requestPb.ImageItem
	err = json.Unmarshal(body, &uploadResponse)
	if err != nil {
		s.logger.Error("Failed to unmarshal body", zap.Error(err))
		return nil, err
	}
	return &requestPb.UploadMediaResponse{Info: &uploadResponse}, nil
}
func (s *RequestServer) GenerateText(ctx context.Context, req *requestPb.GenerateTextRequest) (*requestPb.Text, error) {
	s.logger.Info("Generating text from prompt", zap.String("prompt", req.Prompt))

	systemMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "Ты — помощник, генерирующий статью для публикации. Выводи все в markdown. Также добавь эту партнерскую ссылку " + req.Link,
	}
	userMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Prompt,
	}

	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    []openai.ChatCompletionMessage{systemMsg, userMsg},
			MaxTokens:   1024,
			Temperature: 0.7,
		},
	)
	if err != nil {
		s.logger.Error("OpenAI request failed", zap.Error(err), zap.String("prompt", req.Prompt))
		return nil, err
	}

	text := resp.Choices[0].Message.Content
	paragraphs := strings.Split(text, "\n\n")
	blocks := make([]*requestPb.Block, 0, len(paragraphs))

	title := ""

	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}

		blockType := "text"
		blockText := trimmed

		if strings.HasPrefix(trimmed, "# ") {
			blockType = "header"
			blockText = strings.TrimPrefix(trimmed, "# ")
			title = blockText
		} else if strings.HasPrefix(trimmed, ">") {
			blockType = "quote"
			blockText = strings.TrimPrefix(trimmed, ">")
		}

		blocks = append(blocks, &requestPb.Block{
			Type:   blockType,
			Cover:  false,
			Hidden: false,
			Anchor: "",
			Data: &requestPb.BlockData{
				Text: string(markdown.ToHTML([]byte(blockText), nil, nil)),
			},
		})
	}

	s.logger.Info("Successfully generated text", zap.String("title", title))
	return &requestPb.Text{
		Title:  title,
		Blocks: blocks,
	}, nil
}
func (s *RequestServer) GetItemLink(ctx context.Context, req *requestPb.GetItemLinkRequest) (*requestPb.ItemLink, error) {
	s.logger.Debug("Received request to get item link", zap.String("item URI", req.URI))

	var endpoint string = "https://api.content.market.yandex.ru/v3/affiliate/partner/link/create"

	uri, err := url.Parse(endpoint)
	if err != nil {
		s.logger.Error("Failed to parse uri", zap.Error(err))
		return nil, err
	}
	q := uri.Query()
	q.Add("url", req.URI)
	q.Add("clid", s.cfg.CLID_MARKET)
	uri.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)
	if err != nil {
		s.logger.Error("Failed to create http request", zap.Error(err))
		return nil, err
	}
	resp, err := s.api.Do(httpReq)
	if err != nil {
		s.logger.Error("Failed to execute http request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Received non-OK response", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("Error status code %s", resp.StatusCode)
	}

	var apiResponse struct {
		Link struct {
			URL      string `json:"url"`
			ShortURL string `json:"shortUrl"`
		} `json:"link"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}
	return &requestPb.ItemLink{
		URI: apiResponse.Link.URL,
	}, nil
}
func (s *RequestServer) getItemLinkDTF(ctx context.Context, url string, apiUrl)
func (s *RequestServer) fetchRefreshToken(ctx context.Context, req *requestPb.GetRefreshTokenRequest, baseURL string) (string, error) {
	s.logger.Info("Starting to fetch refresh token", zap.String("email", req.Email), zap.String("baseURL", baseURL))

	if req.Email == "" || req.Password == "" {
		s.logger.Warn("Email or password is empty", zap.String("email", req.Email))
		return "", errors.ErrInvalidCredentials
	}

	// Собираем multipart payload
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	_ = w.WriteField("email", req.Email)
	_ = w.WriteField("password", req.Password)
	if err := w.Close(); err != nil {
		s.logger.Error("Failed to close multipart writer", zap.Error(err))
		return "", errors.ErrInvalidCredentials
	}

	// Формируем HTTP-запрос
	endpoint := fmt.Sprintf("%s/v3.4/auth/email/login", baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, buf)
	if err != nil {
		s.logger.Error("Failed to create HTTP request", zap.Error(err), zap.String("endpoint", endpoint))
		return "", errors.ErrInternalServer
	}
	httpReq.Header.Set("Accept", "*/*")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0")
	httpReq.Header.Set("Content-Type", w.FormDataContentType())

	s.logger.Info("Sending request to external API", zap.String("endpoint", endpoint))

	// Выполняем HTTP-запрос
	res, err := s.api.Do(httpReq)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.Error(err), zap.String("endpoint", endpoint))
		return "", errors.ErrBadRequest
	}
	defer res.Body.Close()

	s.logger.Info("Received response", zap.Int("status_code", res.StatusCode), zap.String("email", req.Email))

	// Читаем тело ответа
	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.Error(err))
		return "", errors.ErrInternalServer
	}

	// Парсим JSON
	var parsed struct {
		Data struct {
			RefreshToken string `json:"refreshToken"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		s.logger.Error("Failed to parse JSON response", zap.Error(err), zap.ByteString("body", body))
		return "", errors.ErrInternalServer
	}

	if parsed.Data.RefreshToken == "" {
		s.logger.Error("Refresh token not found in response", zap.ByteString("body", body))
		return "", errors.ErrInternalServer
	}

	s.logger.Info("Successfully fetched refresh token", zap.String("email", req.Email))
	return parsed.Data.RefreshToken, nil
}
func (s *RequestServer) GetRefreshTokenDTF(ctx context.Context, req *requestPb.GetRefreshTokenRequest) (*requestPb.GetRefreshTokenResponse, error) {
	s.logger.Info("Handling GetRefreshTokenDTF", zap.String("email", req.Email))
	token, err := s.fetchRefreshToken(ctx, req, "https://api.dtf.ru")
	if err != nil {
		s.logger.Error("Failed to get refresh token for DTF", zap.Error(err), zap.String("email", req.Email))
		return nil, err
	}
	return &requestPb.GetRefreshTokenResponse{RefreshToken: token}, nil
}
func (s *RequestServer) GetRefreshTokenVC(ctx context.Context, req *requestPb.GetRefreshTokenRequest) (*requestPb.GetRefreshTokenResponse, error) {
	s.logger.Info("Handling GetRefreshTokenVC", zap.String("email", req.Email))
	token, err := s.fetchRefreshToken(ctx, req, "https://api.vc.ru")
	if err != nil {
		s.logger.Error("Failed to get refresh token for VC", zap.Error(err), zap.String("email", req.Email))
		return nil, err
	}
	return &requestPb.GetRefreshTokenResponse{RefreshToken: token}, nil
}
