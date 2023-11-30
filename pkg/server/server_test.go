package server_test

import (
	"fmt"
	"github.com/zegmic/powserver/pkg/challange"
	"github.com/zegmic/powserver/pkg/db"
	"github.com/zegmic/powserver/pkg/quotes"
	"github.com/zegmic/powserver/pkg/server"
	"net"
	"strings"
	"testing"
)

func TestMessages(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		validate func(string) bool
	}{
		{
			name: "Handle invalid message",
			msg:  "ping",
			validate: func(s string) bool {
				return strings.HasPrefix(s, "requested message is invalid")
			},
		},
		{
			name: "Handle challenge message",
			msg:  "challenge",
			validate: func(s string) bool {
				return !strings.HasPrefix(s, "requested message is invalid")
			},
		},
		{
			name: "Handle solve message",
			msg:  "solve",
			validate: func(s string) bool {
				return !strings.HasPrefix(s, "requested message is invalid")
			},
		},
	}
	port, shutdown := runServer()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%d", port))
			if err != nil {
				t.Error(err)
			}
			defer conn.Close()

			_, err = conn.Write([]byte(test.msg))
			if err != nil {
				t.Error(err)
			}

			buf := make([]byte, 100)
			_, err = conn.Read(buf)
			if err != nil {
				t.Error(err)
			}

			if !test.validate(string(buf)) {
				t.Fail()
			}
		})
	}
	close(shutdown)

}

func runServer() (int, chan interface{}) {
	startup := make(chan interface{})
	shutdown := make(chan interface{})

	s := server.NewServer(&challange.Empty{}, db.New(), quotes.Single{}, 0, shutdown)
	go s.Start(startup)
	<-startup

	return s.Port(), shutdown
}
