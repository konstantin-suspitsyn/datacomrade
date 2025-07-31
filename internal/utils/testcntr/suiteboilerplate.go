package testcntr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/stretchr/testify/suite"
)

type TestContainerSuiteBoilerPlate struct {
	suite.Suite
	PgContainter *PostgresContainer
	ctx          context.Context
}

func (suite *TestContainerSuiteBoilerPlate) SetupSuite() {
	ctx := context.Background()
	container, err := New(ctx)

	if err != nil {
		panic(err.Error())
	}

	suite.ctx = ctx
	suite.PgContainter = container

}

func (suite *TestContainerSuiteBoilerPlate) TearDownSuite() {
	if err := suite.PgContainter.Terminate(suite.ctx); err != nil {
		slog.Debug(fmt.Sprintf("error terminating postgres container: %s", err))
	}
}
