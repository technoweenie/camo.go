package fetcher

import (
  "net/http"
  "io"
  "github.com/ngmoco/falcore"
  "strings"
)

type Response struct {
  Request *falcore.Request
  Status int
  Body io.ReadCloser
  BodyLength int64
  Header http.Header
}

func (fetched *Response) HttpResponse() *http.Response {
  res := new(http.Response)
  res.StatusCode = fetched.Status
  res.ProtoMajor = 1
  res.ProtoMinor = 1
  res.ContentLength = fetched.BodyLength
  res.Request = fetched.Request.HttpRequest
  if fetched.Header == nil {
    res.Header = make(http.Header)
  } else {
    res.Header = fetched.Header
  }
  res.Body = fetched.Body
  return res
}

func Respond(req *falcore.Request, status int, body io.ReadCloser) *Response {
  return &Response{req, status, body, 0, nil}
}

func Error(req *falcore.Request, status int, body string) *Response {
  return &Response{req, status, StringToReadCloser(body), int64(len(body)), nil}
}

func StringToReadCloser(body string) io.ReadCloser {
  return (*fixedResBody)(strings.NewReader(body))
}

type fixedResBody strings.Reader

func (s *fixedResBody) Close() error {
  return nil
}

func (s *fixedResBody) Read(b []byte) (int, error) {
  return (*strings.Reader)(s).Read(b)
}