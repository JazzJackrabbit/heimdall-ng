package models

const OpenaiProvider = "openai"

const (
	O3MiniAlias         = "o3-mini-2025-01-31"
	GPT4OAlias          = "gpt-4o-2024-11-20"
	GPT4OMiniAlias      = "gpt-4o-mini-2024-07-18"
	O1Alias             = "o1-2024-12-17"
	GPT4Alias           = "gpt-4-0613"
	GPT4TurboAlias      = "gpt-4-turbo"
	GPT41Alias          = "gpt-4.1-2025-04-14"
	GPT41MiniAlias      = "gpt-4.1-mini-2025-04-14"
	GPT41NanoAlias      = "gpt-4.1-nano-2025-04-14"
	GPT5Alias           = "gpt-5-2025-08-07"
	GPT5MiniAlias       = "gpt-5-mini-2025-08-07"
	GPT5NanoAlias       = "gpt-5-nano-2025-08-07"
	GPT5ChatAlias       = "gpt-5-chat-latest"
	GPT51Alias          = "gpt-5.1"
	GPT51ChatAlias      = "gpt-5.1-chat-latest"
	GPT51CodexAlias     = "gpt-5.1-codex"
	GPT51CodexMiniAlias = "gpt-5.1-codex-mini"
	GPT52Alias          = "gpt-5.2"
	GPT53CodexAlias     = "gpt-5.3-codex"
	GPT54Alias          = "gpt-5.4"
	GPT54MiniAlias      = "gpt-5.4-mini"
	GPT54NanoAlias      = "gpt-5.4-nano"
	GPT55Alias          = "gpt-5.5"
	GPT56SolAlias       = "gpt-5.6-sol"
	GPT56TerraAlias     = "gpt-5.6-terra"
	GPT56LunaAlias      = "gpt-5.6-luna"
	GPT6AstraAlias      = "gpt-6-astra"
	GPT6SolAlias        = "gpt-6-sol"
	GPT6LunaAlias       = "gpt-6-luna"
	O3Alias             = "o3"
	O4MiniAlias         = "o4-mini"
)

type OpenaiImagePayload struct {
	// Url can be ether that, an url or a base64 encoding of the image .
	// If using base64, it must follow this format: data:image/jpeg;base64,{base64_image}
	Url string
	// Detail determines the level detail to use when processing and understanding the image. Can be either: high, low or auto. If nothing is specified, it will default to auto.
	Detail string
}

type GPT41 struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any
	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT41) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000200
}

func (GPT41) GetName() string {
	return GPT41Alias
}

func (GPT41) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT41)

type GPT41Mini struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any
	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT41Mini) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000040
}

func (GPT41Mini) GetName() string {
	return GPT41MiniAlias
}

func (GPT41Mini) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT41Mini)

// Scheduled for shutdown by OpenAI on October 23, 2026. Use GPT56Luna (gpt-5.6-luna) as a replacement.
type GPT41Nano struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any
	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT41Nano) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000010
}

func (GPT41Nano) GetName() string {
	return GPT41NanoAlias
}

func (GPT41Nano) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT41Nano)

// Scheduled for shutdown by OpenAI on October 23, 2026. Use GPT56Sol (gpt-5.6-sol) as a replacement.
type O3Mini struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any
	// Note: O3Mini does not support vision/images in the API as of 2025
}

func (o O3Mini) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000110
}

func (o O3Mini) GetName() string {
	return O3MiniAlias
}

func (o O3Mini) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(O3Mini)

// Scheduled for shutdown by OpenAI on October 23, 2026. Use GPT56Sol (gpt-5.6-sol) as a replacement.
type O1 struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any

	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (o O1) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00001500
}

func (o O1) GetName() string {
	return O1Alias
}

func (o O1) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(O1)

// Scheduled for shutdown by OpenAI on October 23, 2026. Use GPT56Sol (gpt-5.6-sol) as a replacement.
type GPT4 struct {
	// Note: GPT-4 (gpt-4-0613) does not support vision/images or PDFs
	// This is a text-only model released before vision capabilities were added
}

func (g GPT4) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00006000
}

func (g GPT4) GetName() string {
	return GPT4Alias
}

func (g GPT4) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT4)

// Scheduled for shutdown by OpenAI on October 23, 2026. Use GPT56Sol (gpt-5.6-sol) as a replacement.
type GPT4Turbo struct {
	// ImageFile enables vision for the request
	// Note: GPT-4 Turbo supports vision but NOT direct PDF input
	// PDF support was only added to newer models (GPT-4o, GPT-4.1, o1) in March 2025
	ImageFile []OpenaiImagePayload
}

func (g GPT4Turbo) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00001000
}

func (g GPT4Turbo) GetName() string {
	return GPT4TurboAlias
}

func (g GPT4Turbo) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT4Turbo)

type GPT4O struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any

	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT4O) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000250
}

func (g GPT4O) GetName() string {
	return GPT4OAlias
}

func (g GPT4O) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT4O)

type (
	GPT4OMini struct {
		// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
		//
		//  var schema = map[string]any{
		//  	"name": "navidia_valuation",
		//  	"schema": map[string]any{
		//  		"type": "object",
		//  		"properties": map[string]any{
		//  			"final_answer": map[string]any{"type": "string"},
		//  			"valuation": map[string]any{
		//  				"type": "number",
		//  			},
		//  		},
		//  	},
		//  }
		StructuredOutput map[string]any

		// PdfFile let's you include a PDF file in your request to the LLM.
		// The expected format:
		//
		// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
		// Only provide a pdf file or an image file, not both.
		PdfFile map[string]string
		// ImageFile enables vision for the request
		ImageFile []OpenaiImagePayload
	}
)

func (g GPT4OMini) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000015
}

func (g GPT4OMini) GetName() string {
	return GPT4OMiniAlias
}

func (g GPT4OMini) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT4OMini)

// Scheduled for shutdown by OpenAI on December 11, 2026. Use GPT56Sol (gpt-5.6-sol) as a replacement.
type GPT5 struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any

	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT5) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000125
}

func (g GPT5) GetName() string {
	return GPT5Alias
}

func (g GPT5) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT5)

// Scheduled for shutdown by OpenAI on December 11, 2026. Use GPT56Terra (gpt-5.6-terra) as a replacement.
type GPT5Mini struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any

	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT5Mini) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000025
}

func (g GPT5Mini) GetName() string {
	return GPT5MiniAlias
}

func (g GPT5Mini) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT5Mini)

// Scheduled for shutdown by OpenAI on December 11, 2026. Use GPT56Luna (gpt-5.6-luna) as a replacement.
type GPT5Nano struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any

	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT5Nano) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 5e-8
}

func (g GPT5Nano) GetName() string {
	return GPT5NanoAlias
}

func (g GPT5Nano) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT5Nano)

type GPT5Chat struct {
	// StructuredOutput represents a subset of the JSON Schema Language. Refer to openai documentation for complete and up-to-date information. An example structure could be:
	//
	//  var schema = map[string]any{
	//  	"name": "navidia_valuation",
	//  	"schema": map[string]any{
	//  		"type": "object",
	//  		"properties": map[string]any{
	//  			"final_answer": map[string]any{"type": "string"},
	//  			"valuation": map[string]any{
	//  				"type": "number",
	//  			},
	//  		},
	//  	},
	//  }
	StructuredOutput map[string]any

	// PdfFile let's you include a PDF file in your request to the LLM.
	// The expected format:
	//
	// map["file-name.pdf"]"data:application/pdf;base64," + encodedString
	// Only provide a pdf file or an image file, not both.
	PdfFile map[string]string
	// ImageFile enables vision for the request
	ImageFile []OpenaiImagePayload
}

func (g GPT5Chat) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000125
}

func (g GPT5Chat) GetName() string {
	return GPT5ChatAlias
}

func (g GPT5Chat) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT5Chat)

type GPT51 struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT51) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000125
}

func (g GPT51) GetName() string {
	return GPT51Alias
}

func (g GPT51) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT51)

type GPT51Chat struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT51Chat) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000125
}

func (g GPT51Chat) GetName() string {
	return GPT51ChatAlias
}

func (g GPT51Chat) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT51Chat)

type GPT51Codex struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT51Codex) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000125
}

func (g GPT51Codex) GetName() string {
	return GPT51CodexAlias
}

func (g GPT51Codex) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT51Codex)

type GPT51CodexMini struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT51CodexMini) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000025
}

func (g GPT51CodexMini) GetName() string {
	return GPT51CodexMiniAlias
}

func (g GPT51CodexMini) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT51CodexMini)

type GPT52 struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT52) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000175
}

func (g GPT52) GetName() string {
	return GPT52Alias
}

func (g GPT52) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT52)

type GPT53Codex struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT53Codex) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000175
}

func (g GPT53Codex) GetInputCostPer1M() float64 {
	return 1.75
}

func (g GPT53Codex) GetOutputCostPer1M() float64 {
	return 14.0
}

func (g GPT53Codex) GetName() string {
	return GPT53CodexAlias
}

func (g GPT53Codex) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT53Codex)
var _ CostBreakdown = new(GPT53Codex)

type GPT54 struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT54) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.0000025
}

func (g GPT54) GetInputCostPer1M() float64 {
	return 2.5
}

func (g GPT54) GetOutputCostPer1M() float64 {
	return 15.0
}

func (g GPT54) GetName() string {
	return GPT54Alias
}

func (g GPT54) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT54)
var _ CostBreakdown = new(GPT54)

type GPT54Mini struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT54Mini) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000075
}

func (g GPT54Mini) GetInputCostPer1M() float64 {
	return 0.75
}

func (g GPT54Mini) GetOutputCostPer1M() float64 {
	return 4.5
}

func (g GPT54Mini) GetName() string {
	return GPT54MiniAlias
}

func (g GPT54Mini) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT54Mini)
var _ CostBreakdown = new(GPT54Mini)

type GPT54Nano struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT54Nano) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.0000002
}

func (g GPT54Nano) GetInputCostPer1M() float64 {
	return 0.20
}

func (g GPT54Nano) GetOutputCostPer1M() float64 {
	return 1.25
}

func (g GPT54Nano) GetName() string {
	return GPT54NanoAlias
}

func (g GPT54Nano) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT54Nano)
var _ CostBreakdown = new(GPT54Nano)

type GPT55 struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT55) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000005
}

func (g GPT55) GetInputCostPer1M() float64 {
	return 5.0
}

func (g GPT55) GetOutputCostPer1M() float64 {
	return 30.0
}

func (g GPT55) GetName() string {
	return GPT55Alias
}

func (g GPT55) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT55)
var _ CostBreakdown = new(GPT55)

type GPT56Sol struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT56Sol) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000004
}

func (g GPT56Sol) GetInputCostPer1M() float64 {
	return 4.0
}

func (g GPT56Sol) GetOutputCostPer1M() float64 {
	return 20.0
}

func (g GPT56Sol) GetName() string {
	return GPT56SolAlias
}

func (g GPT56Sol) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT56Sol)
var _ CostBreakdown = new(GPT56Sol)

type GPT56Terra struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT56Terra) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000002
}

func (g GPT56Terra) GetInputCostPer1M() float64 {
	return 2.0
}

func (g GPT56Terra) GetOutputCostPer1M() float64 {
	return 12.0
}

func (g GPT56Terra) GetName() string {
	return GPT56TerraAlias
}

func (g GPT56Terra) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT56Terra)
var _ CostBreakdown = new(GPT56Terra)

type GPT56Luna struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (g GPT56Luna) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.0000002
}

func (g GPT56Luna) GetInputCostPer1M() float64 {
	return 0.20
}

func (g GPT56Luna) GetOutputCostPer1M() float64 {
	return 1.20
}

func (g GPT56Luna) GetName() string {
	return GPT56LunaAlias
}

func (g GPT56Luna) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT56Luna)
var _ CostBreakdown = new(GPT56Luna)

// GPT6Astra does not accept temperature or top_p; the provider drops those request fields.
type GPT6Astra struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

// Prompts over 272K input tokens are billed at a higher rate.
func (g GPT6Astra) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00001
}

func (g GPT6Astra) GetInputCostPer1M() float64 {
	return 10.0
}

func (g GPT6Astra) GetOutputCostPer1M() float64 {
	return 50.0
}

func (g GPT6Astra) GetName() string {
	return GPT6AstraAlias
}

func (g GPT6Astra) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT6Astra)
var _ CostBreakdown = new(GPT6Astra)

// GPT6Sol accepts temperature and top_p only at reasoning effort none, which heimdall does
// not set, so the provider drops those request fields.
type GPT6Sol struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

// Prompts over 272K input tokens are billed at a higher rate.
func (g GPT6Sol) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000002
}

func (g GPT6Sol) GetInputCostPer1M() float64 {
	return 2.0
}

func (g GPT6Sol) GetOutputCostPer1M() float64 {
	return 10.0
}

func (g GPT6Sol) GetName() string {
	return GPT6SolAlias
}

func (g GPT6Sol) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT6Sol)
var _ CostBreakdown = new(GPT6Sol)

// GPT6Luna accepts temperature and top_p only at reasoning effort none, which heimdall does
// not set, so the provider drops those request fields.
type GPT6Luna struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

// Prompts over 272K input tokens are billed at a higher rate.
func (g GPT6Luna) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.0000001
}

func (g GPT6Luna) GetInputCostPer1M() float64 {
	return 0.10
}

func (g GPT6Luna) GetOutputCostPer1M() float64 {
	return 0.50
}

func (g GPT6Luna) GetName() string {
	return GPT6LunaAlias
}

func (g GPT6Luna) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPT6Luna)
var _ CostBreakdown = new(GPT6Luna)

// Scheduled for shutdown by OpenAI on December 11, 2026. Use GPT56Sol (gpt-5.6-sol) as a replacement.
type O3 struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (o O3) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000200
}

func (o O3) GetName() string {
	return O3Alias
}

func (o O3) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(O3)

// Scheduled for shutdown by OpenAI on October 23, 2026. Use GPT56Terra (gpt-5.6-terra) as a replacement.
type O4Mini struct {
	StructuredOutput map[string]any
	PdfFile          map[string]string
	ImageFile        []OpenaiImagePayload
}

func (o O4Mini) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.00000110
}

func (o O4Mini) GetName() string {
	return O4MiniAlias
}

func (o O4Mini) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(O4Mini)

const (
	ImageModelAlias  = "gpt-image-1"
	ImageModel2Alias = "gpt-image-2"
)

// The sizes gpt-image-1 accepts. The 1792x1024 / 1024x1792 pair that used to
// sit here are DALL·E 3 sizes; gpt-image-1 rejects them outright.
const (
	GPTImageSizeAuto      = "auto"
	GPTImageSize1024x1024 = "1024x1024"
	GPTImageSize1536x1024 = "1536x1024"
	GPTImageSize1024x1536 = "1024x1536"

	GPTImageQualityHigh   = "high"
	GPTImageQualityMedium = "medium"
	GPTImageQualityLow    = "low"
)

// Scheduled for shutdown by OpenAI on October 23, 2026. Use GPTImage2 (gpt-image-2) as a replacement.
type GPTImage struct {
	// Allows to set transparency for the background of the generated image(s).
	// Must be one of transparent, opaque or auto (default value).
	// When auto is used, the model will automatically determine the best background for the image.
	// If transparent, the output format needs to support transparency, so it should be set to either png (default value) or webp.
	Background string

	// N is the number of images to generate. Must be 1 for DALL·E 3.
	// Although the API docs mention 'n', DALL-E 3 currently only supports n=1.
	// We keep it for potential future compatibility but default/validate to 1.
	N int

	// Size of the generated images: one of the GPTImageSize constants. Left
	// empty the parameter is not sent at all, so the API's own "auto" applies
	// and the model picks an aspect ratio to suit the prompt — which is what
	// you want for anything whose composition is not square.
	Size string

	// Quality of the image that will be generated. Defaults to "auto".
	Quality string

	// The compression level (0-100%) for the generated images.
	// This parameter is only supported the webp or jpeg output formats, and defaults to 100.
	OutputCompression string

	// The format in which the generated images are returned.
	// Must be one of png, jpeg, or webp.
	OutputFormat string

	// Must be either low for less restrictive filtering or auto (default value).
	Moderation string

	// User is an optional unique identifier representing your end-user,
	// which can help OpenAI monitor and detect abuse.
	User string

	// ImageFile enables image editing with gpt-image-1 (Image Edit API).
	// Input image must be less than 50 MB in size and must be a PNG or JPG file.
	// Note: Only 1 image is supported per request.
	ImageFile []OpenaiImagePayload
}

func (d GPTImage) EstimateCost(text string) float64 {
	// Text input tokens cost $5.00/1M, image output tokens cost $40.00/1M.
	// Use text input cost as a rough estimate for prompt size.
	return (float64(len(text)) / 4) * 0.000005
}

func (d GPTImage) GetName() string {
	return ImageModelAlias
}

func (d GPTImage) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPTImage)

// GPTImage2 targets gpt-image-2, the successor to gpt-image-1, and takes the
// same parameters. Size additionally accepts any WIDTHxHEIGHT string with both
// edges divisible by 16, an aspect ratio between 1:3 and 3:1 and at most
// 3840x2160; the GPTImageSize constants remain valid.
type GPTImage2 GPTImage

func (d GPTImage2) EstimateCost(text string) float64 {
	// Text input tokens cost $5.00/1M, image output tokens cost $30.00/1M.
	// Use text input cost as a rough estimate for prompt size.
	return (float64(len(text)) / 4) * 0.000005
}

func (d GPTImage2) GetName() string {
	return ImageModel2Alias
}

func (d GPTImage2) GetProvider() string {
	return OpenaiProvider
}

var _ Model = new(GPTImage2)
