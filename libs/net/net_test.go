package net

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtocolAndAddress(t *testing.T) {

	cases := []struct {
		fullAddr string
		proto    string
		addr     string
	}{
		{
			"tcp://mydomain:80",
			"tcp",
			"mydomain:80",
		},
		{
			"mydomain:80",
			"tcp",
			"mydomain:80",
		},
		{
			"unix://mydomain:80",
			"unix",
			"mydomain:80",
		},
	}

	for _, c := range cases {
		proto, addr := ProtocolAndAddress(c.fullAddr)
		assert.Equal(t, proto, c.proto)
		assert.Equal(t, addr, c.addr)
	}
}

// TestGetFreePort verifies that GetFreePort returns a valid, available port.
func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()
	require.NoError(t, err)
	assert.Greater(t, port, 0)
	assert.LessOrEqual(t, port, 65535)
}

// TestConnect verifies that Connect can establish a TCP connection to a local server.
func TestConnect(t *testing.T) {
	// Start a simple TCP listener on a free port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	addr := ln.Addr().String()

	// Accept connections in the background so Connect does not block.
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()

	conn, err := Connect(fmt.Sprintf("tcp://%s", addr))
	require.NoError(t, err)
	require.NotNil(t, conn)
	conn.Close()
}

// TestConnectFailure verifies that Connect returns an error for an unreachable address.
func TestConnectFailure(t *testing.T) {
	// Port 1 is reserved and should be unreachable in test environments.
	_, err := Connect("tcp://127.0.0.1:1")
	require.Error(t, err)
}

