package pqerr

import "github.com/lib/pq"

// Postgres error codes
const (
	UniqueViolation     = pq.ErrorCode("23505")
	CaseNotFound        = pq.ErrorCode("20000")
	ForeignKeyViolation = pq.ErrorCode("23503")
)

// Eq returns whether an error contains a specific PQ error code
func Eq(err error, code pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == code
	}
	return false
}
