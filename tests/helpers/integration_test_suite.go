package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/env"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	cassandra *cassandra.Client
	Host      string
	Port      string
}

// Singleton for shared access during local dev (helps boot up times)
var singletonClient *cassandra.Client

// Init initializes the test suite
func (suite *IntegrationTestSuite) Init(host string, port string) error {
	suite.Host = host
	suite.Port = port
	if singletonClient == nil {
		cassandraURL := env.GetEnv("CASSANDRA_URL", "127.0.0.1")
		cassandraKeyspace := env.GetEnv("CASSANDRA_KEYSPACE", "twtr")
		cassandra, err := cassandra.NewClient(cassandraURL, cassandraKeyspace)
		if err != nil {
			return err
		}
		singletonClient = cassandra
		suite.cassandra = singletonClient
	} else {
		return nil
	}
	return nil
}

// Truncate truncates the specified tables in cassandra
func (suite *IntegrationTestSuite) Truncate(tables []string) error {
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE %s", table)
		err := suite.cassandra.Session.Query(query).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetBaseURL retrieves the base url for the current test suite for http requests
func (suite *IntegrationTestSuite) GetBaseURL() string {
	return fmt.Sprintf("http://%s:%s", suite.Host, suite.Port)
}

// GetBaseURLWithSuffix retrieves the base url for the current test suite with a suffix for http requests
func (suite *IntegrationTestSuite) GetBaseURLWithSuffix(suffix string) string {
	return fmt.Sprintf("http://%s:%s%s", suite.Host, suite.Port, suffix)
}

// CreateUserViaHTTP creates a user via http. This function is used across all tests,
// as the application is predicated around users
func (suite *IntegrationTestSuite) CreateUserViaHTTP(request map[string]string) (statusCode int, createUserResponse map[string]interface{}) {
	marshalled, err := json.Marshal(request)
	body := bytes.NewBuffer(marshalled)

	resp, err := http.Post((suite.GetBaseURLWithSuffix("/users")), "application/json", body)
	if err != nil {
		suite.Fail("Received no response from /users")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&createUserResponse); err != nil {
		suite.Fail("Failed parsing response body")
	}

	return resp.StatusCode, createUserResponse
}
