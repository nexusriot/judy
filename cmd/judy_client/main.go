package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"

	pb "github.com/nexusriot/judy/proto"
)

const idFile = "judy.dat"

type Client struct {
	Id      string
	Client  *http.Client
	Logger  *zap.Logger
	Session string
}

//func (c *Client) execCommand(command string) (string, error) {
//	cmd := exec.Command(command)
//	//cmd.Stdin = strings.NewReader()
//	var out bytes.Buffer
//	cmd.Stdin = &out
//	err := cmd.Run()
//	if err != nil {
//		c.Logger.Error("Failed to execute command", zap.Error(err))
//		return "", err
//	}
//	output := out.String()
//	c.Logger.Info("output", zap.String("output", output))
//	return output, err
//}

func (c *Client) readClientId() (string, error) {
	file, err := os.Open(idFile)
	if err != nil {
		return "", err
	}
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return "", statsErr
	}

	var size = stats.Size()
	bytes := make([]byte, size)

	buf := bufio.NewReader(file)
	_, err = buf.Read(bytes)
	id, err := uuid.FromBytes(bytes)
	if err != nil {
		c.Logger.Fatal("Couldn't parse uuid")
		file.Close()
		_ = os.Remove(idFile)
		return "", err
	}
	file.Close()
	return id.String(), err
}

func (c *Client) generateClientId() string {
	c.Logger.Debug("Generating new client id")
	f, err := os.Create(idFile)
	defer f.Close()
	if err != nil {
		c.Logger.Fatal("Couldn't create id file")
	}
	idBytes, err := uuid.New().MarshalBinary()
	w := bufio.NewWriter(f)
	_, err = w.Write(idBytes)
	if err != nil {
		c.Logger.Fatal("Couldn't write id file")
	}
	err = w.Flush()
	if err != nil {
		c.Logger.Fatal("Couldn't write id file")
	}
	id, err := uuid.FromBytes(idBytes)
	if err != nil {
		c.Logger.Fatal("Couldn't parse uuid")
	}
	return id.String()
}

func (c *Client) getClientId() (string, error) {
	res, err := c.readClientId()
	if os.IsNotExist(err) {
		c.Logger.Warn("client id not found")
		res = c.generateClientId()
	}
	return res, err
}

func (c *Client) Heartbeat() {
	hb := pb.Heartbeat{ClientId: c.Id}
	data, err := proto.Marshal(&hb)
	responseBody := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:1337/hb", responseBody)
	if err != nil {
		//  TODO:
	}
	req.Header.Set("User-Agent", "ncl")
	req.Header.Set("content-type", "application/x-protobuf")
	res, getErr := c.Client.Do(req)
	if getErr != nil {
		// TODO:
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		// TODO:
	}
	hbResponse := new(pb.HeartbeatResponse)
	err = proto.Unmarshal(body, hbResponse)
	if err != nil {
		c.Logger.Error("Failed to unmarshal heartbeat", zap.Error(err))
	}
	c.Logger.Info(hbResponse.ClientId)
}

func (c *Client) run() {
	client_id, err := c.getClientId()
	if err != nil {
		c.Logger.Panic("Failed to get client id")
	}
	c.Id = client_id
	c.Logger.Info("Client id %s", zap.String("client_id", client_id))
	for {
		c.Heartbeat()
		time.Sleep(10 * time.Second)
	}
}

func NewClient(logger *zap.Logger) Client {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	return Client{Client: &client, Logger: logger}
}

func main() {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if nil != err {
		panic(err)
	}
	logger.Info("Starting client")
	client := NewClient(logger)
	client.run()
}
