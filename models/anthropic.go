package models

const AnthropicProvider = "anthropic"

const (
	// Deprecated: Claude 3 Opus was retired on January 5, 2026. Requests to this model will fail.
	AnthropicClaude3OpusAlias = "claude-3-opus-20240229"
	// Deprecated: Claude 3.5 Sonnet was retired on October 28, 2025. Requests to this model will fail.
	AnthropicClaude35SonnetAlias = "claude-3-5-sonnet-20241022"
	// Deprecated: Claude 3.5 Haiku was retired on February 19, 2026. Requests to this model will fail.
	// Use Claude45Haiku (claude-haiku-4-5) as a replacement.
	AnthropicClaude35HaikuAlias = "claude-3-5-haiku-20241022"
	// Deprecated: Claude 3.7 Sonnet was retired on February 19, 2026. Requests to this model will fail.
	// Use Claude46Sonnet (claude-sonnet-4-6) as a replacement.
	AnthropicClaude37SonnetAlias = "claude-3-7-sonnet-20250219"
	AnthropicClaude4SonnetAlias  = "claude-sonnet-4-20250514"
	AnthropicClaude4OpusAlias    = "claude-opus-4-20250514"
	AnthropicClaude45HaikuAlias  = "claude-haiku-4-5"
	AnthropicClaude45SonnetAlias = "claude-sonnet-4-5-20250929"
	AnthropicClaude45OpusAlias   = "claude-opus-4-5-20251101"
	AnthropicClaude46OpusAlias   = "claude-opus-4-6"
	AnthropicClaude46SonnetAlias = "claude-sonnet-4-6"
	AnthropicClaude47OpusAlias   = "claude-opus-4-7"
	AnthropicClaude48OpusAlias   = "claude-opus-4-8"
	AnthropicClaude5OpusAlias    = "claude-opus-5"
	AnthropicClaude5SonnetAlias  = "claude-sonnet-5"
	AnthropicClaudeFable5Alias   = "claude-fable-5"
)

type (
	AnthropicImageType string
	AnthropicPdf       string
)

const (
	AnthropicImageJpeg AnthropicImageType = "image/jpeg"
	AnthropicImagePng  AnthropicImageType = "image/png"
	AnthropicImageGif  AnthropicImageType = "image/gif"
	AnthropicImageWebp AnthropicImageType = "image/webp"
)

type Claude3Opus struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude3Opus) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000015
}

func (c Claude3Opus) GetName() string {
	return AnthropicClaude3OpusAlias
}

func (c Claude3Opus) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude3Opus)

type Claude35Sonnet struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude35Sonnet) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000003
}

func (c Claude35Sonnet) GetName() string {
	return AnthropicClaude35SonnetAlias
}

func (c Claude35Sonnet) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude35Sonnet)

type Claude35Haiku struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude35Haiku) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.0000008
}

func (c Claude35Haiku) GetName() string {
	return AnthropicClaude35HaikuAlias
}

func (c Claude35Haiku) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude35Haiku)

type Claude37Sonnet struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude37Sonnet) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000003
}

func (c Claude37Sonnet) GetName() string {
	return AnthropicClaude37SonnetAlias
}

func (c Claude37Sonnet) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude37Sonnet)

type Claude4Sonnet struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude4Sonnet) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000003
}

func (c Claude4Sonnet) GetName() string {
	return AnthropicClaude4SonnetAlias
}

func (c Claude4Sonnet) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude4Sonnet)

type Claude4Opus struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude4Opus) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000015
}

func (c Claude4Opus) GetName() string {
	return AnthropicClaude4OpusAlias
}

func (c Claude4Opus) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude4Opus)

type Claude45Haiku struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude45Haiku) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000001
}

func (c Claude45Haiku) GetName() string {
	return AnthropicClaude45HaikuAlias
}

func (c Claude45Haiku) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude45Haiku)

type Claude45Sonnet struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude45Sonnet) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000003
}

func (c Claude45Sonnet) GetName() string {
	return AnthropicClaude45SonnetAlias
}

func (c Claude45Sonnet) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude45Sonnet)

type Claude45Opus struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
}

func (c Claude45Opus) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000005
}

func (c Claude45Opus) GetInputCostPer1M() float64 {
	return 5.0
}

func (c Claude45Opus) GetOutputCostPer1M() float64 {
	return 25.0
}

func (c Claude45Opus) GetName() string {
	return AnthropicClaude45OpusAlias
}

func (c Claude45Opus) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude45Opus)
var _ CostBreakdown = new(Claude45Opus)

type Claude46Opus struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
	// ExtendedContext enables the 1M token context window (beta).
	// Requires the context-1m-2025-08-07 beta header.
	ExtendedContext bool
	// MaxOutputTokens sets the maximum output tokens (up to 128K).
	// Defaults to 4096 when zero.
	MaxOutputTokens int
}

func (c Claude46Opus) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000005
}

func (c Claude46Opus) GetInputCostPer1M() float64 {
	return 5.0
}

func (c Claude46Opus) GetOutputCostPer1M() float64 {
	return 25.0
}

func (c Claude46Opus) GetName() string {
	return AnthropicClaude46OpusAlias
}

func (c Claude46Opus) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude46Opus)
var _ CostBreakdown = new(Claude46Opus)

type Claude46Sonnet struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
	// ExtendedContext enables the 1M token context window (beta).
	ExtendedContext bool
	// MaxOutputTokens sets the maximum output tokens (up to 64K).
	// Defaults to 4096 when zero.
	MaxOutputTokens int
}

func (c Claude46Sonnet) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000003
}

func (c Claude46Sonnet) GetInputCostPer1M() float64 {
	return 3.0
}

func (c Claude46Sonnet) GetOutputCostPer1M() float64 {
	return 15.0
}

func (c Claude46Sonnet) GetName() string {
	return AnthropicClaude46SonnetAlias
}

func (c Claude46Sonnet) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude46Sonnet)
var _ CostBreakdown = new(Claude46Sonnet)

// Claude47Opus does not accept temperature or top_p; leave those request fields unset.
type Claude47Opus struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
	// MaxOutputTokens sets the maximum output tokens (up to 128K).
	// Defaults to 4096 when zero.
	MaxOutputTokens int
}

func (c Claude47Opus) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000005
}

func (c Claude47Opus) GetInputCostPer1M() float64 {
	return 5.0
}

func (c Claude47Opus) GetOutputCostPer1M() float64 {
	return 25.0
}

func (c Claude47Opus) GetName() string {
	return AnthropicClaude47OpusAlias
}

func (c Claude47Opus) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude47Opus)
var _ CostBreakdown = new(Claude47Opus)

// Claude48Opus does not accept temperature or top_p; leave those request fields unset.
type Claude48Opus struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
	// MaxOutputTokens sets the maximum output tokens (up to 128K).
	// Defaults to 4096 when zero.
	MaxOutputTokens int
}

func (c Claude48Opus) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000005
}

func (c Claude48Opus) GetInputCostPer1M() float64 {
	return 5.0
}

func (c Claude48Opus) GetOutputCostPer1M() float64 {
	return 25.0
}

func (c Claude48Opus) GetName() string {
	return AnthropicClaude48OpusAlias
}

func (c Claude48Opus) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude48Opus)
var _ CostBreakdown = new(Claude48Opus)

// Claude5Opus does not accept temperature or top_p; leave those request fields unset.
type Claude5Opus struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
	// MaxOutputTokens sets the maximum output tokens (up to 128K).
	// Defaults to 4096 when zero.
	MaxOutputTokens int
}

func (c Claude5Opus) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000005
}

func (c Claude5Opus) GetInputCostPer1M() float64 {
	return 5.0
}

func (c Claude5Opus) GetOutputCostPer1M() float64 {
	return 25.0
}

func (c Claude5Opus) GetName() string {
	return AnthropicClaude5OpusAlias
}

func (c Claude5Opus) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude5Opus)
var _ CostBreakdown = new(Claude5Opus)

// Claude5Sonnet does not accept non-default temperature or top_p; leave those request fields unset.
type Claude5Sonnet struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
	// MaxOutputTokens sets the maximum output tokens (up to 128K).
	// Defaults to 4096 when zero.
	MaxOutputTokens int
}

func (c Claude5Sonnet) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000002
}

func (c Claude5Sonnet) GetInputCostPer1M() float64 {
	return 2.0
}

func (c Claude5Sonnet) GetOutputCostPer1M() float64 {
	return 10.0
}

func (c Claude5Sonnet) GetName() string {
	return AnthropicClaude5SonnetAlias
}

func (c Claude5Sonnet) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(Claude5Sonnet)
var _ CostBreakdown = new(Claude5Sonnet)

// ClaudeFable5 requires the organization to have 30-day data retention;
// requests from zero-data-retention organizations are rejected by Anthropic.
// It does not accept temperature or top_p; leave those request fields unset.
type ClaudeFable5 struct {
	ImageFile        map[AnthropicImageType]string
	PdfFiles         []AnthropicPdf
	StructuredOutput map[string]any
	// MaxOutputTokens sets the maximum output tokens (up to 128K).
	// Defaults to 4096 when zero.
	MaxOutputTokens int
}

func (c ClaudeFable5) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00001
}

func (c ClaudeFable5) GetInputCostPer1M() float64 {
	return 10.0
}

func (c ClaudeFable5) GetOutputCostPer1M() float64 {
	return 50.0
}

func (c ClaudeFable5) GetName() string {
	return AnthropicClaudeFable5Alias
}

func (c ClaudeFable5) GetProvider() string {
	return AnthropicProvider
}

var _ Model = new(ClaudeFable5)
var _ CostBreakdown = new(ClaudeFable5)
