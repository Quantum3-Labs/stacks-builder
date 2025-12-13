package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Quantum3-Labs/stacks-builder/backend/internal/querylog"
)

// ListQueryLogs returns paginated query logs with optional filters.
func ListQueryLogs(repo *querylog.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		params := querylog.ListParams{
			Page:          page,
			Limit:         limit,
			Status:        c.Query("status"),
			Endpoint:      c.Query("endpoint"),
			ModelProvider: c.Query("model_provider"),
		}

		if userID, ok := parseInt64Ptr(c.Query("user_id")); ok {
			params.UserID = userID
		}
		if apiKeyID, ok := parseInt64Ptr(c.Query("api_key_id")); ok {
			params.APIKeyID = apiKeyID
		}
		if start, ok := parseDate(c.Query("start_date")); ok {
			params.StartDate = &start
		}
		if end, ok := parseDate(c.Query("end_date")); ok {
			params.EndDate = &end
		}

		logs, total, err := repo.List(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list query logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"logs":  logs,
			"total": total,
			"page":  params.Page,
			"limit": params.Limit,
		})
	}
}

// GetQueryLog returns a single query log by ID.
func GetQueryLog(repo *querylog.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		logEntry, err := repo.GetByID(id)
		if err != nil {
			if err == querylog.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "query log not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch query log"})
			return
		}

		c.JSON(http.StatusOK, logEntry)
	}
}

// GetQueryLogStats returns aggregated statistics over a date range.
func GetQueryLogStats(repo *querylog.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var startDate, endDate time.Time
		if start, ok := parseDate(c.Query("start_date")); ok {
			startDate = start
		}
		if end, ok := parseDate(c.Query("end_date")); ok {
			endDate = end
		}

		stats, err := repo.GetStats(startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch query log stats"})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

func parseInt64Ptr(val string) (*int64, bool) {
	if val == "" {
		return nil, false
	}
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, false
	}
	return &parsed, true
}

// parseDate parses YYYY-MM-DD or RFC3339 timestamps.
func parseDate(val string) (time.Time, bool) {
	if val == "" {
		return time.Time{}, false
	}
	if t, err := time.Parse("2006-01-02", val); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return t, true
	}
	return time.Time{}, false
}
