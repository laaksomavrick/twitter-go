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
	// TODO: have this be set by an ENV when k8s is up; test against k8s
	// Will need to create a new keyspace + tables for above use case to not blow up prod?
	suite.Init("localhost", "3002")
	suite.Truncate([]string{"users"})
}

func (suite *UsersTestSuite) TestCreateUserSuccess() {
	statusCode, createUserResponse := suite.CreateUserViaHTTP(map[string]string{
		"username":             "username",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	assert.Equal(suite.T(), 200, statusCode)
	assert.NotNil(suite.T(), createUserResponse["data"].(map[string]interface{})["accessToken"])
	assert.NotNil(suite.T(), createUserResponse["data"].(map[string]interface{})["refreshToken"])
}

func (suite *UsersTestSuite) TestCreateUserBadRequestFailure() {
	statusCode, _ := suite.CreateUserViaHTTP(map[string]string{})
	assert.Equal(suite.T(), 422, statusCode)
}

func (suite *UsersTestSuite) TestCreateUserAlreadyExistsFailure() {
	request := map[string]string{
		"username":             "anotherUsername",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	}
	statusCodeUnique, _ := suite.CreateUserViaHTTP(request)
	statusCodeDuplicate, _ := suite.CreateUserViaHTTP(request)

	assert.Equal(suite.T(), 200, statusCodeUnique)
	assert.Equal(suite.T(), 409, statusCodeDuplicate)
}

func (suite *UsersTestSuite) TestAuthenticateUserSuccess() {
	username := "username2"
	password := "password"

	_, _ = suite.CreateUserViaHTTP(map[string]string{
		"username":             username,
		"password":             password,
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	statusCode, authenticateUserResponse := suite.authenticateUserViaHTTP(map[string]string{
		"username": username,
		"password": password,
	})

	assert.Equal(suite.T(), 200, statusCode)
	assert.NotNil(suite.T(), authenticateUserResponse["data"].(map[string]interface{})["accessToken"])
	assert.NotNil(suite.T(), authenticateUserResponse["data"].(map[string]interface{})["refreshToken"])
}

func (suite *UsersTestSuite) TestAuthenticateUserBadInputFailure() {

	statusCode, _ := suite.authenticateUserViaHTTP(map[string]string{})

	assert.Equal(suite.T(), 422, statusCode)
}

func (suite *UsersTestSuite) TestAuthenticateUserBadPasswordFailure() {
	username := "username3"
	password := "password"
	badPassword := "aardvark"

	_, _ = suite.CreateUserViaHTTP(map[string]string{
		"username":             username,
		"password":             password,
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	statusCode, _ := suite.authenticateUserViaHTTP(map[string]string{
		"username": username,
		"password": badPassword,
	})

	assert.Equal(suite.T(), 422, statusCode)
}

func (suite *UsersTestSuite) TestAuthenticateUserMissingUserFailure() {
	username := "someoneWhoDoesntExist"
	password := "password"

	statusCode, _ := suite.authenticateUserViaHTTP(map[string]string{
		"username": username,
		"password": password,
	})

	assert.Equal(suite.T(), 404, statusCode)
}

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

func (suite *UsersTestSuite) authenticateUserViaHTTP(request map[string]string) (statusCode int, authenticateUserResponse map[string]interface{}) {
	marshalled, err := json.Marshal(request)
	body := bytes.NewBuffer(marshalled)

	resp, err := http.Post((suite.GetBaseURLWithSuffix("/users/authorize")), "application/json", body)
	if err != nil {
		suite.Fail("Received no response from /authorize")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&authenticateUserResponse); err != nil {
		suite.Fail("Failed parsing response body")
	}

	return resp.StatusCode, authenticateUserResponse
}
