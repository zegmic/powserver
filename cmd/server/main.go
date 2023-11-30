package main

import (
	"github.com/zegmic/powserver/pkg/challange"
	"github.com/zegmic/powserver/pkg/db"
	"github.com/zegmic/powserver/pkg/server"
)

func main() {
	shutdown := make(chan interface{})

	s := server.NewServer(&challange.Hash{}, db.New(), 8080, shutdown)
	startup := make(chan interface{})
	go s.Start(startup)
	<-startup
	select {}
}
