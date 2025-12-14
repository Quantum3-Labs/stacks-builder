package codegen

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"
)

const (
	defaultGeminiModel     = "gemini-2.5-flash"
	defaultGeminiMaxTokens = 8192
)

// GeminiService handles code generation using Gemini API
type GeminiService struct {
	client *genai.Client
}

// NewGeminiService creates a new Gemini service
func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &GeminiService{client: client}, nil
}

// NewGeminiServiceFromEnv creates a new Gemini service using environment variables
func NewGeminiServiceFromEnv() (*GeminiService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	return NewGeminiService(apiKey)
}

// GenerateCode generates Clarity code using Gemini with provided context
func (s *GeminiService) GenerateCode(ctx context.Context, query string, codeContexts []string, docContexts []string, temperature float64, maxTokens int) (*CodeGenerationResponse, error) {
	// Assemble prompt with context
	prompt := buildCodeGenerationInstruction(query, codeContexts, docContexts)

	// Set defaults
	if temperature == 0 {
		temperature = 0.7
	}
	if maxTokens == 0 {
		maxTokens = defaultGeminiMaxTokens
	}

	// Count input tokens
	inputTokenCount, err := s.countTokens(ctx, prompt)
	if err != nil {
		log.Printf("Warning: failed to count input tokens: %v", err)
		inputTokenCount = 0 // fallback to 0 if counting fails
	}

	// Call Gemini API
	geminiResponse, err := s.callGemini(ctx, prompt, temperature, maxTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}

	// Count output tokens
	outputTokenCount, err := s.countTokens(ctx, geminiResponse)
	if err != nil {
		log.Printf("Warning: failed to count output tokens: %v", err)
		outputTokenCount = 0
	}

	// Parse and return response
	parsedResponse, err := s.parseGeminiResponse(geminiResponse)
	if err != nil {
		return nil, err
	}

	// Add token counts
	parsedResponse.InputTokens = inputTokenCount
	parsedResponse.OutputTokens = outputTokenCount

	return parsedResponse, nil
}

// callGemini calls the Gemini API using the go-genai SDK
func (s *GeminiService) callGemini(ctx context.Context, prompt string, temperature float64, maxTokens int) (string, error) {
	config := &genai.GenerateContentConfig{
		Temperature:     genai.Ptr(float32(temperature)),
		MaxOutputTokens: int32(maxTokens),
	}

	result, err := s.client.Models.GenerateContent(
		ctx,
		defaultGeminiModel,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	return result.Text(), nil
}

// parseGeminiResponse extracts code and explanation from Gemini's response
func (s *GeminiService) parseGeminiResponse(response string) (*CodeGenerationResponse, error) {
	// Try to extract code block
	code := extractCodeBlock(response, "clarity")
	if code == "" {
		code = extractCodeBlock(response, "")
	}

	// Extract explanation (everything outside code blocks)
	explanation := removeCodeBlocks(response)

	return &CodeGenerationResponse{
		Code:         code,
		Explanation:  strings.TrimSpace(explanation),
		InputTokens:  0, // Will be set by GenerateCode
		OutputTokens: 0, // Will be set by GenerateCode
	}, nil
}

// countTokens counts tokens in text using Gemini's CountTokens API
func (s *GeminiService) countTokens(ctx context.Context, text string) (int, error) {
	result, err := s.client.Models.CountTokens(
		ctx,
		defaultGeminiModel,
		genai.Text(text),
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("token counting failed: %w", err)
	}

	return int(result.TotalTokens), nil
}
