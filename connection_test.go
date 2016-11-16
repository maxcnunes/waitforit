package main_test

import (
	"reflect"
	"testing"

	. "github.com/maxcnunes/waitforit"
)

func TestBuildConn(t *testing.T) {
	type input struct {
		host     string
		port     int
		fullConn string
	}

	testCases := []struct {
		title    string
		data     input
		expected *Connection
	}{
		{
			"Should create a default connection when only host and port are given",
			input{host: "localhost", port: 80},
			&Connection{Type: "tcp", Scheme: "", Port: 80, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a connection with different host",
			input{host: "localhost", port: 80},
			&Connection{Type: "tcp", Scheme: "", Port: 80, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a connection with different port",
			input{host: "localhost", port: 90},
			&Connection{Type: "tcp", Scheme: "", Port: 90, Host: "localhost", Path: ""},
		},
		{
			"Should ignore the fullConn when the host is given",
			input{host: "localhost", port: 90, fullConn: "tcp://remotehost:10"},
			&Connection{Type: "tcp", Scheme: "", Port: 90, Host: "localhost", Path: ""},
		},
		{
			"Should be able to craete a connection given a fullConn",
			input{fullConn: "tcp://remotehost:10"},
			&Connection{Type: "tcp", Scheme: "", Port: 10, Host: "remotehost", Path: ""},
		},
		{
			"Should be able to create a http connection through the fullConn",
			input{fullConn: "http://localhost"},
			&Connection{Type: "tcp", Scheme: "http", Port: 80, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a https connection through the fullConn",
			input{fullConn: "https://localhost"},
			&Connection{Type: "tcp", Scheme: "https", Port: 443, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a http connection with a path through the fullConn",
			input{fullConn: "https://localhost/cars"},
			&Connection{Type: "tcp", Scheme: "https", Port: 443, Host: "localhost", Path: "/cars"},
		},
		{
			"Should fail when host and full connection are not provided",
			input{},
			nil,
		},
		{
			"Should fail when full connection is not a valid address format",
			input{fullConn: ":/localhost;80"},
			nil,
		},
	}

	for _, v := range testCases {
		conn := BuildConn(v.data.host, v.data.port, v.data.fullConn)
		t.Run(v.title, func(t *testing.T) {
			if !reflect.DeepEqual(conn, v.expected) {
				t.Errorf("Expected to %#v to be deep equal %#v", conn, v.expected)
			}
		})
	}
}
