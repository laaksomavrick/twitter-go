package tests

import (
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
	suite.Init("localhost", "3000")
	suite.Truncate([]string{"users"})
}

func (suite *UsersTestSuite) TestExample() {
	resp, err := http.Get(suite.GetBaseURLWithSuffix("/hello"))
	if err != nil {
		suite.Fail("Received no response from /hello")
	}
	assert.Equal(suite.T(), 200, resp.StatusCode)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
