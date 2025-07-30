package usermodel

import "github.com/stretchr/testify/assert"

func (suite *UserModelSuite) TestOne() {
	t := suite.T()

	assert.Equal(t, 1, 1, "The result should be equal to the expected value")
}
