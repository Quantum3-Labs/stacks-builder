package codegen

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

const (
	defaultOpenAIModel         = openai.ChatModelGPT4o
	defaultOpenAISystemMessage = "You are a Clarity expert."
	defaultOpenAIMaxTokens     = 4096
)

// OpenAIService handles code generation using OpenAI chat completions API.
type OpenAIService struct {
	client        openai.Client
	model         string
	systemMessage string
}

// NewOpenAIService creates a new OpenAI service instance.
func NewOpenAIService(apiKey, model, baseURL, systemMessage string) *OpenAIService {
	if model == "" {
		model = defaultOpenAIModel
	}
	if systemMessage == "" {
		systemMessage = defaultOpenAISystemMessage
	}

	// Build client options
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(opts...)

	return &OpenAIService{
		client:        client,
		model:         model,
		systemMessage: systemMessage,
	}
}

// NewOpenAIServiceFromEnv loads OpenAI configuration from environment variables.
func NewOpenAIServiceFromEnv() (*OpenAIService, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	model := os.Getenv("OPENAI_MODEL")
	baseURL := os.Getenv("OPENAI_BASE_URL")
	systemMessage := os.Getenv("OPENAI_SYSTEM_MESSAGE")

	return NewOpenAIService(apiKey, model, baseURL, systemMessage), nil
}

// GenerateCode calls the OpenAI API to generate code using provided contexts.
func (s *OpenAIService) GenerateCode(ctx context.Context, query string, codeContexts []string, docContexts []string, temperature float64, maxTokens int) (*CodeGenerationResponse, error) {
	if temperature == 0 {
		temperature = 0.7
	}
	if maxTokens == 0 {
		maxTokens = defaultOpenAIMaxTokens
	}

	prompt := buildCodeGenerationInstruction(query, codeContexts, docContexts)

	// Build the chat completion request
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(s.systemMessage),
			openai.UserMessage(prompt),
		},
		Model:       s.model,
		Temperature: param.NewOpt(temperature),
		MaxTokens:   param.NewOpt(int64(maxTokens)),
	}

	// Call the OpenAI API
	chatCompletion, err := s.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("openai response contained no choices")
	}

	assistantText := chatCompletion.Choices[0].Message.Content

	code := extractCodeBlock(assistantText, "clarity")
	if code == "" {
		code = extractCodeBlock(assistantText, "")
	}

	explanation := removeCodeBlocks(assistantText)

	return &CodeGenerationResponse{
		Code:         code,
		Explanation:  explanation,
		InputTokens:  int(chatCompletion.Usage.PromptTokens),
		OutputTokens: int(chatCompletion.Usage.CompletionTokens),
	}, nil
}
