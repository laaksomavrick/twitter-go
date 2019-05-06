package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"twitter-go/tests/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UsersTestSuite struct {
	helpers.IntegrationTestSuite
}

func (suite *UsersTestSuite) SetupSuite() {
	// TODO-13: have this be set by an ENV when k8s is up; test against k8s
	// Will need to create a new keyspace + tables for above use case to not blow up prod?
	suite.Init("localhost", "3000")
	suite.Truncate([]string{"users"})
}

// Authorizing a user bad input  422

// Authorizing a user password mismatch 400

// Authorizing a user doesn't exist 404

func (suite *UsersTestSuite) TestCreateUserSuccess() {
	statusCode, createUserResponse := suite.createUserViaHTTP(map[string]string{
		"username":             "username",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	assert.Equal(suite.T(), 200, statusCode)
	assert.NotNil(suite.T(), createUserResponse["accessToken"])
	assert.NotNil(suite.T(), createUserResponse["refreshToken"])
}

func (suite *UsersTestSuite) TestCreateUserBadRequestFailure() {
	statusCode, _ := suite.createUserViaHTTP(map[string]string{})
	assert.Equal(suite.T(), 400, statusCode)
}

func (suite *UsersTestSuite) TestCreateUserAlreadyExistsFailure() {
	request := map[string]string{
		"username":             "anotherUsername",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	}
	statusCodeUnique, _ := suite.createUserViaHTTP(request)
	statusCodeDuplicate, _ := suite.createUserViaHTTP(request)

	assert.Equal(suite.T(), 200, statusCodeUnique)
	assert.Equal(suite.T(), 409, statusCodeDuplicate)
}

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

func (suite *UsersTestSuite) createUserViaHTTP(request map[string]string) (statusCode int, createUserResponse map[string]interface{}) {
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
