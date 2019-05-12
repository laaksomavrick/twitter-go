package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"twitter-go/tests/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TweetsTestSuite struct {
	helpers.IntegrationTestSuite
	UserA map[string]interface{}
}

func (suite *TweetsTestSuite) SetupSuite() {
	// TODO-13: have this be set by an ENV when k8s is up; test against k8s
	// Will need to create a new keyspace + tables for above use case to not blow up prod?
	suite.Init("localhost", "3000")
	suite.Truncate([]string{"users", "tweets", "tweets_by_user"})

	// Create a new user
	statusCode, user := suite.CreateUserViaHTTP(map[string]string{
		"username":             "username",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	if statusCode != 200 {
		suite.Fail("Unable to create a user for tweets_test")
	}

	suite.UserA = user
}

func (suite *TweetsTestSuite) TestCreateTweetSuccess() {
	statusCode, createTweetResponse := suite.createTweetViaHTTP(map[string]string{
		"username": suite.getUsername(),
		"content":  "this is a tweet",
	})

	assert.Equal(suite.T(), 200, statusCode)
	assert.NotNil(suite.T(), createTweetResponse["content"])
	assert.NotNil(suite.T(), createTweetResponse["id"])
}

func (suite *TweetsTestSuite) TestCreateTweetForbidden() {
	statusCode, _ := suite.createTweetViaHTTP(map[string]string{
		"username": "someInvalidUsername",
		"content":  "this is a tweet",
	})

	assert.Equal(suite.T(), 403, statusCode)
}

func (suite *TweetsTestSuite) TestCreateTweetInvalid() {
	statusCode, _ := suite.createTweetViaHTTP(map[string]string{
		"username": suite.getUsername(),
		"content":  "",
	})

	assert.Equal(suite.T(), 400, statusCode)
}

func (suite *TweetsTestSuite) TestGetTweetsSuccess() {
	statusCode, _ := suite.getTweetsViaHTTP(suite.getUsername())
	assert.Equal(suite.T(), 200, statusCode)
}

func (suite *TweetsTestSuite) TestGetTweetsNotFound() {
	statusCode, _ := suite.getTweetsViaHTTP("someOtherUsername")
	assert.Equal(suite.T(), 404, statusCode)
}

func TestTweetsTestSuite(t *testing.T) {
	suite.Run(t, new(TweetsTestSuite))
}

func (suite *TweetsTestSuite) getUsername() string {
	return suite.UserA["username"].(string)
}

func (suite *TweetsTestSuite) getAccessToken() string {
	return suite.UserA["accessToken"].(string)
}

func (suite *TweetsTestSuite) createTweetViaHTTP(request map[string]string) (statusCode int, createTweetResponse map[string]interface{}) {
	marshalled, _ := json.Marshal(request)
	body := bytes.NewBuffer(marshalled)
	req, _ := http.NewRequest("POST", suite.GetBaseURLWithSuffix("/tweets"), body)

	req.Header.Add("Authorization", "Bearer "+suite.getAccessToken())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		suite.Fail("Received no response from /tweets")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&createTweetResponse); err != nil {
		suite.Fail("Failed parsing response body")
	}

	return resp.StatusCode, createTweetResponse
}

func (suite *TweetsTestSuite) getTweetsViaHTTP(username string) (statusCode int, getTweetsResponse []interface{}) {
	url := fmt.Sprintf("/tweets/%s", username)
	resp, err := http.Get(suite.GetBaseURLWithSuffix(url))
	if err != nil {
		suite.Fail("Received no response from /tweets")
	}

	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(&getTweetsResponse)

	return resp.StatusCode, getTweetsResponse
}
