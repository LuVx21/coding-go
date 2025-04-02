package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MyTestSuite struct {
	suite.Suite
}

func (s *MyTestSuite) SetupTest() {
	fmt.Println("每个测试前的初始化")
}

func (s *MyTestSuite) TearDownTest() {
	fmt.Println("每个测试后的清理")
}

func (suite *MyTestSuite) Test_00() {
	assert.Equal(suite.T(), 5, 5)
}
func (suite *MyTestSuite) Test_01() {
	assert.Equal(suite.T(), 5, 5)
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
