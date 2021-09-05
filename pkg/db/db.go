package db

import (
	"database/sql"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

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

type Task struct {
	Id        string        `json:"id"`
	Task      string        `json:"task"`
	Created   *time.Time    `json:"created"`
	Taken     *sql.NullTime `json:"taken, omitempty"`
	Completed *sql.NullTime `json:"completed, omitempty"`
}

func (t *Task) String() string {
	out, _ := json.Marshal(t)
	return string(out)
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

func getDateTime() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

func NewJudyDb(logger *zap.Logger) *JudyDb {
	db := &JudyDb{logger: logger, db: InitDatabase(logger)}
	db.createTables()
	return db
}

func (j *JudyDb) createTables() {
	createClientTableSQL := `CREATE TABLE IF NOT EXISTS client (
		"id" varchar(36) NOT NULL PRIMARY KEY,		
		"addr" varchar(45),
		"seen" TIMESTAMP
	  );` // SQL Statement for Create Table
	createTasksTableSQL := `CREATE TABLE IF NOT EXISTS task (
		"id" varchar(36) NOT NULL PRIMARY KEY,
		"task" TEXT,
        "created" TIMESTAMP,
        "taken" TIMESTAMP,
        "completed" TIMESTAMP 
	  );`
	j.logger.Debug("Create client table")
	statement, err := j.db.Prepare(createClientTableSQL) // Prepare SQL Statement
	if err != nil {
		j.logger.Fatal("Failed to create client table")
	}
	statement.Exec() // Execute SQL Statements
	j.logger.Info("client table created")
	statement, err = j.db.Prepare(createTasksTableSQL) // Prepare SQL Statement
	if err != nil {
		j.logger.Fatal("Failed to create tasks table")
	}
	statement.Exec() // Execute SQL Statements
	j.logger.Info("task table created")
}

func (j *JudyDb) Close() {
	j.db.Close()
}

func (j *JudyDb) AddClient(uuid string, ip_addr string) {
	j.mx.Lock()
	defer j.mx.Unlock()
	dateTime := getDateTime()
	fetch, err := j.db.Exec("insert or ignore into client (id, addr, seen) values ($1, $2, $3)", uuid, ip_addr, dateTime)
	if nil != err {
		j.logger.Warn("Error inserting client info", zap.Error(err))
	}
	affected, _ := fetch.RowsAffected()
	if affected == 0 {
		setTakenSQL := "UPDATE client SET seen = $1 WHERE id = $2"
		_, err := j.db.Exec(setTakenSQL, dateTime, uuid)
		if nil != err {
			j.logger.Warn("Error update seen for client", zap.String("uuid", uuid))
		} else {
			j.logger.Debug("Update seen for client succeed", zap.String("uuid", uuid), zap.Time("seen", dateTime))
		}
	}
	j.logger.Debug("Client updated successfully", zap.String("uuid", uuid), zap.String("ip_addr", ip_addr), zap.Time("time", dateTime))
}

func (j *JudyDb) AddTask(task string) {
	j.mx.Lock()
	defer j.mx.Unlock()
	dateTime := getDateTime()
	taskUUID := uuid.New().String()
	_, err := j.db.Exec("insert or ignore into task (id, task, created) values ($1, $2, $3)", taskUUID, task, dateTime)
	if nil != err {
		j.logger.Warn("Error inserting client info", zap.Error(err))
	}
	j.logger.Debug("Task added successfully", zap.String("uuid", taskUUID), zap.String("task", task), zap.Time("time", dateTime))
}

func (j *JudyDb) TakeTask() (*Task, error) {
	j.mx.Lock()
	defer j.mx.Unlock()
	dateTime := getDateTime()
	statement := "SELECT * FROM task where taken is NULL ORDER BY created ASC LIMIT 1"
	var task Task
	row := j.db.QueryRow(statement, 3)
	err := row.Scan(&task.Id, &task.Task, &task.Created, &task.Taken, &task.Completed)
	if err == nil {
		setTakenSQL := "UPDATE task SET taken = $1 WHERE id = $2"
		_, err := j.db.Exec(setTakenSQL, dateTime, task.Id)
		if nil != err {
			j.logger.Warn("Error set task taking", zap.String("uuid", task.Id))
		} else {
			j.logger.Debug("TakeTask succeed", zap.String("uuid", task.Id))
		}
		return &task, nil
	}
	return nil, err
}
