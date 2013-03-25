package camo

import (
	"./fetcher"
	"fmt"
	"github.com/ngmoco/falcore"
	"net/http"
	"strings"
)

type Fetcher interface {
	Fetch() *fetcher.Response
}

type CamoFilter struct {
	userAgent  *string
	httpClient *http.Client
}

func NewCamoFilter(ua string) *CamoFilter {
	return &CamoFilter{&ua, &http.Client{}}
}

func (filter *CamoFilter) UserAgent() string {
	return *filter.userAgent
}

func (filter *CamoFilter) FilterRequest(req *falcore.Request) *http.Response {
	req.HttpRequest.Header.Del("Cookie")
	url := filter.getUrlFromRequest(req)
	f := fetcher.NewHttpFetcher(filter)
	res := f.Fetch(req, url).HttpResponse()
	return res
}

func (filter *CamoFilter) getUrlFromRequest(req *falcore.Request) string {
	urlPieces := strings.SplitN(req.HttpRequest.URL.Path[1:], "/", 2)
	digest := urlPieces[0]

	query := req.HttpRequest.URL.Query()
	fmt.Println(digest)

	return query.Get("url")
}
