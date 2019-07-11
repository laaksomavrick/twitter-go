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
	UserB map[string]interface{}
}

func (suite *TweetsTestSuite) SetupSuite() {
	// TODO: have this be set by an ENV when k8s is up; test against k8s
	// Will need to create a new keyspace + tables for above use case to not blow up prod?
	suite.Init("localhost", "3002")
	suite.Truncate([]string{"users", "tweets", "tweets_by_user"})

	// Create a new user
	statusCode, userA := suite.CreateUserViaHTTP(map[string]string{
		"username":             "username",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	if statusCode != 200 {
		suite.Fail("Unable to create a user for tweets_test")
	}

	suite.UserA = userA["data"].(map[string]interface{})

	// Create another new user
	statusCode, userB := suite.CreateUserViaHTTP(map[string]string{
		"username":             "username2",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email2@gmail.com",
		"displayName":          "displayName2",
	})

	if statusCode != 200 {
		suite.Fail("Unable to create a user for tweets_test")
	}

	suite.UserB = userB["data"].(map[string]interface{})
}

func (suite *TweetsTestSuite) TestCreateTweetSuccess() {
	accessToken := suite.UserA["accessToken"].(string)
	statusCode, createTweetResponse := suite.createTweetViaHTTP(map[string]string{
		"content": "this is a tweet",
	}, accessToken)

	assert.Equal(suite.T(), 200, statusCode)
	assert.NotNil(suite.T(), createTweetResponse["data"].(map[string]interface{})["content"])
	assert.NotNil(suite.T(), createTweetResponse["data"].(map[string]interface{})["id"])
}

func (suite *TweetsTestSuite) TestCreateTweetNotAuthorized() {
	accessToken := ""
	statusCode, _ := suite.createTweetViaHTTP(map[string]string{
		"content": "this is a tweet",
	}, accessToken)

	assert.Equal(suite.T(), 401, statusCode)
}

func (suite *TweetsTestSuite) TestCreateTweetInvalid() {
	accessToken := suite.UserA["accessToken"].(string)
	statusCode, _ := suite.createTweetViaHTTP(map[string]string{
		"content": "",
	}, accessToken)

	assert.Equal(suite.T(), 422, statusCode)
}

func (suite *TweetsTestSuite) TestGetTweetsSuccess() {
	username := suite.UserA["username"].(string)
	statusCode, _ := suite.getTweetsViaHTTP(username)
	assert.Equal(suite.T(), 200, statusCode)
}

func (suite *TweetsTestSuite) TestGetTweetsNotFound() {
	statusCode, _ := suite.getTweetsViaHTTP("someOtherUsername")
	assert.Equal(suite.T(), 404, statusCode)
}

func TestTweetsTestSuite(t *testing.T) {
	suite.Run(t, new(TweetsTestSuite))
}

func (suite *TweetsTestSuite) createTweetViaHTTP(request map[string]string, accessToken string) (statusCode int, createTweetResponse map[string]interface{}) {
	marshalled, _ := json.Marshal(request)
	body := bytes.NewBuffer(marshalled)
	req, _ := http.NewRequest("POST", suite.GetBaseURLWithSuffix("/tweets"), body)

	req.Header.Add("Authorization", "Bearer "+accessToken)

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

func (suite *TweetsTestSuite) getTweetsViaHTTP(username string) (statusCode int, getTweetsResponse map[string]interface{}) {
	url := fmt.Sprintf("/tweets/%s", username)
	resp, err := http.Get(suite.GetBaseURLWithSuffix(url))
	if err != nil {
		suite.Fail("Received no response from /tweets")
	}

	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(&getTweetsResponse)

	return resp.StatusCode, getTweetsResponse
}
