package camo

import (
	"fmt"
	"github.com/ngmoco/falcore"
	"net/http"
	"strconv"
	"strings"
)

type CamoFilter struct {
	UserAgent  string
	httpClient *http.Client
}

func NewCamoFilter(ua string) *CamoFilter {
	return &CamoFilter{ua, &http.Client{}}
}

func (filter *CamoFilter) FilterRequest(req *falcore.Request) *http.Response {
	req.HttpRequest.Header.Del("Cookie")
	url := "https://f.cloud.github.com/assets/21/228013/635f6528-8672-11e2-8426-7517f5480715.gif"
	return filter.processUrl(req, url)
}

func (filter *CamoFilter) processUrl(req *falcore.Request, url string) *http.Response {
	clientReq := filter.buildClientRequest(req, url)
	clientRes, err := filter.httpClient.Do(clientReq)
	if err != nil {
		fmt.Printf("Client request error: %s\n", err)
		return falcore.SimpleResponse(req.HttpRequest, 500, nil, "Error")
	}

	return filter.handleResponse(req, clientRes)
}

func (filter *CamoFilter) buildClientRequest(req *falcore.Request, url string) *http.Request {
	cli, _ := http.NewRequest("GET", url, nil)

	accept := req.HttpRequest.Header.Get("Accept")
	if accept == "" {
		accept = "image/*"
	}

	cli.Header.Set("User-Agent", filter.UserAgent)
	cli.Header.Set("Via", filter.UserAgent)
	cli.Header.Set("X-Content-Type-Options", "nosniff")
	cli.Header.Set("Accept", accept)
	cli.Header.Set("Accept-Encoding", req.HttpRequest.Header.Get("Accept-Encoding"))
	cli.Header.Set("X-Forwarded-For", req.HttpRequest.Header.Get("X-Forwarded-For"))

	return cli
}

func (filter *CamoFilter) handleResponse(req *falcore.Request, clientRes *http.Response) *http.Response {
	switch clientRes.StatusCode {
	case 200:
		return filter.proxyResponse(req, clientRes)
	}

	return falcore.SimpleResponse(req.HttpRequest, 500, nil, "responded poorly")
}

func (filter *CamoFilter) proxyResponse(req *falcore.Request, clientRes *http.Response) *http.Response {
	contentLength := clientRes.Header.Get("Content-Length")
	len, _ := strconv.ParseInt(contentLength, 10, 64)

	res := new(http.Response)
	res.StatusCode = 200
	res.ProtoMajor = 1
	res.ProtoMinor = 1
	res.ContentLength = len
	res.Request = req.HttpRequest
	res.Header = make(http.Header)
	res.Body = clientRes.Body

	if contentLength != "" {
		res.Header.Set("Content-Length", contentLength)
	}

	if transfer := clientRes.Header.Get("Transfer-Encoding"); transfer != "" {
		res.Header.Set("Transfer-Encoding", transfer)
	}
	if content := clientRes.Header.Get("Content-Encoding"); content != "" {
		res.Header.Set("Content-Encoding", content)
	}
	cacheControl := clientRes.Header.Get("Cache-Control")
	if cacheControl == "" {
		cacheControl = "public, max-age=31536000"
	}

	res.Header.Set("Content-Type", clientRes.Header.Get("Content-Type"))
	res.Header.Set("Cache-Control", cacheControl)
	res.Header.Set("X-Content-Type-Options", "nosniff")

	return res
}

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
