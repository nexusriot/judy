package main

import (
	"go.uber.org/zap"
	"net"

	"github.com/golang/protobuf/proto"

	db "github.com/nexusriot/judy/pkg/db"
	pb "github.com/nexusriot/judy/proto"
)

// TODO: temp
func createTask(db *db.JudyDb) {
	db.AddTask("ls")
}

func handleManagementConnection(c net.Conn, db *db.JudyDb, logger *zap.Logger) {
	message := make([]byte, 16384)
	length, err := c.Read(message)
	if err != nil {
		logger.Error("Error read form connection: ", zap.Error(err))
	}
	logger.Debug("Got message")
	msg := new(pb.CommandRequest)
	err = proto.Unmarshal(message[:length], msg)
	if err != nil {
		logger.Error("Unmarshaling error: ", zap.Error(err))
		return
	}
	logger.Debug("message", zap.String("command", msg.Command))
	// TODO: get response

	respMsg := &pb.CommandResponse{
		Response: "hello",
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

func handleConnection(c net.Conn, db *db.JudyDb, logger *zap.Logger) {
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
	task, _ := db.TakeTask()
	logger.Debug("task", zap.String("task_data", task.String()))
	respMsg := &pb.HeartbeatResponse{
		ClientId: msg.ClientId,
		Command:  task.Task,
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
	sync := make(chan struct{})
	judyDB := db.NewJudyDb(logger)
	defer judyDB.Close()
	addr := "0.0.0.0:1337"
	createTask(judyDB)
	l, err := net.Listen("tcp", addr)
	if nil != err {
		logger.Error("Failed to listen ", zap.String("addr", addr), zap.Error(err))
		return
	}
	defer l.Close()
	addr_man := "0.0.0.0:1338"
	lman, err := net.Listen("tcp", addr_man)
	if nil != err {
		logger.Error("Failed to listen management", zap.String("addr", addr_man), zap.Error(err))
		return
	}
	defer lman.Close()
	// TODO: versioning
	logger.Info("Server v. 0.0.1 started")
	// handle heartbeat connection
	go func() {
		for {
			conn, err := l.Accept()
			if nil != err {
				logger.Error("Error accepting connection", zap.Error(err))
			}
			logger.Info("Accepting request", zap.String("ip", l.Addr().String()))
			go handleConnection(conn, judyDB, logger)
		}
	}()
	// handle management connections
	go func() {
		for {
			conn, err := lman.Accept()
			if nil != err {
				logger.Error("Error accepting management connection", zap.Error(err))
			}
			logger.Info("Accepting request", zap.String("ip", l.Addr().String()))
			go handleManagementConnection(conn, judyDB, logger)
		}
	}()
	<-sync
}
