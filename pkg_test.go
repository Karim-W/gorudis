package gorudis

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testClient(
	expected string,
	result string,
	err error,
) Client {
	c := &client{
		conn:  newMock(expected, result, err),
		res:   make(chan string, 5),
		txmtx: sync.Mutex{},
	}
	go c.listen()
	return c
}

func TestPing(t *testing.T) {
	cl := testClient("PING", "pong", nil)
	err := cl.Ping()
	assert.Nil(t, err)
}

func TestSet(t *testing.T) {
	cl := testClient("SET Foo Bar", "OK", nil)
	err := cl.Set("Foo", "Bar")
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	cl := testClient("GET Foo", "Bar", nil)
	res, err := cl.Get("Foo")
	assert.Nil(t, err)
	assert.Equal(t, "Bar", res)
}

func TestDel(t *testing.T) {
	cl := testClient("DEL Foo", "OK", nil)
	err := cl.Del("Foo")
	assert.Nil(t, err)
}

func TestSetError(t *testing.T) {
	cl := testClient("SET Foo Bar", "ERR", nil)
	err := cl.Set("Foo", "Bar")
	assert.NotNil(t, err)
}

func TestAddToSet(t *testing.T) {
	cl := testClient("SADD Foo Bar", "OK", nil)
	err := cl.AddToSet("Foo", "Bar")
	assert.Nil(t, err)
}

func TestGetSMemebers(t *testing.T) {
	cl := testClient("SMEMBERS Foo", `["bar","baz"]`, nil)
	res, err := cl.SetMembers("Foo")
	assert.Nil(t, err)
	assert.Equal(t, []string{"bar", "baz"}, res)
}
