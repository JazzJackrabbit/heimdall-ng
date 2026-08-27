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

var perplexityBaseUrl = "https://api.perplexity.ai/v1/agent"

type Perplexity struct {
	apiKeys []string
}

func NewPerplexity(apiKeys []string) Perplexity {
	return Perplexity{
		apiKeys,
	}
}

// CompleteResponse implements LLMProvider.
func (p Perplexity) CompleteResponse(
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

	for i, key := range p.apiKeys {
		reqLog.Events = append(reqLog.Events, response.Event{
			Timestamp: time.Now(),
			Description: fmt.Sprintf(
				"attempting to complete request with key_number: %v",
				i,
			),
		})
		res, _, err := p.doRequest(ctx, req, client, nil, key)
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
	}

	return p.tryWithBackup(ctx, req, client, nil, reqLog)
}

// perplexityPreset returns the Agent API preset Perplexity recommends in
// place of each Sonar model, together with the model's structured output
// schema. A preset selects the underlying model, tools and token limits; the
// response reports the model that ran.
func perplexityPreset(m models.Model) (string, map[string]any) {
	switch m := m.(type) {
	case models.Sonar:
		return "fast", m.StructuredOutput
	case models.SonarPro:
		return "low", m.StructuredOutput
	case models.SonarReasoningPro:
		return "medium", m.StructuredOutput
	case models.SonarReasoning: //nolint:staticcheck // backward compatibility
		return "medium", m.StructuredOutput
	default:
		return "", nil
	}
}

type perplexityInputItem struct {
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type perplexityRequest struct {
	Preset         string                `json:"preset"`
	Instructions   string                `json:"instructions,omitempty"`
	Input          []perplexityInputItem `json:"input"`
	Stream         bool                  `json:"stream"`
	Temperature    float32               `json:"temperature,omitempty"`
	TopP           float32               `json:"top_p,omitempty"`
	ResponseFormat map[string]any        `json:"response_format,omitempty"`
}

type perplexityError struct {
	Message string `json:"message"`
}

// perplexityEvent is the subset of an Agent API stream event the provider
// reads: text deltas, the completed response with its usage, and failures.
type perplexityEvent struct {
	Type     string           `json:"type"`
	Delta    string           `json:"delta"`
	Error    *perplexityError `json:"error"`
	Response struct {
		Model string `json:"model"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
			TotalTokens  int `json:"total_tokens"`
		} `json:"usage"`
		Error *perplexityError `json:"error"`
	} `json:"response"`
}

// doRequest implements LLMProvider.
func (p Perplexity) doRequest(
	ctx context.Context,
	req request.Completion,
	client http.Client,
	chunkHandler func(chunk string) error,
	key string,
) (response.Completion, int, error) {
	preset, structuredOutput := perplexityPreset(req.Model)
	if preset == "" {
		return response.Completion{}, 0, fmt.Errorf(
			"unsupported Perplexity model: %s",
			req.Model.GetName(),
		)
	}

	input := make([]perplexityInputItem, 0, len(req.History)+1)
	for _, his := range req.History {
		input = append(input, perplexityInputItem{
			Type:    "message",
			Role:    his.Role,
			Content: his.Content,
		})
	}
	input = append(input, perplexityInputItem{
		Type:    "message",
		Role:    "user",
		Content: req.UserMessage,
	})

	apiReq := perplexityRequest{
		Preset:       preset,
		Instructions: req.SystemMessage,
		Input:        input,
		Stream:       true,
		Temperature:  req.Temperature,
		TopP:         req.TopP,
	}

	if len(structuredOutput) > 0 {
		apiReq.ResponseFormat = map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":   "response",
				"schema": structuredOutput,
			},
		}
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return response.Completion{}, 0, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		perplexityBaseUrl,
		bytes.NewReader(body))
	if err != nil {
		return response.Completion{}, 0, fmt.Errorf(
			"create request: %w",
			err,
		)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+key)

	resp, err := client.Do(httpReq) //nolint:gosec // URL is a known API endpoint
	if err != nil {
		return response.Completion{}, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr == nil {
			return response.Completion{}, resp.StatusCode, fmt.Errorf(
				"perplexity returned status %d: %s",
				resp.StatusCode,
				string(bodyBytes),
			)
		}
		return response.Completion{}, resp.StatusCode, fmt.Errorf(
			"perplexity returned status %d",
			resp.StatusCode,
		)
	}

	reader := bufio.NewReader(resp.Body)
	var fullContent strings.Builder
	var usage response.Usage
	var rawEvents []json.RawMessage
	model := req.Model.GetName()
	chunks := 0
	now := time.Now()

	for {
		if chunks == 0 && time.Since(now).Seconds() > 3.0 {
			return response.Completion{}, 0, context.Canceled
		}
		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return response.Completion{}, 0, fmt.Errorf(
				"read line: %w",
				err,
			)
		}

		// Events arrive as "event:" and "data:" lines; the type is repeated
		// inside the data payload, so only data lines are read.
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if line == "" || line == "[DONE]" {
			continue
		}

		var event perplexityEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return response.Completion{}, 0, fmt.Errorf(
				"unmarshal event: %w",
				err,
			)
		}

		rawEvents = append(rawEvents, json.RawMessage(line))
		chunks++

		switch event.Type {
		case "response.output_text.delta":
			fullContent.WriteString(event.Delta)

			if chunkHandler != nil {
				if err := chunkHandler(event.Delta); err != nil {
					return response.Completion{}, 0, err
				}
			}
		case "response.completed":
			if event.Response.Model != "" {
				model = event.Response.Model
			}
			usage = response.Usage{
				PromptTokens:     event.Response.Usage.InputTokens,
				CompletionTokens: event.Response.Usage.OutputTokens,
				TotalTokens:      event.Response.Usage.TotalTokens,
			}
		case "response.failed", "error":
			message := line
			if event.Error != nil && event.Error.Message != "" {
				message = event.Error.Message
			} else if event.Response.Error != nil && event.Response.Error.Message != "" {
				message = event.Response.Error.Message
			}
			return response.Completion{}, 0, fmt.Errorf(
				"perplexity request failed: %s",
				message,
			)
		}
	}

	finalContent := fullContent.String()
	rawResp, err := json.Marshal(rawEvents)
	if err != nil {
		return response.Completion{}, 0, fmt.Errorf("marshal raw response events: %w", err)
	}

	return response.Completion{
		Content:     finalContent,
		Model:       model,
		Usage:       usage,
		RawRequest:  body,
		RawResponse: rawResp,
	}, 0, nil
}

func (p Perplexity) Name() string {
	return models.PerplexityProvider
}

// StreamResponse implements LLMProvider.
func (p Perplexity) StreamResponse(
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

	for i, key := range p.apiKeys {
		reqLog.Events = append(reqLog.Events, response.Event{
			Timestamp: time.Now(),
			Description: fmt.Sprintf(
				"attempting to complete request with key_number: %v",
				i,
			),
		})
		res, _, err := p.doRequest(ctx, req, client, chunkHandler, key)
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
	}

	return p.tryWithBackup(ctx, req, client, chunkHandler, reqLog)
}

// tryWithBackup implements LLMProvider.
func (p Perplexity) tryWithBackup(
	ctx context.Context,
	req request.Completion,
	client http.Client,
	chunkHandler func(chunk string) error,
	requestLog *response.Logging,
) (response.Completion, error) {
	key := p.apiKeys[0]

	maxRetries := 5
	initialBackoff := 100 * time.Millisecond
	maxBackoff := 10 * time.Second

	var lastErr error
	for attempt := range maxRetries {
		requestLog.Events = append(requestLog.Events, response.Event{
			Timestamp: time.Now(),
			Description: fmt.Sprintf(
				"attempting to complete request with expoential backoff. attempt: %v",
				attempt,
			),
		})

		select {
		case <-ctx.Done():
			requestLog.Events = append(requestLog.Events, response.Event{
				Timestamp: time.Now(),
				Description: fmt.Sprintf(
					"context was called with error: %v",
					ctx.Err(),
				),
			})
			return response.Completion{}, ctx.Err()
		default:
			res, resCode, err := p.doRequest(
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

var _ LLMProvider = new(Perplexity)
