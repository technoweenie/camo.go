package fetcher

import (
  "fmt"
  "net/http"
  "github.com/ngmoco/falcore"
  "strconv"
)

type UserAgentFilter interface {
  UserAgent() string
}

type HttpFetcher struct {
  filter UserAgentFilter
  httpClient *http.Client
}

func NewHttpFetcher(filter UserAgentFilter) *HttpFetcher {
  return &HttpFetcher{filter, &http.Client{}}
}

func (fetcher *HttpFetcher) Fetch(req *falcore.Request, url string) *Response {
  clientReq, err := fetcher.buildClientRequest(req, url)
  if err != nil {
    fmt.Printf("Client request error: %s\n", err)
    return Error(req, 500, "Error")
  }

  clientRes, err := fetcher.httpClient.Do(clientReq)
  if err != nil {
    fmt.Printf("Client error: %s\n", err)
    return Error(req, 500, "Error")
  }

  return fetcher.handleResponse(req, clientRes)
}

func (fetcher *HttpFetcher) buildClientRequest(req *falcore.Request, url string) (*http.Request, error) {
  cli, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return nil, err
  }

  accept := req.HttpRequest.Header.Get("Accept")
  if accept == "" {
    accept = "image/*"
  }

  ua := fetcher.filter.UserAgent()
  cli.Header.Set("User-Agent", ua)
  cli.Header.Set("Via", ua)
  cli.Header.Set("X-Content-Type-Options", "nosniff")
  cli.Header.Set("Accept", accept)
  cli.Header.Set("Accept-Encoding", req.HttpRequest.Header.Get("Accept-Encoding"))
  cli.Header.Set("X-Forwarded-For", req.HttpRequest.Header.Get("X-Forwarded-For"))

  return cli, nil
}

func (fetcher *HttpFetcher) handleResponse(req *falcore.Request, clientRes *http.Response) *Response {
  switch clientRes.StatusCode {
  case 200:
    return fetcher.proxyResponse(req, clientRes)
  }

  return Error(req, 500, "responded poorly")
}

func (fetcher *HttpFetcher) proxyResponse(req *falcore.Request, clientRes *http.Response) *Response {
  contentLength := clientRes.Header.Get("Content-Length")
  len, _ := strconv.ParseInt(contentLength, 10, 64)

  res := Respond(req, 200, clientRes.Body)
  res.BodyLength = len

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

  return res
}
