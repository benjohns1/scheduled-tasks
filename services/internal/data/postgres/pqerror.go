package postgres

import "github.com/lib/pq"

// Postgres error codes
const (
	ErrUniqueViolation = pq.ErrorCode("23505")
	ErrCaseNotFound    = pq.ErrorCode("20000")
)

func isPqErr(err error, code pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == code
	}
	return false
}
