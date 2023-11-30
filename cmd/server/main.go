package main

import (
	"github.com/zegmic/powserver/pkg/challange"
	"github.com/zegmic/powserver/pkg/db"
	"github.com/zegmic/powserver/pkg/quotes"
	"github.com/zegmic/powserver/pkg/server"
)

func main() {
	shutdown := make(chan interface{})
	q := quotes.New()
	s := server.NewServer(&challange.Hash{}, db.New(), q, 8080, shutdown)
	startup := make(chan interface{})
	go s.Start(startup)
	<-startup
	select {}
}
