package gorudis

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type (
	Client interface {
		Set(string, string) error
		Get(string) (string, error)
		Del(string) error
		AddToSet(string, string) error
		RemoveFromSet(string, string) error
		SetMembers(string) ([]string, error)
		IsSetMember(string, string) (bool, error)
		Ping() error
	}
	client struct {
		conn  net.Conn
		res   chan string
		txmtx sync.Mutex
	}
)

func Init(
	host string,
	port int,
) Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
	c := &client{
		conn:  conn,
		res:   make(chan string, 10),
		txmtx: sync.Mutex{},
	}
	go c.listen()
	return c
}

func (c *client) listen() {
	for {
		buf := make([]byte, 1024)
		n, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
		}
		c.res <- string(buf[:n])
	}
}

func (c *client) Set(key string, value string) error {
	c.txmtx.Lock()
	err := c.send("SET " + key + " " + value)
	if err != nil {
		c.txmtx.Unlock()
		return err
	}
	res := <-c.res
	c.txmtx.Unlock()
	if res != "OK" {
		return fmt.Errorf(res)
	}
	return nil
}

func (c *client) Get(key string) (response string, err error) {
	c.txmtx.Lock()
	err = c.send("GET " + key)
	if err != nil {
		c.txmtx.Unlock()
		return "", err
	}
	response = <-c.res
	c.txmtx.Unlock()
	return response, nil
}

func (c *client) Del(key string) error {
	c.txmtx.Lock()
	err := c.send("DEL " + key)
	if err != nil {
		c.txmtx.Unlock()
		return err
	}
	res := <-c.res
	c.txmtx.Unlock()
	if res != "OK" {
		return fmt.Errorf(res)
	}
	return nil
}

func (c *client) AddToSet(set string, val string) error {
	c.txmtx.Lock()
	err := c.send("SADD " + set + " " + val)
	if err != nil {
		c.txmtx.Unlock()
		return nil
	}
	res := <-c.res
	c.txmtx.Unlock()
	if res != "OK" {
		return fmt.Errorf(res)
	}
	return nil
}

func (c *client) RemoveFromSet(set string, val string) error {
	c.txmtx.Lock()
	err := c.send("SREM " + set + " " + val)
	if err != nil {
		c.txmtx.Unlock()
		return nil
	}
	res := <-c.res
	c.txmtx.Unlock()
	if res != "OK" {
		return fmt.Errorf(res)
	}
	return nil
}

func (c *client) SetMembers(key string) (response []string, err error) {
	c.txmtx.Lock()
	err = c.send("SMEMBERS " + key)
	if err != nil {
		c.txmtx.Unlock()
		return response, err
	}
	res := <-c.res
	c.txmtx.Unlock()
	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		fmt.Println("cries: " + err.Error())
		return response, err
	}
	return response, nil
}

func (c *client) IsSetMember(set string, mem string) (exists bool, err error) {
	c.txmtx.Lock()
	err = c.send("SISMEMBER " + set + " " + mem)
	if err != nil {
		c.txmtx.Unlock()
		return false, err
	}
	res := <-c.res
	c.txmtx.Unlock()
	return res == "true", nil
}

func (c *client) Ping() error {
	c.txmtx.Lock()
	err := c.send("PING")
	if err != nil {
		c.txmtx.Unlock()
		return err
	}
	res := <-c.res
	c.txmtx.Unlock()
	if strings.ToUpper(res) != "PONG" {
		return fmt.Errorf("failed to ping service")
	}
	return nil
}

func (c *client) send(str string) error {
	_, err := c.conn.Write([]byte(str))
	time.Sleep(1000 * time.Nanosecond)
	return err
}
