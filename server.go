package camo

import (
	"github.com/ngmoco/falcore"
)

type Server struct {
	Camo *CamoFilter
}

func NewServer() *Server {
	return &Server{NewCamoFilter("camo.go")}
}

func (server *Server) ListenAndServe(port int) error {
	pipe := falcore.NewPipeline()

	methodFilter := NewRequestMethodFilter()
	methodFilter.Allow("GET")
	pipe.Upstream.PushBack(methodFilter)

	emptyFilter := NewSimplePathFilter()
	emptyFilter.AddPath("/")
	emptyFilter.AddPath("/favicon.ico")
	pipe.Upstream.PushBack(emptyFilter)
	pipe.Upstream.PushBack(NewViaFilter(server.Camo.UserAgent()))
	pipe.Upstream.PushBack(server.Camo)

	return falcore.NewServer(port, pipe).ListenAndServe()
}

func (server *Server) SetDigest(digest DigestCalculator) {
	server.Camo.Digest = digest
}

func (server *Server) SetDigestKey(key string) {
	server.SetDigest(NewDigest(key))
}
