# Heimdall

[![Go Reference](https://pkg.go.dev/badge/github.com/JazzJackrabbit/heimdall-ng.svg)](https://pkg.go.dev/github.com/JazzJackrabbit/heimdall-ng)
[![Release](https://img.shields.io/github/v/release/JazzJackrabbit/heimdall-ng)](https://github.com/JazzJackrabbit/heimdall-ng/releases)
[![Build](https://github.com/JazzJackrabbit/heimdall-ng/actions/workflows/build.yml/badge.svg)](https://github.com/JazzJackrabbit/heimdall-ng/actions/workflows/build.yml)
[![License: BSD-3](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](LICENSE)

Heimdall is a Go library that sends LLM requests through one interface to OpenAI, Anthropic, Google Gemini, Vertex AI, xAI Grok, Perplexity and OpenRouter. A request names a primary model and a list of fallbacks. The router tries each model in order, and each provider tries its API keys in order, until one succeeds.

## Features

- One request type for every provider, with streaming and non-streaming calls
- Fallback models across providers and API key failover within a provider
- PDF, image and file inputs
- JSON structured output
- Token usage and a per-request event log on every response
- Typed model definitions with per-model options and pricing

## Installation

```bash
go get github.com/JazzJackrabbit/heimdall-ng
```

Requires Go 1.26 or later.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	heimdall "github.com/JazzJackrabbit/heimdall-ng"
	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/JazzJackrabbit/heimdall-ng/providers"
	"github.com/JazzJackrabbit/heimdall-ng/request"
)

func main() {
	router := heimdall.New(60*time.Second, []heimdall.LLMProvider{
		providers.NewOpenAI([]string{os.Getenv("OPENAI_API_KEY")}),
		providers.NewAnthropic([]string{os.Getenv("ANTHROPIC_API_KEY")}),
		providers.NewGoogle([]string{os.Getenv("GOOGLE_API_KEY")}),
	})

	res, err := router.Complete(context.Background(), request.Completion{
		Model:         models.GPT6Sol{},
		Fallback:      []models.Model{models.Claude5Sonnet{}, models.Gemini38Flash{}},
		SystemMessage: "You are a concise assistant.",
		UserMessage:   "Explain what a vector database is in two sentences.",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Content)
	fmt.Printf("model: %s, tokens: %d\n", res.Model, res.Usage.TotalTokens)
}
```

## Usage

### Providers

Each provider is created with a list of API keys, which it tries in order. Register the providers you need with `heimdall.New`. A model is routed to the provider it belongs to, and models whose provider is not registered are skipped. The timeout passed to `heimdall.New` applies to each HTTP request, including the whole streamed response.

| Provider | Constructor | Model types |
| --- | --- | --- |
| OpenAI | `providers.NewOpenAI(keys)` | `models.GPT*` |
| Anthropic | `providers.NewAnthropic(keys)` | `models.Claude*` |
| Google Gemini | `providers.NewGoogle(keys)` | `models.Gemini*` |
| Vertex AI | `providers.NewVertexAI(ctx, projectID, location, credentialsJSON)` | `models.VertexGemini*` |
| xAI Grok | `providers.NewGrok(keys)` | `models.Grok*` |
| Perplexity | `providers.NewPerplexity(keys)` | `models.Sonar*` |
| OpenRouter | `providers.NewOpenRouter(keys)` | `models.OpenRouterModel{ModelName: "..."}` |

Vertex AI authenticates with a service account key instead of API keys:

```go
credentials, err := os.ReadFile("service-account.json")
if err != nil {
	panic(err)
}

vertex, err := providers.NewVertexAI(ctx, "my-project-id", "global", credentials)
if err != nil {
	panic(err)
}
```

### Model options

Options specific to a model are fields on its type, for example `MaxOutputTokens` on Claude 4.6 and later models (4096 when unset) and `ThinkingLevel` on Gemini 3 text models. The [models package reference](https://pkg.go.dev/github.com/JazzJackrabbit/heimdall-ng/models) lists the fields of each type.

```go
model := models.Claude5Sonnet{MaxOutputTokens: 16000}
```

`Temperature` and `TopP` on the request are sent to Anthropic models that accept them. Other providers ignore them.

Model types that implement `models.CostBreakdown` report their input and output prices per million tokens.

### Streaming

`router.Stream` takes the same request and calls the handler with each chunk of text as it arrives. It returns the full response when the stream ends. If a stream fails partway, the next key or model starts a new stream, and the handler has already received the chunks from the failed attempt.

```go
res, err := router.Stream(ctx, req, func(chunk string) error {
	fmt.Print(chunk)
	return nil
})
```

### Conversation history

Earlier turns go in `History`. The new message goes in `UserMessage`.

```go
req := request.Completion{
	Model:         models.Claude5Sonnet{},
	SystemMessage: "You are a travel assistant.",
	History: []request.Message{
		{Role: "user", Content: "I'm planning a trip to Japan."},
		{Role: "assistant", Content: "When are you planning to go?"},
	},
	UserMessage: "In April. What should I know?",
}
```

### Files and images

Files are set on the model value, since each provider accepts them in a different form. Content is base64 encoded.

```go
pdf, err := os.ReadFile("report.pdf")
if err != nil {
	panic(err)
}
encoded := base64.StdEncoding.EncodeToString(pdf)

req := request.Completion{
	Model: models.Claude5Sonnet{
		PdfFiles: []models.AnthropicPdf{models.AnthropicPdf(encoded)},
	},
	Fallback: []models.Model{
		models.Gemini38Flash{
			PdfFiles: []models.GooglePdf{models.GooglePdf(encoded)},
		},
		models.GPT6Sol{
			PdfFile: map[string]string{"report.pdf": "data:application/pdf;base64," + encoded},
		},
	},
	UserMessage: "Summarize this report.",
}
```

Images use the `ImageFile` field of the same models. OpenAI models take a data URL, Gemini models take raw base64 with its MIME type and Anthropic models take a map of media type to base64.

### Structured output

Set a JSON schema on the model's `StructuredOutput` field. OpenAI models take a named schema:

```go
model := models.GPT6Sol{
	StructuredOutput: map[string]any{
		"name": "city",
		"schema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":    map[string]any{"type": "string"},
				"country": map[string]any{"type": "string"},
			},
			"required": []string{"name", "country"},
		},
	},
}
```

Anthropic and Gemini models take the schema object directly, without the `name` and `schema` wrapper.

### Calling a provider directly

Every provider also implements `CompleteResponse` and `StreamResponse`. Use them to call one provider without the router's fallbacks.

```go
anthropic := providers.NewAnthropic([]string{os.Getenv("ANTHROPIC_API_KEY")})

res, err := anthropic.CompleteResponse(ctx, req, http.Client{Timeout: time.Minute}, nil)
```

### Errors

When no provider is registered for the model or any of its fallbacks, the router returns `heimdall.ErrUnsupportedProvider`. `router.Stream` returns `heimdall.ErrNoChunkHandler` when the handler is nil. Other errors come from the last model that was tried.

```go
res, err := router.Complete(ctx, req)
switch {
case errors.Is(err, heimdall.ErrUnsupportedProvider):
	// Register a provider for one of the request's models.
case err != nil:
	// The request failed on every model.
}
```

Every response also carries `RequestLog`, which lists each model and key that was tried.

## Supported models

The tables list current models. Models with an announced shutdown date and retired models remain in the `models` package and are marked in their doc comments. The [package reference](https://pkg.go.dev/github.com/JazzJackrabbit/heimdall-ng/models) lists every type.

<details>
<summary>OpenAI</summary>

| Type | Model ID |
| --- | --- |
| `GPT6Astra` | gpt-6-astra |
| `GPT6Sol` | gpt-6-sol |
| `GPT6Luna` | gpt-6-luna |
| `GPT56Sol` | gpt-5.6-sol |
| `GPT56Terra` | gpt-5.6-terra |
| `GPT56Luna` | gpt-5.6-luna |
| `GPT55` | gpt-5.5 |
| `GPT54` | gpt-5.4 |
| `GPT54Mini` | gpt-5.4-mini |
| `GPT54Nano` | gpt-5.4-nano |
| `GPT53Codex` | gpt-5.3-codex |
| `GPT52` | gpt-5.2 |
| `GPT51` | gpt-5.1 |
| `GPT51Chat` | gpt-5.1-chat-latest |
| `GPT5Chat` | gpt-5-chat-latest |
| `GPT41` | gpt-4.1-2025-04-14 |
| `GPT41Mini` | gpt-4.1-mini-2025-04-14 |
| `GPT4O` | gpt-4o-2024-11-20 |
| `GPT4OMini` | gpt-4o-mini-2024-07-18 |
| `GPTImage2` | gpt-image-2 |

</details>

<details>
<summary>Anthropic</summary>

| Type | Model ID |
| --- | --- |
| `ClaudeFable51` | claude-fable-5-1 |
| `Claude55Opus` | claude-opus-5-5 |
| `ClaudeFable5` | claude-fable-5 |
| `Claude5Opus` | claude-opus-5 |
| `Claude5Sonnet` | claude-sonnet-5 |
| `Claude48Opus` | claude-opus-4-8 |
| `Claude47Opus` | claude-opus-4-7 |
| `Claude46Opus` | claude-opus-4-6 |
| `Claude46Sonnet` | claude-sonnet-4-6 |
| `Claude45Opus` | claude-opus-4-5-20251101 |
| `Claude45Sonnet` | claude-sonnet-4-5-20250929 |
| `Claude45Haiku` | claude-haiku-4-5 |

Claude Fable models require 30-day data retention on the Anthropic organization.

</details>

<details>
<summary>Google Gemini and Vertex AI</summary>

Each Gemini type has a Vertex AI counterpart with the `Vertex` prefix, for example `VertexGemini38Flash`. `Gemini25ProPreview` and `Gemini25FlashPreview` keep their original names for compatibility. Their Vertex AI types are `VertexGemini25Pro` and `VertexGemini25Flash`.

| Type | Model ID |
| --- | --- |
| `Gemini38Flash` | gemini-3.8-flash |
| `Gemini37Flash` | gemini-3.7-flash |
| `Gemini36Flash` | gemini-3.6-flash |
| `Gemini35Flash` | gemini-3.5-flash |
| `Gemini35FlashLite` | gemini-3.5-flash-lite |
| `Gemini31ProPreview` | gemini-3.1-pro-preview |
| `Gemini31FlashImage` | gemini-3.1-flash-image |
| `Gemini3FlashPreview` | gemini-3-flash-preview |
| `Gemini3ProImage` | gemini-3-pro-image |
| `Gemini25ProPreview` | gemini-2.5-pro |
| `Gemini25FlashPreview` | gemini-2.5-flash |
| `Gemini25FlashLite` | gemini-2.5-flash-lite |

</details>

<details>
<summary>xAI Grok, Perplexity and OpenRouter</summary>

| Type | Model ID |
| --- | --- |
| `Grok46` | grok-4.6 |
| `Grok45` | grok-4.5 |
| `Grok43` | grok-4.3 |
| `SonarReasoningPro` | sonar-reasoning-pro |
| `SonarPro` | sonar-pro |
| `Sonar` | sonar |
| `OpenRouterModel` | any OpenRouter model, set in `ModelName` |

</details>

## Development

```bash
go test ./...
golangci-lint run
```

Request-building and response-parsing tests run offline against local test servers. Tests that call a provider's API are skipped unless its credentials are set in `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `GOOGLE_API_KEY`, `GROK_API_KEY` or `PERPLEXITY_API_KEY`. Vertex AI tests use `VERTEX_PROJECT_ID` and `VERTEX_LOCATION`. The [justfile](justfile) has per-provider test recipes and loads these variables from a `.env` file.

## License

BSD 3-Clause. See [LICENSE](LICENSE).
