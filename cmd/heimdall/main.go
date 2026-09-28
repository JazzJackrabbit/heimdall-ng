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
	ctx := context.Background()

	google := providers.NewGoogle([]string{os.Getenv("GOOGLE_API_KEY")})
	router := heimdall.New(30*time.Second, []heimdall.LLMProvider{google})

	res, err := router.Complete(ctx, request.Completion{
		Model:         models.Gemini38Flash{},
		SystemMessage: "You are a helpful assistant.",
		UserMessage:   "What is the capital of France?",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println(res.Content)
}
