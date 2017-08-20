package main_test

import (
	"net"
	"net/http"
	"net/http/httptest"
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
	if s.conn == nil {
		return nil
	}

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
	if s.conn == nil {
		return nil
	}

	if s.conn.Scheme == "http" {
		if s.server != nil {
			s.server.Close()
		}
	} else {
		if s.listener != nil {
			err = s.listener.Close()
		}
	}
	return err
}

func TestDialConn(t *testing.T) {
	print := func(a ...interface{}) {}

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

			err = DialConn(&v.conn, defaultTimeout, print)
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

func TestDialConfigs(t *testing.T) {
	print := func(a ...interface{}) {}

	type testItem struct {
		conf          Config
		allowStart    bool
		openConnAfter int
		finishOk      bool
		serverHanlder http.Handler
	}
	testCases := []struct {
		title string
		items []testItem
	}{
		{
			"Should successfully check a single connection.",
			[]testItem{
				{
					Config{Port: 8080, Host: "localhost", Timeout: 5},
					true,
					0,
					true,
					nil,
				},
			},
		},
		{
			"Should successfully check all connections.",
			[]testItem{
				{
					Config{Port: 8080, Host: "localhost", Timeout: 5},
					true,
					0,
					true,
					nil,
				},
				{
					Config{Address: "http://localhost:8081", Timeout: 5},
					true,
					0,
					true,
					nil,
				},
			},
		},
		{
			"Should fail when at least a single connection is not available.",
			[]testItem{
				{
					Config{Port: 8080, Host: "localhost", Timeout: 5},
					true,
					0,
					true,
					nil,
				},
				{
					Config{Port: 8081, Host: "localhost", Timeout: 5},
					false,
					0,
					false,
					nil,
				},
			},
		},
		{
			"Should fail when at least a single connection is not valid.",
			[]testItem{
				{
					Config{Port: 8080, Host: "localhost", Timeout: 5},
					true,
					0,
					true,
					nil,
				},
				{
					Config{Address: "http:/localhost;8081", Timeout: 5},
					false,
					0,
					false,
					nil,
				},
			},
		},
	}

	for _, v := range testCases {
		t.Run(v.title, func(t *testing.T) {
			confs := []Config{}
			finishAllOk := true

			for _, item := range v.items {
				confs = append(confs, item.conf)
				if finishAllOk && !item.finishOk {
					finishAllOk = false
				}

				conn := BuildConn(&item.conf)

				s := NewServer(conn, item.serverHanlder)
				defer s.Close()

				if item.allowStart {
					go func() {
						if item.openConnAfter > 0 {
							time.Sleep(time.Duration(item.openConnAfter) * time.Second)
						}

						if err := s.Start(); err != nil {
							t.Error(err)
						}
					}()
				}
			}

			err := DialConfigs(confs, print)
			if err != nil && finishAllOk {
				t.Errorf("Expected to connect successfully %#v. But got error %v.", confs, err)
				return
			}

			if err == nil && !finishAllOk {
				t.Errorf("Expected to not connect successfully %#v.", confs)
			}
		})
	}
}
