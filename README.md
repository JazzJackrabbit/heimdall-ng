# Heimdall

Heimdall is a Go library for making LLM requests more consistent by providing model fallbacks and API key rotation. It acts as a router between your application and various LLM providers. This is a maintained fork of [flyx-ai/heimdall](https://github.com/flyx-ai/heimdall).

[![Go Reference](https://pkg.go.dev/badge/github.com/JazzJackrabbit/heimdall-ng.svg)](https://pkg.go.dev/github.com/JazzJackrabbit/heimdall-ng)
[![Release](https://img.shields.io/github/v/release/JazzJackrabbit/heimdall-ng)](https://github.com/JazzJackrabbit/heimdall-ng/releases)
[![License: BSD-3](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](LICENSE)

## Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Provider Setup](#provider-setup)
- [Multi-turn Conversations](#multi-turn-conversations)
- [Working with PDF Files](#working-with-pdf-files)
- [Streaming Responses](#streaming-responses)
- [Structured Output](#structured-output)
- [Advanced Router Configuration](#advanced-router-configuration)
- [Working with Images](#working-with-images)
- [Error Handling](#error-handling)
- [Supported Models](#supported-models)
- [License](#license)

## Features

- **Provider Abstraction**: Unified interface for multiple LLM providers (OpenAI, Anthropic, Google/Gemini, Grok, VertexAI, OpenRouter, Perplexity)
- **Model Fallbacks**: Configurable fallback models if the primary model fails
- **API Key Rotation**: Pass multiple keys per provider; requests rotate to a backup key on failure
- **Streaming Support**: Fully supports streaming responses for real-time applications
- **Multimodal Inputs**: Support for PDFs, images, and other file types
- **Structured Output**: Format responses as JSON using provider-specific schema formats
- **Request Logging**: Built-in request logging for monitoring and debugging

## Installation

```bash
go get github.com/JazzJackrabbit/heimdall-ng
```

## Quick Start

### Basic Usage

Here's a simple example of how to use Heimdall with OpenAI:

```go
package main

import (
	"context"
	"os"
	"time"

	"github.com/JazzJackrabbit/heimdall-ng"
	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/JazzJackrabbit/heimdall-ng/providers"
	"github.com/JazzJackrabbit/heimdall-ng/request"
)

func main() {
	ctx := context.Background()
	
	// Create a provider with your API key
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIProvider := providers.NewOpenAI([]string{openAIAPIKey})
	
	// Setup the router with a timeout and the provider
	timeout := 30 * time.Second
	router := heimdall.New(timeout, []heimdall.LLMProvider{openAIProvider})
	
	// Create a completion request
	req := request.Completion{
		Model: models.GPT4O{},
		SystemMessage: "You are a helpful assistant.",
		UserMessage:   "What's the capital of France?",
		Fallback: []models.Model{
			models.GPT4OMini{},
		},
		Temperature: 0.7,
		TopP:        1.0,
		Tags: map[string]string{
			"env":  "production",
			"type": "geography",
		},
	}
	
	// Get a completion
	response, err := router.Complete(ctx, req)
	if err != nil {
		panic(err)
	}
	
	// Use the response
	println(response.Content)
}
```

> **Two ways to call.** The recommended entry point is the **router** (`router.Complete` / `router.Stream`), which adds model fallbacks on top of your providers. Every provider also exposes the lower-level `CompleteResponse` / `StreamResponse` methods directly — several examples below use these to keep the snippet self-contained, but you can route any of them through the router the same way.

## Provider Setup

Heimdall supports multiple LLM providers. Here's how to set them up:

### OpenAI

```go
// Single API key
openAIProvider := providers.NewOpenAI([]string{"your-api-key"})

// Multiple API keys for load balancing
openAIProvider := providers.NewOpenAI([]string{"key1", "key2", "key3"})
```

### Anthropic

```go
anthropicProvider := providers.NewAnthropic([]string{"your-api-key"})
```

### Google/Gemini

```go
googleProvider := providers.NewGoogle([]string{"your-api-key"})
```

Heimdall also supports caching tokens for Google provider to reduce latency and improve performance for repeated requests. 

```go
func main() {
	ctx := context.Background()

	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	g := providers.NewGoogle([]string{googleApiKey})
    
	key, err := g.CacheContent(
		ctx,
		models.Gemini25FlashPreview{}.GetName(),
		providers.CacheContentPayload{
			// Map of MIME type -> file URI (e.g. a Cloud Storage URI)
			FileData: map[string]string{
				"application/pdf": "gs://my-bucket/contract.pdf",
			},
		},
		"You are a contract analysis assistant.",
		10*time.Minute,
	)
	if err != nil {
		panic(err)
	}

	// Reuse the returned id in later requests to skip re-sending the cached content
	fmt.Println("cached content id:", key)
}

```

### Perplexity

```go
perplexityProvider := providers.NewPerplexity([]string{"your-api-key"})
```

### VertexAI

```go
vertexAIProvider := providers.NewVertexAI([]string{"your-api-key"})
```

### Grok

```go
grokProvider := providers.NewGrok([]string{"your-api-key"})
```

### OpenRouter

```go
openRouterProvider := providers.NewOpenRouter([]string{"your-api-key"})
```

## Multi-turn Conversations

Pass prior turns through the `History` field to give the model conversational context. Each `request.Message` has a `Role` (`"user"` or `"assistant"`) and `Content`:

```go
req := request.Completion{
	Model:         models.GPT4O{},
	SystemMessage: "You are a helpful travel assistant.",
	History: []request.Message{
		{Role: "user", Content: "I'm planning a trip to Japan."},
		{Role: "assistant", Content: "Great! When are you planning to go?"},
	},
	UserMessage: "Sometime in April — what should I know?",
}

response, err := router.Complete(ctx, req)
if err != nil {
	panic(err)
}

fmt.Println(response.Content)
```

## Working with PDF Files

### OpenAI with PDF Input

```go
func main() {
	ctx := context.Background()
	
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	
	// Read and encode the PDF file
	fileBytes, err := os.ReadFile("document.pdf")
	if err != nil {
		panic(err)
	}
	
	encodedString := base64.StdEncoding.EncodeToString(fileBytes)
	dataURL := "data:application/pdf;base64," + encodedString
	
	openAIProvider := providers.NewOpenAI([]string{openAIAPIKey})
	res, err := openAIProvider.CompleteResponse(
		ctx,
		request.Completion{
			Model: models.GPT4O{
				PdfFile: map[string]string{
					"document.pdf": dataURL,
				},
			},
			UserMessage: "Summarize the contents of this PDF",
		},
		http.Client{},
		nil,
	)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res.Content)
}
```

### Google/Gemini with PDF Input

```go
func main() {
	ctx := context.Background()
	
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	
	// Read and encode a PDF file
	fileBytes, err := os.ReadFile("document.pdf")
	if err != nil {
		panic(err)
	}
	encodedPdf := base64.StdEncoding.EncodeToString(fileBytes)
	
	// Create a Google provider
	googleProvider := providers.NewGoogle([]string{googleAPIKey})
	
	// Example with base64 encoded PDF
	res, err := googleProvider.CompleteResponse(
		ctx,
		request.Completion{
			Model: models.Gemini25FlashPreview{
				PdfFiles: []models.GooglePdf{
					models.GooglePdf(encodedPdf), // Base64 encoded PDF
				},
			},
			SystemMessage: "You are a helpful assistant that analyzes documents.",
			UserMessage:   "Summarize the main points from this PDF document",
			Temperature:   0,
			TopP:          0,
			Tags:          map[string]string{},
		},
		http.Client{},
		nil,
	)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res.Content)
}
```

### Anthropic with PDF Input

```go
func main() {
	ctx := context.Background()
	
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	
	// Read and encode a PDF file
	fileBytes, err := os.ReadFile("document.pdf")
	if err != nil {
		panic(err)
	}
	encodedPdf := base64.StdEncoding.EncodeToString(fileBytes)
	
	// Create an Anthropic provider
	anthropicProvider := providers.NewAnthropic([]string{anthropicAPIKey})
	
	// Request with PDF
	res, err := anthropicProvider.CompleteResponse(
		ctx,
		request.Completion{
			Model: models.Claude46Sonnet{
				PdfFiles: []models.AnthropicPdf{
					models.AnthropicPdf(encodedPdf),
				},
			},
			SystemMessage: "You are a helpful document analyst.",
			UserMessage:   "Analyze this PDF and provide key insights",
			Temperature:   0,
			TopP:          0,
			Tags:          map[string]string{},
		},
		http.Client{},
		nil,
	)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res.Content)
}
```

## Streaming Responses

For streaming responses, use the `Stream` method:

```go
func main() {
	ctx := context.Background()
	
	// Create provider and router
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIProvider := providers.NewOpenAI([]string{openAIAPIKey})
	
	timeout := 30 * time.Second
	router := heimdall.New(timeout, []heimdall.LLMProvider{openAIProvider})
	
	// Create request
	req := request.Completion{
		Model: models.GPT4O{},
		SystemMessage: "You are a helpful assistant.",
		UserMessage:   "Write a short story about a space explorer",
		Temperature: 0.7,
		Tags: map[string]string{
			"env":  "production",
			"type": "creative",
		},
	}
	
	// Handle streaming chunks
	chunkHandler := func(chunk string) error {
		fmt.Print(chunk)
		return nil
	}
	
	// Get streaming response
	_, err := router.Stream(ctx, req, chunkHandler)
	if err != nil {
		panic(err)
	}
}
```

## Structured Output

You can request structured output from supported models:

### OpenAI Structured Output

```go
var schema = map[string]any{
	"name": "financial_analysis",
	"schema": map[string]any{
		"type": "object",
		"properties": map[string]any{
			"summary":     map[string]any{"type": "string"},
			"marketCap":   map[string]any{"type": "number"},
			"advantages":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"risks":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"recommendation": map[string]any{"type": "string", "enum": []string{"buy", "hold", "sell"}},
		},
		"required": []string{"summary", "recommendation"},
	},
}

func main() {
	ctx := context.Background()
	
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	
	openAIProvider := providers.NewOpenAI([]string{openAIAPIKey})
	res, err := openAIProvider.CompleteResponse(
		ctx,
		request.Completion{
			Model: models.GPT4O{
				StructuredOutput: schema,
			},
			UserMessage: "Create a financial analysis of Nvidia at its current valuation",
		},
		http.Client{},
		nil,
	)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res.Content)
}
```

### Google/Gemini Structured Output

```go
var schemaGoogle = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"summary":     map[string]any{"type": "string"},
		"marketCap":   map[string]any{"type": "number"},
		"advantages":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		"risks":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		"recommendation": map[string]any{"type": "string", "enum": []string{"buy", "hold", "sell"}},
	},
	"required": []string{"summary", "recommendation"},
}

func main() {
	ctx := context.Background()
	
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	
	googleProvider := providers.NewGoogle([]string{googleAPIKey})
	res, err := googleProvider.CompleteResponse(
		ctx,
		request.Completion{
			Model: models.Gemini25FlashPreview{
				StructuredOutput: schemaGoogle,
			},
			SystemMessage: "You are a financial analyst.",
			UserMessage:   "Create a financial analysis of Nvidia at its current valuation",
			Temperature:   0,
			TopP:          0,
			Tags:          map[string]string{},
		},
		http.Client{},
		nil,
	)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res.Content)
}
```

## Advanced Router Configuration

You can configure Heimdall with multiple providers and fallback options:

```go
func main() {
	ctx := context.Background()
	
	// Configure providers
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	
	openAIProvider := providers.NewOpenAI([]string{openAIAPIKey})
	googleProvider := providers.NewGoogle([]string{googleAPIKey})
	anthropicProvider := providers.NewAnthropic([]string{anthropicAPIKey})
	
	// Create router with all providers
	timeout := 60 * time.Second
	router := heimdall.New(timeout, []heimdall.LLMProvider{
		openAIProvider,
		googleProvider,
		anthropicProvider,
	})
	
	// Create request with primary model and fallbacks
	req := request.Completion{
		// Primary model - will be tried first
		Model: models.GPT4O{},
		
		// Fallbacks - will be tried in order if primary fails
		Fallback: []models.Model{
			models.Claude46Sonnet{},
			models.Gemini25FlashPreview{},
		},
		
		SystemMessage: "You are a helpful assistant.",
		UserMessage:   "Explain quantum computing in simple terms",
		Temperature: 0.3,
		TopP:        1.0,
		Tags: map[string]string{
			"env":     "production",
			"type":    "science",
			"purpose": "education",
		},
	}
	
	// Get completion, which will try fallbacks if needed
	response, err := router.Complete(ctx, req)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(response.Content)
	
	// Access request logs for debugging
	fmt.Printf("Request logs: %+v\n", response.RequestLog)
}
```

## Working with Images

### OpenAI with Image Input

```go
func main() {
	ctx := context.Background()
	
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	
	// Read and encode the image file
	fileBytes, err := os.ReadFile("image.jpg")
	if err != nil {
		panic(err)
	}
	
	encodedString := base64.StdEncoding.EncodeToString(fileBytes)
	dataURL := "data:image/jpeg;base64," + encodedString
	
	openAIProvider := providers.NewOpenAI([]string{openAIAPIKey})
	res, err := openAIProvider.CompleteResponse(
		ctx,
		request.Completion{
			Model: models.GPT4O{
				ImageFile: []models.OpenaiImagePayload{
					{
						Url:    dataURL,
						Detail: "high",
					},
				},
			},
			UserMessage: "Describe what you see in this image",
		},
		http.Client{},
		nil,
	)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res.Content)
}
```

### Anthropic with Image Input

```go
func main() {
	ctx := context.Background()
	
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	
	// Read and encode the image file
	fileBytes, err := os.ReadFile("image.jpg")
	if err != nil {
		panic(err)
	}
	
	encodedString := base64.StdEncoding.EncodeToString(fileBytes)
	
	anthropicProvider := providers.NewAnthropic([]string{anthropicAPIKey})
	res, err := anthropicProvider.CompleteResponse(
		ctx,
		request.Completion{
			Model: models.Claude46Sonnet{
				ImageFile: map[models.AnthropicImageType]string{
					models.AnthropicImageJpeg: encodedString,
				},
			},
			SystemMessage: "You are a helpful visual assistant.",
			UserMessage:   "Describe what you see in this image",
			Temperature:   0,
			TopP:          0,
			Tags:          map[string]string{},
		},
		http.Client{},
		nil,
	)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res.Content)
}
```

## Error Handling

Heimdall provides comprehensive error handling. Here's an example of how to handle errors:

```go
response, err := router.Complete(ctx, req)
if err != nil {
    switch {
    case errors.Is(err, heimdall.ErrNoChunkHandler):
        // Handle missing chunk handler for streaming
        fmt.Println("Streaming requires a chunk handler")
    case errors.Is(err, heimdall.ErrUnsupportedProvider):
        // Handle case where provider for model is not registered
        fmt.Println("No provider registered for this model")
    default:
        // Handle other errors
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Supported Models

Heimdall supports various models from different providers:

### OpenAI Models
- GPT-5.6 Sol (gpt-5.6-sol)
- GPT-5.6 Terra (gpt-5.6-terra)
- GPT-5.6 Luna (gpt-5.6-luna)
- GPT-5.5 (gpt-5.5)
- GPT-5.4 (gpt-5.4)
- GPT-5.4 Mini (gpt-5.4-mini)
- GPT-5.4 Nano (gpt-5.4-nano)
- GPT-5.3 Codex (gpt-5.3-codex)
- GPT-4 (gpt-4-0613)
- GPT-4 Turbo (gpt-4-turbo)
- GPT-4o (gpt-4o-2024-11-20)
- GPT-4o Mini (gpt-4o-mini-2024-07-18)
- GPT-4.1 (gpt-4.1-2025-04-14)
- GPT-4.1 Mini (gpt-4.1-mini-2025-04-14)
- GPT-4.1 Nano (gpt-4.1-nano-2025-04-14)
- GPT-5 (gpt-5-2025-08-07)
- GPT-5 Mini (gpt-5-mini-2025-08-07)
- GPT-5 Nano (gpt-5-nano-2025-08-07)
- GPT-5 Chat (gpt-5-chat-latest)
- GPT-5.1 (gpt-5.1)
- GPT-5.1 Chat (gpt-5.1-chat-latest)
- GPT-5.1 Codex (gpt-5.1-codex)
- GPT-5.1 Codex Mini (gpt-5.1-codex-mini)
- GPT-5.2 (gpt-5.2)
- O1 (o1-2024-12-17)
- O3 Mini (o3-mini-2025-01-31)
- O3 (o3)
- O4 Mini (o4-mini)
- GPT Image 2 (gpt-image-2), same parameters as gpt-image-1 plus arbitrary sizes up to 3840x2160
- GPT Image (gpt-image-1) — pass `ImageFile` to edit reference images instead of generating from the prompt alone

### Anthropic Models
- Claude Fable 5 (claude-fable-5) — requires 30-day data retention on the organization
- Claude Opus 5 (claude-opus-5)
- Claude Sonnet 5 (claude-sonnet-5)
- Claude 4.8 Opus (claude-opus-4-8)
- Claude 4.7 Opus (claude-opus-4-7)
- Claude 4.6 Opus (claude-opus-4-6)
- Claude 4.6 Sonnet (claude-sonnet-4-6)
- Claude 4.5 Opus (claude-opus-4-5-20251101)
- Claude 4.5 Sonnet (claude-sonnet-4-5-20250929)
- Claude 4.5 Haiku (claude-haiku-4-5)

### Google/Gemini Models
- Gemini 3.7 Flash (gemini-3.7-flash)
- Gemini 3.6 Flash (gemini-3.6-flash)
- Gemini 3.5 Flash (gemini-3.5-flash)
- Gemini 3.5 Flash Lite (gemini-3.5-flash-lite)
- Gemini 3.1 Pro Preview (gemini-3.1-pro-preview)
- Gemini 3.1 Flash Lite (gemini-3.1-flash-lite)
- Gemini 3.1 Flash Image (gemini-3.1-flash-image)
- Gemini 3 Flash Preview (gemini-3-flash-preview)
- Gemini 3 Pro Image (gemini-3-pro-image)
- Gemini 2.5 Pro (gemini-2.5-pro)
- Gemini 2.5 Flash (gemini-2.5-flash)
- Gemini 2.5 Flash Lite (gemini-2.5-flash-lite)
- Gemini 2.5 Flash Image (gemini-2.5-flash-image)

### VertexAI Models
Gemini models served through Google Cloud Vertex AI (use the `VertexGemini*` model types):
- Gemini 3.7 Flash (gemini-3.7-flash)
- Gemini 3.6 Flash (gemini-3.6-flash)
- Gemini 3.5 Flash (gemini-3.5-flash)
- Gemini 3.5 Flash Lite (gemini-3.5-flash-lite)
- Gemini 3.1 Pro Preview (gemini-3.1-pro-preview)
- Gemini 3.1 Flash Lite (gemini-3.1-flash-lite)
- Gemini 3.1 Flash Image (gemini-3.1-flash-image)
- Gemini 3 Flash Preview (gemini-3-flash-preview)
- Gemini 3 Pro Image (gemini-3-pro-image)
- Gemini 2.5 Pro (gemini-2.5-pro)
- Gemini 2.5 Flash (gemini-2.5-flash)
- Gemini 2.5 Flash Lite (gemini-2.5-flash-lite)
- Gemini 2.5 Flash Image (gemini-2.5-flash-image)

### Grok Models
- Grok 4.6 (grok-4.6)
- Grok 4.5 (grok-4.5)
- Grok 4.3 (grok-4.3)

### Perplexity Models
- Sonar Reasoning Pro (sonar-reasoning-pro)
- Sonar Pro (sonar-pro)
- Sonar (sonar)

### OpenRouter
- Supports any model available on OpenRouter via dynamic model names

## License

BSD-3-Clause — see [LICENSE](LICENSE). Based on [flyx-ai/heimdall](https://github.com/flyx-ai/heimdall).
