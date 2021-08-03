package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net"

	"github.com/golang/protobuf/proto"

	db "github.com/nexusriot/judy/pkg/db"
	pb "github.com/nexusriot/judy/proto"
)

func handleConnection(id string, c net.Conn, db *db.JudyDb, logger *zap.Logger) {
	message := make([]byte, 4096)
	length, err := c.Read(message)
	if err != nil {
		logger.Error("Error read form connection: ", zap.Error(err))
	}
	logger.Debug("Got message")
	msg := new(pb.Heartbeat)
	err = proto.Unmarshal(message[:length], msg)
	if err != nil {
		logger.Error("Unmarshaling error: ", zap.Error(err))
		return
	}
	logger.Debug("message", zap.String("client_id", msg.ClientId))
	db.AddClient(msg.ClientId, c.RemoteAddr().String())
	respMsg := &pb.HeartbeatResponse{
		ClientId: msg.ClientId,
	}
	data, err := proto.Marshal(respMsg)
	if err != nil {
		logger.Error("Marshaling error: ", zap.Error(err))
		return
	}
	length, err = c.Write(data)
	if err != nil {
		logger.Error("Sending error: ", zap.Error(err))
	}
}

func main() {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if nil != err {
		panic(err)
	}
	// TODO: make configurable
	judyDB := db.NewJudyDb(logger)
	defer judyDB.Close()
	addr := "0.0.0.0:1337"
	l, err := net.Listen("tcp", addr)
	if nil != err {
		logger.Error("Failed to listen ", zap.String("addr", addr), zap.Error(err))
		return
	}
	defer l.Close()
	// TODO: versioning
	logger.Info("Server v. 0.0.1 started")
	for {
		conn, err := l.Accept()
		if nil != err {
			logger.Error("Error accepting connection", zap.Error(err))
		}
		conn_id := uuid.New().String()
		logger.Info("Accepting request id", zap.String("conn_id", conn_id), zap.String("ip", l.Addr().String()))
		go handleConnection(conn_id, conn, judyDB, logger)
	}
}
