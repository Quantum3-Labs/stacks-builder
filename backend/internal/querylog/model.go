package querylog

import "time"

// QueryLog represents a single tracked request/response cycle for analytics and debugging.
type QueryLog struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	APIKeyID         *int64    `json:"api_key_id,omitempty"`
	Endpoint         string    `json:"endpoint"`
	Query            string    `json:"query"`
	Response         string    `json:"response,omitempty"`
	ModelProvider    string    `json:"model_provider,omitempty"`
	RAGContextsCount int       `json:"rag_contexts_count"`
	InputTokens      int       `json:"input_tokens"`
	OutputTokens     int       `json:"output_tokens"`
	LatencyMs        int64     `json:"latency_ms"`
	Status           string    `json:"status"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	ConversationID   *int64    `json:"conversation_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// QueryLogStats aggregates query log metrics for reporting.
type QueryLogStats struct {
	TotalQueries      int64            `json:"total_queries"`
	SuccessCount      int64            `json:"success_count"`
	ErrorCount        int64            `json:"error_count"`
	AvgLatencyMs      float64          `json:"avg_latency_ms"`
	TotalInputTokens  int64            `json:"total_input_tokens"`
	TotalOutputTokens int64            `json:"total_output_tokens"`
	QueriesByEndpoint map[string]int64 `json:"queries_by_endpoint"`
	QueriesByProvider map[string]int64 `json:"queries_by_provider"`
}
