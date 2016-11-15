package main_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	. "github.com/maxcnunes/waitforit"
)

const regexPort string = `:(\d+)$`

type Server struct {
	conn          *Connection
	listener      net.Listener
	server        *httptest.Server
	serverHandler http.Handler
}

func NewServer(c *Connection, h http.Handler) *Server {
	return &Server{conn: c, serverHandler: h}
}

func (s *Server) Start() (err error) {
	addr := net.JoinHostPort(s.conn.Host, strconv.Itoa(s.conn.Port))
	s.listener, err = net.Listen(s.conn.Type, addr)

	if s.conn.Scheme == "http" {
		s.server = &httptest.Server{
			Listener: s.listener,
			Config:   &http.Server{Handler: s.serverHandler},
		}

		s.server.Start()
	}
	return err
}

func (s *Server) Close() (err error) {
	if s.conn.Scheme == "http" {
		if s.server != nil {
			s.server.Close()
		}
	} else {
		err = s.listener.Close()
	}
	return err
}

func TestBuildConn(t *testing.T) {
	type input struct {
		host     string
		port     int
		fullConn string
	}

	testCases := []struct {
		title    string
		data     input
		expected Connection
	}{
		{
			"Should create a default connection when only host and port are given",
			input{host: "localhost", port: 80},
			Connection{Type: "tcp", Scheme: "", Port: 80, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a connection with different host",
			input{host: "localhost", port: 80},
			Connection{Type: "tcp", Scheme: "", Port: 80, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a connection with different port",
			input{host: "localhost", port: 90},
			Connection{Type: "tcp", Scheme: "", Port: 90, Host: "localhost", Path: ""},
		},
		{
			"Should ignore the fullConn when the host is given",
			input{host: "localhost", port: 90, fullConn: "tcp://remotehost:10"},
			Connection{Type: "tcp", Scheme: "", Port: 90, Host: "localhost", Path: ""},
		},
		{
			"Should be able to craete a connection given a fullConn",
			input{fullConn: "tcp://remotehost:10"},
			Connection{Type: "tcp", Scheme: "", Port: 10, Host: "remotehost", Path: ""},
		},
		{
			"Should be able to create a http connection through the fullConn",
			input{fullConn: "http://localhost"},
			Connection{Type: "tcp", Scheme: "http", Port: 80, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a https connection through the fullConn",
			input{fullConn: "https://localhost"},
			Connection{Type: "tcp", Scheme: "https", Port: 443, Host: "localhost", Path: ""},
		},
		{
			"Should be able to create a http connection with a path through the fullConn",
			input{fullConn: "https://localhost/cars"},
			Connection{Type: "tcp", Scheme: "https", Port: 443, Host: "localhost", Path: "/cars"},
		},
	}

	for _, v := range testCases {
		conn := BuildConn(v.data.host, v.data.port, v.data.fullConn)
		t.Run(v.title, func(t *testing.T) {
			if !reflect.DeepEqual(*conn, v.expected) {
				t.Errorf("Expected to %#v to be deep equal %#v", conn, v.expected)
			}
		})
	}
}

func TestDial(t *testing.T) {
	testCases := []struct {
		title         string
		conn          Connection
		allowStart    bool
		openConnAfter int
		finishOk      bool
		serverHanlder http.Handler
	}{
		{
			"Should successfully check connection that is already available.",
			Connection{Type: "tcp", Scheme: "", Port: 8080, Host: "localhost", Path: ""},
			true,
			0,
			true,
			nil,
		},
		{
			"Should successfully check connection that open before reach the timeout.",
			Connection{Type: "tcp", Scheme: "", Port: 8080, Host: "localhost", Path: ""},
			true,
			2,
			true,
			nil,
		},
		{
			"Should successfully check a HTTP connection that is already available.",
			Connection{Type: "tcp", Scheme: "http", Port: 8080, Host: "localhost", Path: ""},
			true,
			0,
			true,
			nil,
		},
		{
			"Should successfully check a HTTP connection that open before reach the timeout.",
			Connection{Type: "tcp", Scheme: "http", Port: 8080, Host: "localhost", Path: ""},
			true,
			2,
			true,
			nil,
		},
		{
			"Should successfully check a HTTP connection that returns 404 status code.",
			Connection{Type: "tcp", Scheme: "http", Port: 8080, Host: "localhost", Path: ""},
			true,
			0,
			true,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "", 404)
			}),
		},
		{
			"Should fail checking a HTTP connection that returns 500 status code.",
			Connection{Type: "tcp", Scheme: "http", Port: 8080, Host: "localhost", Path: ""},
			true,
			0,
			false,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "", 500)
			}),
		},
	}

	defaultTimeout := 5
	for _, v := range testCases {
		t.Run(v.title, func(t *testing.T) {
			var err error
			s := NewServer(&v.conn, v.serverHanlder)
			defer s.Close()

			if v.allowStart {
				go func() {
					if v.openConnAfter > 0 {
						time.Sleep(time.Duration(v.openConnAfter) * time.Second)
					}

					if err := s.Start(); err != nil {
						t.Error(err)
					}
				}()
			}

			err = Dial(&v.conn, defaultTimeout)
			if err != nil && v.finishOk {
				t.Errorf("Expected to connect successfully %#v. But got error %v.", v.conn, err)
				return
			}

			if err == nil && !v.finishOk {
				t.Errorf("Expected to not connect successfully %#v.", v.conn)
			}
		})
	}
}
