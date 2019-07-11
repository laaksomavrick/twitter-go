package service

import (
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/logger"

	"github.com/gocql/gocql"
)

// Repository is a wrapper around gocql with a set of convenience functions
// for common operations
type Repository struct {
	Cassandra *cassandra.Client
}

// Query wraps makes a query gocql can understand, adding a log message for auditing
func (r *Repository) Query(queryString string, values ...interface{}) gocql.Query {
	query := r.Cassandra.Session.Query(queryString, values...)
	logger.Info(logger.Loggable{Message: "Executing query", Data: query.String()})
	return *query
}
