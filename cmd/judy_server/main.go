package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/nexusriot/judy/pkg/db"
	pb "github.com/nexusriot/judy/proto"
)

// TODO: temp
//func createTask(db *db.JudyDb) {
//	db.AddTask("ls")
//}

func handleHeartbeat(db *db.JudyDb, c *gin.Context, logger *zap.Logger) *pb.HeartbeatResponse {
	body, _ := c.GetRawData()
	hb := new(pb.Heartbeat)
	err := proto.Unmarshal(body, hb)
	if err != nil {
		logger.Error("Failed to unmarshal heartbeat", zap.Error(err))
	}
	logger.Info("Received Heartbeat from", zap.String("client_id", hb.ClientId))
	db.AddClient(hb.ClientId, c.ClientIP())
	// TODO: collect tasks
	hbResponse := new(pb.HeartbeatResponse)
	hbResponse.ClientId = hb.ClientId
	return hbResponse
}

func main() {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if nil != err {
		panic(err)
	}
	// TODO: make configurable
	//sync := make(chan struct{})
	judyDB := db.NewJudyDb(true, logger)
	defer judyDB.Close()
	r := gin.Default()
	r.POST("/hb", func(c *gin.Context) {
		hbResponse := handleHeartbeat(judyDB, c, logger)
		c.ProtoBuf(200, hbResponse)
	})
	r.Run("0.0.0.0:1337")
}
