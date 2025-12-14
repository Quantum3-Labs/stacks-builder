package middleware

import (
	"bytes"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Quantum3-Labs/stacks-builder/backend/internal/querylog"
)

// Context keys for handler-specific data.
const (
	QueryLogModelProvider    = "querylog_model_provider"
	QueryLogInputTokens      = "querylog_input_tokens"
	QueryLogOutputTokens     = "querylog_output_tokens"
	QueryLogRAGContextsCount = "querylog_rag_contexts_count"
	QueryLogConversationID   = "querylog_conversation_id"
	QueryLogErrorMessage     = "querylog_error_message"
)

// responseWriter wraps gin.ResponseWriter to capture the response body.
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	_, _ = w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	_, _ = w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// QueryLogMiddleware captures request/response data for tracked endpoints and logs asynchronously.
func QueryLogMiddleware(service *querylog.Service, trackedEndpoints []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		if !isTrackedEndpoint(path, trackedEndpoints) {
			c.Next()
			return
		}

		requestBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		rw := &responseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = rw

		startTime := time.Now()

		c.Next() // Execute the rest of the chain/handler

		latencyMs := time.Since(startTime).Milliseconds()

		logEntry := &querylog.QueryLog{
			Endpoint:  path,
			Query:     extractQuery(requestBody),
			Response:  truncateResponse(rw.body.String(), 10000),
			LatencyMs: latencyMs,
			Status:    getStatus(c.Writer.Status()),
			CreatedAt: time.Now().UTC(),
		}

		if userID, ok := c.Get("user_id"); ok {
			if id, ok := toInt64(userID); ok {
				logEntry.UserID = id
			}
		}

		if apiKeyID, ok := c.Get("api_key_id"); ok {
			if id, ok := toInt64(apiKeyID); ok {
				logEntry.APIKeyID = &id
			}
		}

		if provider, ok := c.Get(QueryLogModelProvider); ok {
			if v, ok := provider.(string); ok {
				logEntry.ModelProvider = v
			}
		}
		if tokens, ok := c.Get(QueryLogInputTokens); ok {
			if v, ok := toInt(tokens); ok {
				logEntry.InputTokens = v
			}
		}
		if tokens, ok := c.Get(QueryLogOutputTokens); ok {
			if v, ok := toInt(tokens); ok {
				logEntry.OutputTokens = v
			}
		}
		if count, ok := c.Get(QueryLogRAGContextsCount); ok {
			if v, ok := toInt(count); ok {
				logEntry.RAGContextsCount = v
			}
		}
		if convID, ok := c.Get(QueryLogConversationID); ok {
			if id, ok := toInt64(convID); ok {
				logEntry.ConversationID = &id
			}
		}
		if errMsg, ok := c.Get(QueryLogErrorMessage); ok {
			if v, ok := errMsg.(string); ok {
				logEntry.ErrorMessage = v
			}
		}

		// Require user_id to avoid foreign-key failures.
		if logEntry.UserID == 0 {
			log.Printf("querylog: skipping entry for %s, no user_id in context", path)
			return
		}

		service.LogAsync(logEntry)
	}
}

func isTrackedEndpoint(path string, tracked []string) bool {
	for _, e := range tracked {
		if e == path {
			return true
		}
	}
	return false
}

func extractQuery(body []byte) string {
	return strings.TrimSpace(string(body))
}

func truncateResponse(val string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(val) <= maxLen {
		return val
	}
	return val[:maxLen]
}

func getStatus(code int) string {
	if code >= 200 && code < 400 {
		return "success"
	}
	return "error"
}

func toInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int64:
		return val, true
	case int:
		return int64(val), true
	case int32:
		return int64(val), true
	}
	return 0, false
}

func toInt(v any) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int32:
		return int(val), true
	case int64:
		return int(val), true
	}
	return 0, false
}
