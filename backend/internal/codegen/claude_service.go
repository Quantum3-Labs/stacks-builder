package codegen

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const (
	defaultClaudeModel         = "claude-sonnet-4-5-20250514"
	defaultClaudeSystemMessage = "You are an expert Clarity programmer."
	defaultClaudeMaxTokens     = 4096
	defaultClaudeTemperature   = 0.7
)

// ClaudeService handles code generation using Anthropic Claude API.
type ClaudeService struct {
	client        anthropic.Client
	model         string
	systemMessage string
}

// NewClaudeService creates a new Claude service instance.
func NewClaudeService(apiKey, model, baseURL, apiVersion, systemMessage string) *ClaudeService {
	if model == "" {
		model = defaultClaudeModel
	}
	if systemMessage == "" {
		systemMessage = defaultClaudeSystemMessage
	}

	// Build client options
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	// Add base URL if provided
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	// Note: API version is handled automatically by the SDK
	// but we keep the parameter for backward compatibility
	_ = apiVersion

	client := anthropic.NewClient(opts...)

	return &ClaudeService{
		client:        client,
		model:         model,
		systemMessage: systemMessage,
	}
}

// NewClaudeServiceFromEnv loads Claude configuration from environment variables.
func NewClaudeServiceFromEnv() (*ClaudeService, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CLAUDE_API_KEY environment variable not set")
	}

	model := os.Getenv("CLAUDE_MODEL")
	baseURL := os.Getenv("CLAUDE_BASE_URL")
	apiVersion := os.Getenv("CLAUDE_API_VERSION")
	systemMessage := os.Getenv("CLAUDE_SYSTEM_MESSAGE")

	return NewClaudeService(apiKey, model, baseURL, apiVersion, systemMessage), nil
}

// GenerateCode calls Anthropic Claude API to generate code with provided contexts.
func (s *ClaudeService) GenerateCode(ctx context.Context, query string, codeContexts []string, docContexts []string, temperature float64, maxTokens int) (*CodeGenerationResponse, error) {
	if temperature == 0 {
		temperature = defaultClaudeTemperature
	}
	if maxTokens == 0 {
		maxTokens = defaultClaudeMaxTokens
	}

	prompt := buildCodeGenerationInstruction(query, codeContexts, docContexts)

	// Create message using SDK types
	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:       anthropic.Model(s.model),
		MaxTokens:   int64(maxTokens),
		Temperature: anthropic.Float(temperature),
		System: []anthropic.TextBlockParam{
			{Text: s.systemMessage},
		},
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					{
						OfText: &anthropic.TextBlockParam{
							Text: prompt,
						},
					},
				},
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate code with Claude: %w", err)
	}

	// Extract text from response
	var assistantText string
	for _, block := range message.Content {
		// Use type assertion to check for TextBlock
		if textBlock, ok := block.AsAny().(anthropic.TextBlock); ok {
			assistantText += textBlock.Text
		}
	}

	if assistantText == "" {
		return nil, fmt.Errorf("claude response contained no text content")
	}

	// Extract code blocks and explanation
	code := extractCodeBlock(assistantText, "clarity")
	if code == "" {
		code = extractCodeBlock(assistantText, "")
	}

	explanation := removeCodeBlocks(assistantText)

	return &CodeGenerationResponse{
		Code:         code,
		Explanation:  explanation,
		InputTokens:  int(message.Usage.InputTokens),
		OutputTokens: int(message.Usage.OutputTokens),
	}, nil
}
