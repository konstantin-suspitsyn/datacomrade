package usermodel

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/testcntr"
	"github.com/stretchr/testify/suite"
)

type UserModelSuite struct {
	suite.Suite
	pgContainer *testcntr.PostgresContainer
	Model       Models
	ctx         context.Context
}

func (suite *UserModelSuite) SetupSuite() {
	ctx := context.Background()
	container, err := testcntr.New(ctx)

	if err != nil {
		panic(err.Error())
	}

	model := NewModel(container.DB)

	suite.ctx = ctx
	suite.pgContainer = container
	suite.Model = model

}

func (suite *UserModelSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		slog.Debug(fmt.Sprintf("error terminating postgres container: %s", err))
	}
}

func TestUserModelTestSuite(t *testing.T) {
	suite.Run(t, new(UserModelSuite))
}
