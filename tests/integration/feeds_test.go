package integration

import (
	"net/http"
	"testing"
	"twitter-go/tests/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FeedsTestSuite struct {
	helpers.IntegrationTestSuite
	UserA map[string]interface{}
	UserB map[string]interface{}
}

func (suite *FeedsTestSuite) SetupSuite() {
	// TODO: have this be set by an ENV when k8s is up; test against k8s
	// Will need to create a new keyspace + tables for above use case to not blow up prod?
	suite.Init("localhost", "3002")
	suite.Truncate([]string{"users", "feed_items"})

	// Create a new user
	statusCode, userA := suite.CreateUserViaHTTP(map[string]string{
		"username":             "username",
		"password":             "password",
		"passwordConfirmation": "password",
		"email":                "email@gmail.com",
		"displayName":          "displayName",
	})

	if statusCode != 200 {
		suite.Fail("Unable to create a user for feeds_test")
	}

	suite.UserA = userA["data"].(map[string]interface{})
}

func (suite *FeedsTestSuite) TestGetFeedSuccess() {
	accessToken := suite.UserA["accessToken"].(string)
	statusCode := suite.getFeedViaHTTP(accessToken)
	assert.Equal(suite.T(), 200, statusCode)
}

func (suite *FeedsTestSuite) TestGetFeedNotAuthorized() {
	accessToken := ""
	statusCode := suite.getFeedViaHTTP(accessToken)
	assert.Equal(suite.T(), 401, statusCode)
}

func TestFeedsTestSuite(t *testing.T) {
	suite.Run(t, new(FeedsTestSuite))
}

func (suite *FeedsTestSuite) getFeedViaHTTP(accessToken string) (statusCode int) {
	req, _ := http.NewRequest("GET", suite.GetBaseURLWithSuffix("/feeds/me"), nil)

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		suite.Fail("Received no response from /feeds/me")
	}

	defer resp.Body.Close()

	return resp.StatusCode
}
