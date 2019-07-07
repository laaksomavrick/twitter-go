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

type FollowersTestSuite struct {
	helpers.IntegrationTestSuite
	UserA map[string]interface{}
	UserB map[string]interface{}
}

func (suite *FollowersTestSuite) SetupSuite() {
	// TODO-13: have this be set by an ENV when k8s is up; test against k8s
	// Will need to create a new keyspace + tables for above use case to not blow up prod?
	suite.Init("localhost", "3002")
	suite.Truncate([]string{"users", "user_followers", "user_followings"})

	// Create a new user
	statusCode, userA := suite.CreateUserViaHTTP(map[string]string{
		"username":             "username",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	if statusCode != 200 {
		suite.Fail("Unable to create a user for followers_test")
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

func (suite *FollowersTestSuite) TestFollowerUserSuccess() {
	accessToken := suite.UserA["accessToken"].(string)
	followingUsername := suite.UserB["username"].(string)
	statusCode, createTweetResponse := suite.followViaHTTP(map[string]string{
		"followingUsername": followingUsername,
	}, accessToken)

	assert.Equal(suite.T(), 200, statusCode)
	assert.NotNil(suite.T(), createTweetResponse["data"].(map[string]interface{})["followingUsername"])
	assert.NotNil(suite.T(), createTweetResponse["data"].(map[string]interface{})["username"])
}

func (suite *FollowersTestSuite) TestFollowNotAuthorized() {
	accessToken := ""
	statusCode, _ := suite.followViaHTTP(map[string]string{}, accessToken)
	assert.Equal(suite.T(), 401, statusCode)
}

func (suite *FollowersTestSuite) TestFollowNotFound() {
	accessToken := suite.UserA["accessToken"].(string)
	followingUsername := "404"
	statusCode, _ := suite.followViaHTTP(map[string]string{
		"followingUsername": followingUsername,
	}, accessToken)

	assert.Equal(suite.T(), 404, statusCode)
}

func (suite *FollowersTestSuite) TestFollowSameUsername() {
	accessToken := suite.UserA["accessToken"].(string)
	followingUsername := suite.UserA["username"].(string) // UserA, not B
	statusCode, _ := suite.followViaHTTP(map[string]string{
		"followingUsername": followingUsername,
	}, accessToken)

	assert.Equal(suite.T(), 422, statusCode)
}

func (suite *FollowersTestSuite) TestFollowBadRequest() {
	accessToken := suite.UserA["accessToken"].(string)
	statusCode, _ := suite.followViaHTTP(map[string]string{
		"followingUsername": "",
	}, accessToken)

	assert.Equal(suite.T(), 422, statusCode)
}

func TestFollowersTestSuite(t *testing.T) {
	suite.Run(t, new(FollowersTestSuite))
}

func (suite *FollowersTestSuite) followViaHTTP(request map[string]string, accessToken string) (statusCode int, createTweetResponse map[string]interface{}) {
	marshalled, _ := json.Marshal(request)
	body := bytes.NewBuffer(marshalled)
	req, _ := http.NewRequest("POST", suite.GetBaseURLWithSuffix("/follow"), body)

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		suite.Fail("Received no response from /follow")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&createTweetResponse); err != nil {
		suite.Fail("Failed parsing response body")
	}

	return resp.StatusCode, createTweetResponse
}
