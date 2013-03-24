package camo

import (
	"github.com/ngmoco/falcore"
)

type Server struct {
	Port      int
	UserAgent string
}

func NewServer(port int) *Server {
	return &Server{port, "camo.go"}
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
	pipe.Upstream.PushBack(NewCamoFilter(server.UserAgent))

	return falcore.NewServer(server.Port, pipe).ListenAndServe()
}
