package camo

import (
	"github.com/ngmoco/falcore"
)

type Server struct {
	key       string
	Port      int
	UserAgent string
}

func NewServer(key string, port int) *Server {
	return &Server{key, port, "camo.go"}
}

func (server *Server) ListenAndServe() error {
	pipe := falcore.NewPipeline()

	methodFilter := NewRequestMethodFilter()
	methodFilter.Allow("GET")
	pipe.Upstream.PushBack(methodFilter)

	emptyFilter := NewSimplePathFilter()
	emptyFilter.AddPath("/")
	emptyFilter.AddPath("/favicon.ico")
	pipe.Upstream.PushBack(emptyFilter)
	pipe.Upstream.PushBack(NewViaFilter(server.UserAgent))
	pipe.Upstream.PushBack(NewCamoFilter(server.key, server.UserAgent))

	return falcore.NewServer(server.Port, pipe).ListenAndServe()
}
