package main_test

import (
	"testing"

	. "github.com/maxcnunes/waitforit"
)

func TestBuildConn(t *testing.T) {
	type input struct {
		proto   string
		host    string
		port    int
		address string
	}

	type expected struct {
		netType string
		host    string
		address string
	}

	testCases := []struct {
		title    string
		data     input
		expected *expected
	}{
		{
			"Should create a default connection when only host and port are given",
			input{host: "localhost", port: 80},
			&expected{netType: "tcp", host: "localhost:80", address: "tcp://localhost:80"},
		},
		{
			"Should be able to create a connection with different port",
			input{host: "localhost", port: 90},
			&expected{netType: "tcp", host: "localhost:90", address: "tcp://localhost:90"},
		},
		{
			"Should ignore protocol, host or port when an address is given",
			input{host: "localhost", port: 90, address: "tcp://remotehost:10"},
			&expected{netType: "tcp", host: "remotehost:10", address: "tcp://remotehost:10"},
		},
		{
			"Should be able to craete a connection given a address",
			input{address: "tcp://remotehost:10"},
			&expected{netType: "tcp", host: "remotehost:10", address: "tcp://remotehost:10"},
		},
		{
			"Should be able to create a http connection through the address",
			input{address: "http://localhost"},
			&expected{netType: "tcp", host: "localhost:80", address: "http://localhost:80"},
		},
		{
			"Should be able to create a https connection through the address",
			input{address: "https://localhost"},
			&expected{netType: "tcp", host: "localhost:443", address: "https://localhost:443"},
		},
		{
			"Should support address to https with custom port",
			input{address: "https://localhost:444"},
			&expected{netType: "tcp", host: "localhost:444", address: "https://localhost:444"},
		},
		{
			"Should be able to create a http connection with a path through the address",
			input{address: "https://localhost/cars"},
			&expected{netType: "tcp", host: "localhost:443", address: "https://localhost:443/cars"},
		},
		{
			"Should be able to create a http connection with a path with inner paths",
			input{address: "http://backend:8182/backend/tunnel/tunnel.nocache.js"},
			&expected{
				netType: "tcp",
				host:    "backend:8182",
				address: "http://backend:8182/backend/tunnel/tunnel.nocache.js",
			},
		},
		{
			"Should be able to create a ipv6 connection",
			input{address: "http://[2001:41d0:8:6a52:298:2dff:fef3:8ce1]:8182/cars"},
			&expected{
				netType: "tcp",
				host:    "[2001:41d0:8:6a52:298:2dff:fef3:8ce1]:8182",
				address: "http://[2001:41d0:8:6a52:298:2dff:fef3:8ce1]:8182/cars",
			},
		},
		{
			"Should be able to create a ipv6 connection without a provided scheme",
			input{address: "[2001:41d0:8:6a52:298:2dff:fef3:8ce1]:8182/cars"},
			&expected{
				netType: "tcp",
				host:    "[2001:41d0:8:6a52:298:2dff:fef3:8ce1]:8182",
				address: "tcp://[2001:41d0:8:6a52:298:2dff:fef3:8ce1]:8182/cars",
			},
		},
		{
			"Should fail when host and full connection are not provided",
			input{},
			nil,
		},
		{
			"Should fail when full connection is not a valid address format",
			input{address: ":/localhost;80"},
			nil,
		},
	}

	for _, v := range testCases {
		t.Run(v.title, func(t *testing.T) {
			cfg := &Config{
				Host:    v.data.host,
				Port:    v.data.port,
				Address: v.data.address,
			}

			conn, err := BuildConn(cfg)
			if v.expected == nil {
				if conn == nil {
					return
				}

				t.Fatalf("Expected connection build to fail, instead it successed with network type %s and url %s", conn.NetworkType, conn.URL)
			}

			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, "network type", conn.NetworkType, v.expected.netType)
			assertEqual(t, "host", conn.URL.Host, v.expected.host)
			assertEqual(t, "address", conn.URL.String(), v.expected.address)
		})
	}
}

func assertEqual(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %s %#v to be deep equal to %#v", name, a, b)
	}
}
