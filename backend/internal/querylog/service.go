package querylog

import "log"

// Service provides asynchronous logging over a buffered channel.
type Service struct {
	repo    *Repository
	logChan chan *QueryLog
}

// NewService constructs a Service with a buffered channel and background worker.
func NewService(repo *Repository) *Service {
	s := &Service{
		repo:    repo,
		logChan: make(chan *QueryLog, 1000),
	}
	go s.processLogs()
	return s
}

// LogAsync enqueues a log entry without blocking callers.
func (s *Service) LogAsync(log *QueryLog) {
	select {
	case s.logChan <- log:
	default:
		// Drop when buffer is full to avoid backpressure on request path.
	}
}

func (s *Service) processLogs() {
	for logEntry := range s.logChan {
		if err := s.repo.Create(logEntry); err != nil {
			log.Printf("querylog: failed to persist query log: %v", err)
		}
	}
}
