package gorudis

import (
	"fmt"
	"net"
	"time"
)

type mockConn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

func newMock(
	expected string,
	result string,
	err error,
) mockConn {
	return &mock{
		expected: expected,
		result:   []byte(result),
		err:      err,
	}
}

type mock struct {
	err      error
	expected string
	result   []byte
}

func (m *mock) Read(b []byte) (n int, err error) {
	if m.err != nil {
		return 1, m.err
	}
	return copy(b, m.result), nil
}

func (m *mock) Write(b []byte) (n int, err error) {
	if string(b) != m.expected {
		return 1, fmt.Errorf("expected %s, got %s", m.expected, string(b))
	}

	return len(b), nil
}

func (m *mock) Close() error                       { return nil }
func (m *mock) LocalAddr() net.Addr                { return nil }
func (m *mock) RemoteAddr() net.Addr               { return nil }
func (m *mock) SetDeadline(t time.Time) error      { return nil }
func (m *mock) SetReadDeadline(t time.Time) error  { return nil }
func (m *mock) SetWriteDeadline(t time.Time) error { return nil }
