package providers

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/JazzJackrabbit/heimdall-ng/request"
	"github.com/JazzJackrabbit/heimdall-ng/response"
)

var aiandBaseURL = "https://api.aiand.com/v1"

type aiandRequest struct {
	Model               string         `json:"model"`
	Messages            any            `json:"messages"`
	Stream              bool           `json:"stream"`
	StreamOptions       streamOptions  `json:"stream_options"`
	Temperature         float32        `json:"temperature,omitempty"`
	TopP                float32        `json:"top_p,omitempty"`
	MaxCompletionTokens int            `json:"max_completion_tokens,omitempty"`
	ReasoningEffort     string         `json:"reasoning_effort,omitempty"`
	ResponseFormat      map[string]any `json:"response_format,omitempty"`
}

// aiandChunk is one streamed chat completion event. Reasoning models send
// their reasoning first, under both reasoning_content and reasoning (the same
// text), then the answer as content.
type aiandChunk struct {
	Choices []struct {
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type aiandMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type aiandTextPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type aiandFile struct {
	Filename string `json:"filename"`
	FileData string `json:"file_data"`
}

type aiandFilePart struct {
	Type string    `json:"type"`
	File aiandFile `json:"file"`
}

type AIAnd struct {
	apiKeys []string
}

func NewAIAnd(apiKeys []string) AIAnd {
	return AIAnd{
		apiKeys: apiKeys,
	}
}

func (a AIAnd) Name() string {
	return models.AIAndProvider
}

func (a AIAnd) doRequest(
	ctx context.Context,
	req request.Completion,
	client http.Client,
	chunkHandler func(chunk string) error,
	key string,
) (response.Completion, int, error) {
	model, ok := req.Model.(models.AIAndModel)
	if !ok {
		return response.Completion{}, http.StatusBadRequest, fmt.Errorf(
			"model %s is not an ai& model",
			req.Model.GetName(),
		)
	}
	params := model.GetAIAndParams()

	aiandReq := aiandRequest{
		Model:               model.GetName(),
		Stream:              true,
		StreamOptions:       streamOptions{IncludeUsage: true},
		Temperature:         req.Temperature,
		TopP:                req.TopP,
		MaxCompletionTokens: params.MaxOutputTokens,
		ReasoningEffort:     string(params.ReasoningEffort),
		Messages: prepareAIAndMessages(
			params,
			req.SystemMessage,
			req.UserMessage,
			req.History,
		),
	}

	if len(params.StructuredOutput) > 0 {
		aiandReq.ResponseFormat = map[string]any{
			"type":        "json_schema",
			"json_schema": params.StructuredOutput,
		}
	}

	body, err := json.Marshal(aiandReq)
	if err != nil {
		return response.Completion{}, 0, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/chat/completions", aiandBaseURL),
		bytes.NewReader(body))
	if err != nil {
		return response.Completion{}, 0, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+key)

	resp, err := client.Do(httpReq) //nolint:gosec // URL is a known API endpoint
	if err != nil {
		return response.Completion{}, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return response.Completion{}, resp.StatusCode, fmt.Errorf(
			"received status code %d: %s", resp.StatusCode, aiandErrorMessage(bodyBytes))
	}

	reader := bufio.NewReader(resp.Body)
	var fullContent strings.Builder
	var thoughts strings.Builder
	var usage response.Usage
	var rawEvents []json.RawMessage
	finishReason := ""

	for {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return response.Completion{}, 0, fmt.Errorf("read line: %w", err)
		}

		// Only data lines carry events; blank separators, "event:" lines and
		// ":" keep-alive comments are skipped.
		data, isData := strings.CutPrefix(strings.TrimSpace(line), "data:")
		data = strings.TrimSpace(data)
		if isData && data != "" && data != "[DONE]" {
			var chunk aiandChunk
			if jsonErr := json.Unmarshal([]byte(data), &chunk); jsonErr != nil {
				return response.Completion{}, 0, fmt.Errorf("unmarshal chunk: %w", jsonErr)
			}
			if chunk.Error != nil {
				// Mid-stream errors are reported after a 200, so they are
				// treated as server errors and retried.
				return response.Completion{}, http.StatusBadGateway, fmt.Errorf(
					"stream error: %s", chunk.Error.Message)
			}

			rawEvents = append(rawEvents, json.RawMessage(data))

			if len(chunk.Choices) > 0 {
				choice := chunk.Choices[0]
				thoughts.WriteString(choice.Delta.ReasoningContent)
				if choice.FinishReason != "" {
					finishReason = choice.FinishReason
				}
				if choice.Delta.Content != "" {
					fullContent.WriteString(choice.Delta.Content)
					if chunkHandler != nil {
						if err := chunkHandler(choice.Delta.Content); err != nil {
							return response.Completion{}, 0, err
						}
					}
				}
			}

			if chunk.Usage.TotalTokens != 0 {
				usage = response.Usage{
					PromptTokens:     chunk.Usage.PromptTokens,
					CompletionTokens: chunk.Usage.CompletionTokens,
					TotalTokens:      chunk.Usage.TotalTokens,
				}
			}
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	// max_completion_tokens covers reasoning and answer together, so a low
	// cap can end the reply before any answer text.
	if fullContent.Len() == 0 && finishReason == "length" {
		return response.Completion{}, http.StatusBadRequest, errors.New(
			"the output token limit was reached during reasoning, before any answer: raise MaxOutputTokens or lower ReasoningEffort",
		)
	}

	rawResp, err := json.Marshal(rawEvents)
	if err != nil {
		return response.Completion{}, 0, fmt.Errorf("marshal raw response events: %w", err)
	}

	return response.Completion{
		Content:     fullContent.String(),
		Thoughts:    thoughts.String(),
		Model:       model.GetName(),
		Usage:       usage,
		RawRequest:  body,
		RawResponse: rawResp,
	}, 0, nil
}

// aiandErrorMessage pulls the message out of ai&'s OpenAI-shaped error body,
// falling back to the raw body.
func aiandErrorMessage(body []byte) string {
	var parsed struct {
		Error struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.Error.Message == "" {
		return string(body)
	}
	if parsed.Error.Code != "" {
		return fmt.Sprintf("%s (%s)", parsed.Error.Message, parsed.Error.Code)
	}
	return parsed.Error.Message
}

// prepareAIAndMessages builds the chat messages. Images and PDFs go on the
// final user message, ahead of its text.
func prepareAIAndMessages(
	params models.AIAndParams,
	systemInst string,
	userMsg string,
	history []request.Message,
) []aiandMessage {
	messages := make([]aiandMessage, 0, len(history)+2)

	if systemInst != "" {
		messages = append(messages, aiandMessage{Role: "system", Content: systemInst})
	}

	for _, his := range history {
		messages = append(messages, aiandMessage{Role: his.Role, Content: his.Content})
	}

	if len(params.ImageFile) == 0 && len(params.PdfFiles) == 0 {
		messages = append(messages, aiandMessage{Role: "user", Content: userMsg})
		return messages
	}

	parts := make([]any, 0, len(params.ImageFile)+len(params.PdfFiles)+1)
	for _, img := range params.ImageFile {
		detail := "auto"
		if img.Detail != "" {
			detail = img.Detail
		}
		parts = append(parts, imageInput{
			Type:     "image_url",
			ImageURL: imageURL{URL: img.URL, Detail: detail},
		})
	}
	for i, pdf := range params.PdfFiles {
		parts = append(parts, aiandFilePart{
			Type: "file",
			File: aiandFile{
				Filename: fmt.Sprintf("document_%d.pdf", i+1),
				FileData: "data:application/pdf;base64," + string(pdf),
			},
		})
	}
	parts = append(parts, aiandTextPart{Type: "text", Text: userMsg})

	messages = append(messages, aiandMessage{Role: "user", Content: parts})
	return messages
}

func (a AIAnd) tryWithBackup(
	ctx context.Context,
	req request.Completion,
	client http.Client,
	chunkHandler func(chunk string) error,
	requestLog *response.Logging,
) (response.Completion, error) {
	key := a.apiKeys[0]

	maxRetries := 5
	initialBackoff := 100 * time.Millisecond
	maxBackoff := 10 * time.Second

	var lastErr error
	for attempt := range maxRetries {
		requestLog.Events = append(requestLog.Events, response.Event{
			Timestamp: time.Now(),
			Description: fmt.Sprintf(
				"attempting to complete request with exponential backoff. attempt: %v",
				attempt,
			),
		})

		select {
		case <-ctx.Done():
			requestLog.Events = append(requestLog.Events, response.Event{
				Timestamp: time.Now(),
				Description: fmt.Sprintf(
					"context was cancelled with error: %v",
					ctx.Err(),
				),
			})
			return response.Completion{}, ctx.Err()
		default:
			res, resCode, err := a.doRequest(
				ctx,
				req,
				client,
				chunkHandler,
				key,
			)
			if err == nil {
				return res, nil
			}
			requestLog.Events = append(requestLog.Events, response.Event{
				Timestamp: time.Now(),
				Description: fmt.Sprintf(
					"request could not be completed, err: %v",
					err,
				),
			})

			if !isRetryableError(resCode) {
				requestLog.Events = append(requestLog.Events, response.Event{
					Timestamp: time.Now(),
					Description: fmt.Sprintf(
						"request was not retryable due to err: %v",
						err,
					),
				})
				return response.Completion{}, err
			}

			lastErr = err

			backoff := min(initialBackoff*time.Duration(
				1<<attempt,
			), maxBackoff)

			var randomBytes [8]byte
			var jitter time.Duration
			if _, err := rand.Read(randomBytes[:]); err != nil {
				jitter = backoff
			} else {
				randFloat := float64(binary.LittleEndian.Uint64(randomBytes[:])) / (1 << 64)
				jitter = time.Duration(float64(backoff) * (0.8 + 0.4*randFloat))
			}

			timer := time.NewTimer(jitter)
			select {
			case <-ctx.Done():
				timer.Stop()
				return response.Completion{}, ctx.Err()
			case <-timer.C:
				continue
			}
		}
	}

	return response.Completion{}, fmt.Errorf(
		"max retries exceeded: %w",
		lastErr,
	)
}

func (a AIAnd) CompleteResponse(
	ctx context.Context,
	req request.Completion,
	client http.Client,
	requestLog *response.Logging,
) (response.Completion, error) {
	reqLog := &response.Logging{}
	if requestLog == nil {
		if req.Tags == nil {
			req.Tags = make(map[string]string)
		}
		req.Tags["request_type"] = "completion"

		reqLog = &response.Logging{
			Events: []response.Event{
				{
					Timestamp:   time.Now(),
					Description: "start of call to CompleteResponse",
				},
			},
			SystemMsg: req.SystemMessage,
			UserMsg:   req.UserMessage,
			Start:     time.Now(),
		}
	}
	if requestLog != nil {
		reqLog = requestLog
	}

	return a.withKeys(ctx, req, client, nil, reqLog)
}

func (a AIAnd) StreamResponse(
	ctx context.Context,
	client http.Client,
	req request.Completion,
	chunkHandler func(chunk string) error,
	requestLog *response.Logging,
) (response.Completion, error) {
	reqLog := &response.Logging{}
	if requestLog == nil {
		if req.Tags == nil {
			req.Tags = make(map[string]string)
		}
		req.Tags["request_type"] = "streaming"

		reqLog = &response.Logging{
			Events: []response.Event{
				{
					Timestamp:   time.Now(),
					Description: "start of call to StreamResponse",
				},
			},
			SystemMsg: req.SystemMessage,
			UserMsg:   req.UserMessage,
			Start:     time.Now(),
		}
	}
	if requestLog != nil {
		reqLog = requestLog
	}

	return a.withKeys(ctx, req, client, chunkHandler, reqLog)
}

// withKeys tries each API key once, then retries the first key with backoff.
// A request ai& rejected as invalid (400) fails the same way on every key, so
// it returns at once.
func (a AIAnd) withKeys(
	ctx context.Context,
	req request.Completion,
	client http.Client,
	chunkHandler func(chunk string) error,
	reqLog *response.Logging,
) (response.Completion, error) {
	for i, key := range a.apiKeys {
		reqLog.Events = append(reqLog.Events, response.Event{
			Timestamp: time.Now(),
			Description: fmt.Sprintf(
				"attempting to complete request with key_number: %v",
				i,
			),
		})
		res, resCode, err := a.doRequest(ctx, req, client, chunkHandler, key)
		if err == nil {
			return res, nil
		}

		reqLog.Events = append(reqLog.Events, response.Event{
			Timestamp: time.Now(),
			Description: fmt.Sprintf(
				"request could not be completed, err: %v",
				err,
			),
		})

		if resCode == http.StatusBadRequest {
			return response.Completion{}, err
		}
	}

	return a.tryWithBackup(ctx, req, client, chunkHandler, reqLog)
}

var _ LLMProvider = new(AIAnd)
