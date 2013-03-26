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
	userAgent *string
	Digest    DigestCalculator
}

func NewCamoFilter(ua string) *CamoFilter {
	return &CamoFilter{&ua, nil}
}

func (filter *CamoFilter) UserAgent() string {
	return *filter.userAgent
}

func (filter *CamoFilter) FilterRequest(req *falcore.Request) *http.Response {
	req.HttpRequest.Header.Del("Cookie")
	url, err := filter.getUrlFromRequest(req)
	if err != nil {
		return falcore.SimpleResponse(req.HttpRequest, 403, nil, err.Error())
	}

	f := fetcher.NewHttpFetcher(filter)
	res := f.Fetch(req, url).HttpResponse()
	return res
}

func (filter *CamoFilter) getUrlFromRequest(req *falcore.Request) (string, error) {
	urlPieces := strings.SplitN(req.HttpRequest.URL.Path[1:], "/", 2)
	expected := urlPieces[0]
	query := req.HttpRequest.URL.Query()
	url := query.Get("url")
	actual := filter.Digest.Calculate(url)

	if actual == expected {
		return url, nil
	}
	err := fmt.Errorf("checksum mismatch for %s\n%s != %s", url, expected, actual)
	fmt.Printf("%T / %s\n", err, err)
	return url, err
}
