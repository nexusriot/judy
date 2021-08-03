package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"

	"go.uber.org/zap"

	pb "github.com/nexusriot/judy/proto"
)

const idFile = "judy.dat"

type Client struct {
	id         string
	connection *net.Conn
	send       chan []byte
	receive    chan []byte
	logger     *zap.Logger
	session    string
}

func (c *Client) readClientId() (string, error) {
	file, err := os.Open(idFile)
	if err != nil {
		return "", err
	}
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return "", statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	buf := bufio.NewReader(file)
	_, err = buf.Read(bytes)
	id, err := uuid.FromBytes(bytes)
	if err != nil {
		c.logger.Fatal("Couldn't parse uuid")
		file.Close()
		_ = os.Remove(idFile)
		return "", err
	}
	file.Close()
	return id.String(), err
}

func (c *Client) generateClientId() string {
	c.logger.Debug("Generating new client id")
	f, err := os.Create(idFile)
	defer f.Close()
	if err != nil {
		c.logger.Fatal("Couldn't create id file")
	}
	idBytes, err := uuid.New().MarshalBinary()
	w := bufio.NewWriter(f)
	_, err = w.Write(idBytes)
	if err != nil {
		c.logger.Fatal("Couldn't write id file")
	}
	err = w.Flush()
	if err != nil {
		c.logger.Fatal("Couldn't write id file")
	}
	id, err := uuid.FromBytes(idBytes)
	if err != nil {
		c.logger.Fatal("Couldn't parse uuid")
	}
	return id.String()
}

func (c *Client) getClientId() (string, error) {
	res, err := c.readClientId()
	if os.IsNotExist(err) {
		c.logger.Warn("client id not found")
		res = c.generateClientId()
	}
	return res, err
}

func (c *Client) write(message []byte) {
	_, err := (*c.connection).Write(message)
	if err != nil {
		log.Println("Error writing connection", err.Error())
	}
}

func (c *Client) read() (int, []byte) {
	message := make([]byte, 4096)
	length, err := (*c.connection).Read(message)
	if err != nil {
		c.logger.Warn("Failed to read message")
	}
	return length, message
}

func NewClient(logger *zap.Logger) Client {
	return Client{send: make(chan []byte, 4096), receive: make(chan []byte, 4096), logger: logger}
}

func (c *Client) run() {
	client_id, err := c.getClientId()
	c.id = client_id
	c.logger.Info("Client id %s", zap.String("client_id", client_id))

	conn, err := net.Dial("tcp", "localhost:1337")
	defer conn.Close()

	if err != nil {
		c.logger.Error("Connection error:", zap.Error(err))
	}
	c.connection = &conn
	msg1 := &pb.Heartbeat{
		ClientId: c.id,
	}
	data, err := proto.Marshal(msg1)
	if err != nil {
		c.logger.Error("Marshaling error: ", zap.Error(err))
		return
	}
	c.write(data)
	length, message := c.read()
	c.logger.Debug("Got message")
	msg := new(pb.HeartbeatResponse)
	err = proto.Unmarshal(message[:length], msg)
	if err != nil {
		c.logger.Error("Unmarshaling error: ", zap.Error(err))
		return
	}
	c.logger.Debug("message", zap.String("client_id", msg.ClientId))
}

func main() {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if nil != err {
		panic(err)
	}
	logger.Info("Starting client")
	for {
		client := NewClient(logger)
		client.run()
		time.Sleep(10 * time.Second)
	}
}
