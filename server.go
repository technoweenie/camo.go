package camo

import (
	"github.com/ngmoco/falcore"
)

func Server(port int) *falcore.Server {
	pipe := falcore.NewPipeline()

	methodFilter := NewRequestMethodFilter()
	methodFilter.Allow("GET")
	pipe.Upstream.PushBack(methodFilter)

	emptyFilter := NewSimplePathFilter()
	emptyFilter.AddPath("/")
	emptyFilter.AddPath("/favicon.ico")
	pipe.Upstream.PushBack(emptyFilter)

	return falcore.NewServer(port, pipe)
}
