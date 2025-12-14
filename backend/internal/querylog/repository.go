package querylog

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrNotFound is returned when a query log record cannot be located.
var ErrNotFound = errors.New("query log not found")

// Repository persists and queries query log records.
type Repository struct {
	db *sql.DB
}

// NewRepository returns a repository backed by the supplied sql.DB handle.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ListParams defines filters and pagination for listing query logs.
type ListParams struct {
	Page          int
	Limit         int
	UserID        *int64
	APIKeyID      *int64
	Status        string
	Endpoint      string
	ModelProvider string
	StartDate     *time.Time
	EndDate       *time.Time
}

// Create inserts a new query log record.
func (r *Repository) Create(log *QueryLog) error {
	if log == nil {
		return fmt.Errorf("log is nil")
	}

	now := time.Now().UTC()
	log.CreatedAt = now

	var (
		apiKeyID       any
		conversationID any
		response       any
		modelProvider  any
		errorMessage   any
	)

	if log.APIKeyID != nil {
		apiKeyID = *log.APIKeyID
	}
	if log.ConversationID != nil {
		conversationID = *log.ConversationID
	}
	if log.Response != "" {
		response = log.Response
	}
	if log.ModelProvider != "" {
		modelProvider = log.ModelProvider
	}
	if log.ErrorMessage != "" {
		errorMessage = log.ErrorMessage
	}

	const insertQuery = `
		INSERT INTO query_logs (
			user_id, api_key_id, endpoint, query, response, model_provider,
			rag_contexts_count, input_tokens, output_tokens, latency_ms, status,
			error_message, conversation_id, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(insertQuery,
		log.UserID,
		apiKeyID,
		log.Endpoint,
		log.Query,
		response,
		modelProvider,
		log.RAGContextsCount,
		log.InputTokens,
		log.OutputTokens,
		log.LatencyMs,
		log.Status,
		errorMessage,
		conversationID,
		log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert query log: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fetch query log id: %w", err)
	}
	log.ID = id
	return nil
}

// GetByID returns a query log by its identifier.
func (r *Repository) GetByID(id int64) (*QueryLog, error) {
	const query = `
		SELECT
			id, user_id, api_key_id, endpoint, query, response, model_provider,
			rag_contexts_count, input_tokens, output_tokens, latency_ms, status,
			error_message, conversation_id, created_at
		FROM query_logs
		WHERE id = ?
	`

	var (
		log            QueryLog
		apiKeyID       sql.NullInt64
		conversationID sql.NullInt64
		response       sql.NullString
		modelProvider  sql.NullString
		errorMessage   sql.NullString
	)

	err := r.db.QueryRow(query, id).Scan(
		&log.ID,
		&log.UserID,
		&apiKeyID,
		&log.Endpoint,
		&log.Query,
		&response,
		&modelProvider,
		&log.RAGContextsCount,
		&log.InputTokens,
		&log.OutputTokens,
		&log.LatencyMs,
		&log.Status,
		&errorMessage,
		&conversationID,
		&log.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query query log: %w", err)
	}

	if apiKeyID.Valid {
		log.APIKeyID = &apiKeyID.Int64
	}
	if conversationID.Valid {
		log.ConversationID = &conversationID.Int64
	}
	if response.Valid {
		log.Response = response.String
	}
	if modelProvider.Valid {
		log.ModelProvider = modelProvider.String
	}
	if errorMessage.Valid {
		log.ErrorMessage = errorMessage.String
	}

	return &log, nil
}

// List returns paginated query logs matching the provided filters and the total count.
func (r *Repository) List(params ListParams) ([]QueryLog, int64, error) {
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 500 {
		limit = 500
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	whereParts := make([]string, 0)
	args := make([]any, 0)

	if params.UserID != nil {
		whereParts = append(whereParts, "user_id = ?")
		args = append(args, *params.UserID)
	}
	if params.APIKeyID != nil {
		whereParts = append(whereParts, "api_key_id = ?")
		args = append(args, *params.APIKeyID)
	}
	if params.Status != "" {
		whereParts = append(whereParts, "status = ?")
		args = append(args, params.Status)
	}
	if params.Endpoint != "" {
		whereParts = append(whereParts, "endpoint = ?")
		args = append(args, params.Endpoint)
	}
	if params.ModelProvider != "" {
		whereParts = append(whereParts, "model_provider = ?")
		args = append(args, params.ModelProvider)
	}
	if params.StartDate != nil {
		whereParts = append(whereParts, "created_at >= ?")
		args = append(args, *params.StartDate)
	}
	if params.EndDate != nil {
		whereParts = append(whereParts, "created_at <= ?")
		args = append(args, *params.EndDate)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM query_logs %s", whereClause)
	var total int64
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count query logs: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT
			id, user_id, api_key_id, endpoint, query, response, model_provider,
			rag_contexts_count, input_tokens, output_tokens, latency_ms, status,
			error_message, conversation_id, created_at
		FROM query_logs
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, whereClause)

	listArgs := append(append([]any{}, args...), limit, offset)

	rows, err := r.db.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list query logs: %w", err)
	}
	defer rows.Close()

	logs := make([]QueryLog, 0)

	for rows.Next() {
		var (
			log            QueryLog
			apiKeyID       sql.NullInt64
			conversationID sql.NullInt64
			response       sql.NullString
			modelProvider  sql.NullString
			errorMessage   sql.NullString
		)

		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&apiKeyID,
			&log.Endpoint,
			&log.Query,
			&response,
			&modelProvider,
			&log.RAGContextsCount,
			&log.InputTokens,
			&log.OutputTokens,
			&log.LatencyMs,
			&log.Status,
			&errorMessage,
			&conversationID,
			&log.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan query log: %w", err)
		}

		if apiKeyID.Valid {
			log.APIKeyID = &apiKeyID.Int64
		}
		if conversationID.Valid {
			log.ConversationID = &conversationID.Int64
		}
		if response.Valid {
			log.Response = response.String
		}
		if modelProvider.Valid {
			log.ModelProvider = modelProvider.String
		}
		if errorMessage.Valid {
			log.ErrorMessage = errorMessage.String
		}

		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate query logs: %w", err)
	}

	return logs, total, nil
}

// GetStats returns aggregated query log statistics for a date range.
// Zero-value startDate/endDate mean "no bound" for that side of the range.
func (r *Repository) GetStats(startDate, endDate time.Time) (*QueryLogStats, error) {
	whereParts := make([]string, 0)
	args := make([]any, 0)

	if !startDate.IsZero() {
		whereParts = append(whereParts, "created_at >= ?")
		args = append(args, startDate)
	}
	if !endDate.IsZero() {
		whereParts = append(whereParts, "created_at <= ?")
		args = append(args, endDate)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	stats := QueryLogStats{
		QueriesByEndpoint: make(map[string]int64),
		QueriesByProvider: make(map[string]int64),
	}

	aggregateQuery := fmt.Sprintf(`
		SELECT
			COUNT(*) AS total_queries,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) AS success_count,
			SUM(CASE WHEN status = 'error' THEN 1 ELSE 0 END) AS error_count,
			COALESCE(AVG(latency_ms), 0) AS avg_latency_ms,
			COALESCE(SUM(input_tokens), 0) AS total_input_tokens,
			COALESCE(SUM(output_tokens), 0) AS total_output_tokens
		FROM query_logs
		%s
	`, whereClause)

	if err := r.db.QueryRow(aggregateQuery, args...).Scan(
		&stats.TotalQueries,
		&stats.SuccessCount,
		&stats.ErrorCount,
		&stats.AvgLatencyMs,
		&stats.TotalInputTokens,
		&stats.TotalOutputTokens,
	); err != nil {
		return nil, fmt.Errorf("aggregate stats: %w", err)
	}

	endpointQuery := fmt.Sprintf(`
		SELECT endpoint, COUNT(*) FROM query_logs
		%s
		GROUP BY endpoint
	`, whereClause)

	if err := r.collectCounts(endpointQuery, args, stats.QueriesByEndpoint); err != nil {
		return nil, fmt.Errorf("aggregate endpoint stats: %w", err)
	}

	providerQuery := fmt.Sprintf(`
		SELECT COALESCE(model_provider, ''), COUNT(*) FROM query_logs
		%s
		GROUP BY model_provider
	`, whereClause)

	if err := r.collectCounts(providerQuery, args, stats.QueriesByProvider); err != nil {
		return nil, fmt.Errorf("aggregate provider stats: %w", err)
	}

	return &stats, nil
}

// DeleteOlderThan removes query log records older than the provided timestamp.
func (r *Repository) DeleteOlderThan(date time.Time) (int64, error) {
	res, err := r.db.Exec("DELETE FROM query_logs WHERE created_at < ?", date)
	if err != nil {
		return 0, fmt.Errorf("delete query logs: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}
	return rows, nil
}

func (r *Repository) collectCounts(query string, args []any, target map[string]int64) error {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var count int64
		if err := rows.Scan(&key, &count); err != nil {
			return err
		}
		target[key] = count
	}

	return rows.Err()
}
