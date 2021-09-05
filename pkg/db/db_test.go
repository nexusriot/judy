package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"go.uber.org/zap"
)

type JudyDbTestSuite struct {
	suite.Suite
	logger *zap.Logger
	db     *JudyDb
	ctx    context.Context
}

func (s *JudyDbTestSuite) SetupSuite() {
	var loggerConfig = zap.NewDevelopmentConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, _ := loggerConfig.Build()
	s.logger = logger
	s.db = NewJudyDb(true, s.logger)
	s.ctx = context.Background()
}

func (s *JudyDbTestSuite) TestJudyDb_AddClient() {
	s.db.AddClient(uuid.New().String(), "127.0.0.1")
}

func (s *JudyDbTestSuite) TestJudyDb_AddTask() {
	clientUuid := uuid.New().String()
	s.db.AddClient(clientUuid, "127.0.0.1")
	s.db.AddTask(clientUuid, "test")
}

func (s *JudyDbTestSuite) TestJudyDb_TakeTasks() {
	clientUuid := uuid.New().String()
	s.db.AddClient(clientUuid, "127.0.0.1")
	s.db.AddTask(clientUuid, "1")
	s.db.AddTask(clientUuid, "2")
	s.db.AddTask(clientUuid, "3")
	tasks, err := s.db.TakeTasks(clientUuid)
	s.Require().Nil(err)
	s.Require().Equal(3, len(tasks))
}

func TestJudyDbTestSuite(t *testing.T) {
	suite.Run(t, new(JudyDbTestSuite))
}
