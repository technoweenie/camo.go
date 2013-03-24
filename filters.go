package camo

import (
	"github.com/ngmoco/falcore"
	"net/http"
	"strings"
)

type SimplePathFilter struct {
	Paths map[string]bool
	Body  string
}

func NewSimplePathFilter() *SimplePathFilter {
	return &SimplePathFilter{make(map[string]bool), "OK"}
}

func (filter *SimplePathFilter) AddPath(path string) {
	filter.Paths[path] = true
}

func (filter *SimplePathFilter) FilterRequest(req *falcore.Request) *http.Response {
	if filter.Paths[req.HttpRequest.URL.Path] {
		return falcore.SimpleResponse(req.HttpRequest, 200, nil, filter.Body)
	}
	return nil
}

type RequestMethodFilter struct {
	AllowedMethods map[string]bool
	Body           string
}

func NewRequestMethodFilter() *RequestMethodFilter {
	return &RequestMethodFilter{make(map[string]bool), "Only GET/HEAD allowed"}
}

func (filter *RequestMethodFilter) Allow(method string) {
	if method == "GET" {
		filter.Allow("HEAD")
	}

	filter.AllowedMethods[method] = true
}

func (filter *RequestMethodFilter) FilterRequest(req *falcore.Request) *http.Response {
	if filter.AllowedMethods[req.HttpRequest.Method] {
		return nil
	}
	return falcore.SimpleResponse(req.HttpRequest, 405, nil, filter.Body)
}

type ViaFilter struct {
	UserAgent string
}

func NewViaFilter(ua string) *ViaFilter {
	return &ViaFilter{ua}
}

func (filter *ViaFilter) FilterRequest(req *falcore.Request) *http.Response {
	if strings.HasPrefix(req.HttpRequest.Header.Get("Via"), filter.UserAgent) {
		return falcore.SimpleResponse(req.HttpRequest, 403, nil, "Requesting from self")
	}
	return nil
}
