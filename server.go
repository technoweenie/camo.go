package camo

import (
	"github.com/ngmoco/falcore"
)

func Server(port int) *falcore.Server {
	userAgent := "camo.go"
	pipe := falcore.NewPipeline()

	methodFilter := NewRequestMethodFilter()
	methodFilter.Allow("GET")
	pipe.Upstream.PushBack(methodFilter)

	emptyFilter := NewSimplePathFilter()
	emptyFilter.AddPath("/")
	emptyFilter.AddPath("/favicon.ico")
	pipe.Upstream.PushBack(emptyFilter)

	viaFilter := NewViaFilter(userAgent)
	pipe.Upstream.PushBack(viaFilter)

	return falcore.NewServer(port, pipe)
}
