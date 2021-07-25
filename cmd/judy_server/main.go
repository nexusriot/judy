package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net"

	db "github.com/nexusriot/judy/pkg/db"
)

func handleConnection(id string, c net.Conn, db *db.JudyDb, logger *zap.Logger) {
	db.AddClient(id, c.RemoteAddr().String())

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
		logger.Info("Accepting connection id", zap.String("conn_id", conn_id), zap.String("ip", l.Addr().String()))
		go handleConnection(conn_id, conn, judyDB, logger)
	}
}
