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
	return filter.processUrl(req, filter.getUrlFromRequest(req))
}

func (filter *CamoFilter) getUrlFromRequest(req *falcore.Request) string {
	urlPieces := strings.SplitN(req.HttpRequest.URL.Path[1:], "/", 2)
	digest := urlPieces[0]

	query := req.HttpRequest.URL.Query()
	fmt.Println(digest)
	fmt.Println(query.Get("url"))

	return query.Get("url")
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
	defer clientRes.Body.Close()

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
