package db

import (
	"database/sql"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const (
	JUDY_DB = "judy.db"
)

type JudyDb struct {
	logger *zap.Logger
	db     *sql.DB
	mx     sync.Mutex
}

func InitDatabase(logger *zap.Logger) *sql.DB {
	logger.Debug("Creating database...")
	_, err := os.Stat(JUDY_DB)
	if os.IsNotExist(err) {
		_, err := os.Create(JUDY_DB)
		if err != nil {
			logger.Fatal("Failed to create db")
		}
	}
	logger.Info("db created")
	// TODO: Migrations
	sqliteDatabase, _ := sql.Open("sqlite3", "./"+JUDY_DB)
	return sqliteDatabase
}

func NewJudyDb(logger *zap.Logger) *JudyDb {
	db := &JudyDb{logger: logger, db: InitDatabase(logger)}
	db.createTables()
	return db
}

func (j *JudyDb) createTables() {
	createClientTableSQL := `CREATE TABLE IF NOT EXISTS client (
		"id" varchar(36) NOT NULL PRIMARY KEY,		
		"ip_addr" varchar(45)
	  );` // SQL Statement for Create Table
	j.logger.Debug("Create client table")
	statement, err := j.db.Prepare(createClientTableSQL) // Prepare SQL Statement
	if err != nil {
		j.logger.Fatal("Failed to create client table")
	}
	statement.Exec() // Execute SQL Statements
	j.logger.Info("client table created")
}

func (j *JudyDb) Close() {
	j.db.Close()
}

func (j *JudyDb) AddClient(uuid string, ip_addr string) {
	j.mx.Lock()
	defer j.mx.Unlock()
	_, err := j.db.Exec("insert or ignore into client (id, ip_addr) values ($1, $2)", uuid, ip_addr)
	if nil != err {
		j.logger.Warn("Error inserting client info", zap.Error(err))
	}
	j.logger.Debug("Client added successfully", zap.String("uuid", uuid), zap.String("ip_addr", ip_addr))
}
