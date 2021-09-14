package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/nexusriot/judy/pkg/db"
	pb "github.com/nexusriot/judy/proto"
)

type JudyServer struct {
	db     *db.JudyDb
	logger *zap.Logger
	engine *gin.Engine
}

func NewJudyServer() (*JudyServer, error) {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if nil != err {
		return nil, err
	}
	judyDB := db.NewJudyDb(true, logger)
	r := gin.Default()
	return &JudyServer{
		db:     judyDB,
		logger: logger,
		engine: r,
	}, nil
}

func (s *JudyServer) Run(address string) {
	defer s.db.Close()
	s.engine.POST("/hb", func(c *gin.Context) {
		hbResponse := s.handleHeartbeat(c)
		c.ProtoBuf(200, hbResponse)
	})
	s.engine.Run(address)
}

func (s *JudyServer) handleHeartbeat(c *gin.Context) *pb.HeartbeatResponse {
	body, _ := c.GetRawData()
	hb := new(pb.Heartbeat)
	err := proto.Unmarshal(body, hb)
	if err != nil {
		s.logger.Error("Failed to unmarshal heartbeat", zap.Error(err))
	}
	s.logger.Info("Received Heartbeat from", zap.String("client_id", hb.ClientId))
	s.db.AddClient(hb.ClientId, c.ClientIP())
	// TODO: collect tasks
	hbResponse := new(pb.HeartbeatResponse)
	hbResponse.ClientId = hb.ClientId
	return hbResponse
}

func main() {
	judy, err := NewJudyServer()
	if err != nil {
		panic(err)
	}
	judy.Run("0.0.0.0:1337")
}
