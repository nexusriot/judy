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

type Client struct {
	Id   string     `json:"id"`
	Addr string     `json:"addr"`
	Seen *time.Time `json:"seen"`
}

type Task struct {
	Id        string        `json:"id"`
	ClientId  string        `json:"client_id"`
	Task      string        `json:"task"`
	Created   *time.Time    `json:"created"`
	Completed *sql.NullTime `json:"completed, omitempty"`
}

func (t *Task) String() string {
	out, _ := json.Marshal(t)
	return string(out)
}

func InitDatabase(inMemory bool, logger *zap.Logger) *sql.DB {
	var sqliteDatabase *sql.DB
	if !inMemory {
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
		sqliteDatabase, _ = sql.Open("sqlite3", "./"+JUDY_DB)
	} else {
		sqliteDatabase, _ = sql.Open("sqlite3", ":memory:")
	}

	return sqliteDatabase
}

func getDateTime() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

func NewJudyDb(inMemory bool, logger *zap.Logger) *JudyDb {
	db := &JudyDb{logger: logger, db: InitDatabase(inMemory, logger)}
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
		"client_id" varchar(36) NOT NULL,
		"task" TEXT,
        "created" TIMESTAMP,
        "completed" TIMESTAMP,
        FOREIGN KEY (client_id) REFERENCES client (id)
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
		updateSeenSQL := "UPDATE client SET seen = $1 WHERE id = $2"
		_, err := j.db.Exec(updateSeenSQL, dateTime, uuid)
		if nil != err {
			j.logger.Warn("Error update seen for client", zap.String("uuid", uuid))
		} else {
			j.logger.Debug("Update seen for client succeed", zap.String("uuid", uuid), zap.Time("seen", dateTime))
		}
	}
	j.logger.Debug("Client updated successfully", zap.String("uuid", uuid), zap.String("ip_addr", ip_addr), zap.Time("time", dateTime))
}

func (j *JudyDb) GetCLient(clientId string) (*Client, error) {
	j.mx.Lock()
	defer j.mx.Unlock()
	statement, err := j.db.Prepare("SELECT * FROM client where id = ?")
	if err != nil {
		j.logger.Error("Get client: Failed to prepare statement", zap.Error(err), zap.String("client_id", clientId))
		return nil, err
	}
	var client Client
	row := statement.QueryRow(clientId)
	err = row.Scan(&client.Id, &client.Addr, &client.Seen)
	if err != nil {
		j.logger.Error("Get client: Failed to scan", zap.Error(err), zap.String("client_id", clientId))
		return nil, err
	}
	return &client, nil
}

func (j *JudyDb) ListClients() ([]*Client, error) {
	j.logger.Debug("Listing clients")
	j.mx.Lock()
	defer j.mx.Unlock()
	clients := make([]*Client, 0)
	statement, err := j.db.Prepare("SELECT * FROM client order by seen")
	if err != nil {
		j.logger.Error("List clients: Failed to prepare statement", zap.Error(err))
		return nil, err
	}
	rows, err := statement.Query()
	if err != nil {
		j.logger.Error("List clients: Failed to query statement", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		client := new(Client)
		err = rows.Scan(&client.Id, &client.Addr, &client.Seen)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clients, nil
}

func (j *JudyDb) AddTask(clientId string, task string) string {
	j.mx.Lock()
	defer j.mx.Unlock()
	dateTime := getDateTime()
	taskUUID := uuid.New().String()
	_, err := j.db.Exec("insert into task (id, client_id, task, created) values ($1, $2, $3, $4)", taskUUID, clientId, task, dateTime)
	if nil != err {
		j.logger.Warn("Error inserting client info", zap.Error(err))
	}
	j.logger.Debug("Task added successfully", zap.String("uuid", taskUUID), zap.String("task", task), zap.Time("time", dateTime))
	return taskUUID
}

func (j *JudyDb) GetTask(taskId string) (*Task, error) {
	j.mx.Lock()
	defer j.mx.Unlock()
	statement, err := j.db.Prepare("SELECT * FROM task where id = ?")
	if err != nil {
		j.logger.Error("Get task: Failed to prepare statement", zap.Error(err), zap.String("task_id", taskId))
		return nil, err
	}
	var task Task
	row := statement.QueryRow(taskId)
	err = row.Scan(&task.Id, &task.ClientId, &task.Task, &task.Created, &task.Completed)
	if err != nil {
		j.logger.Error("Get task: Failed to scan", zap.Error(err), zap.String("task_id", taskId))
		return nil, err
	}
	return &task, nil
}

func (j *JudyDb) TakeTasks(clientId string) ([]*Task, error) {
	j.logger.Debug("Taking tasks", zap.String("client_id", clientId))
	tasks := make([]*Task, 0)
	j.mx.Lock()
	defer j.mx.Unlock()
	statement, err := j.db.Prepare("SELECT * FROM task where client_id = ? AND completed is NULL ORDER BY created ASC")
	if err != nil {
		j.logger.Error("Take tasks: Failed to prepare statement", zap.String("client_id", clientId))
		return nil, err
	}
	rows, err := statement.Query(clientId)
	if err != nil {
		j.logger.Error("Take tasks: Failed to query statement", zap.String("client_id", clientId))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := new(Task)
		err = rows.Scan(&task.Id, &task.ClientId, &task.Task, &task.Created, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
