package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/elevran/relay/log"
)

const (
	refererHeader   = "Referer"
	xURLErrorHeader = "X-URL-Error"
)

type relay struct {
	logger   log.Logger
	frontend *http.Server
	backend  *httptest.Server
}

func newRelay(address string, logger log.Logger) *relay {

	r := &relay{
		frontend: &http.Server{
			Addr:         address,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 4 * time.Second,
		},
		logger: logger,
	}

	r.frontend.Handler = &httputil.ReverseProxy{
		Director: r.rewriteRequest,
	}
	r.backend = httptest.NewUnstartedServer(http.HandlerFunc(r.malformedTargetURL))
	return r
}

func (r *relay) start() error {
	r.backend.Start()
	return r.frontend.ListenAndServe()
}

func (r *relay) shutdown(ctx context.Context) error {
	r.backend.CloseClientConnections()
	r.backend.Close()
	return r.frontend.Shutdown(ctx)
}

func (r *relay) rewriteRequest(req *http.Request) {
	target, err := url.Parse(req.URL.String()[1:]) // skip leading `/`

	if target.Opaque != "" {
		err = errors.New(target.Opaque)
	}

	if err != nil {
		target, _ = url.Parse(r.backend.URL)
		req.Header.Set(refererHeader, req.URL.String())
		req.Header.Set(xURLErrorHeader, err.Error())
	}

	req.URL = target
	req.Host = req.URL.Host
	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}
}

func (r *relay) malformedTargetURL(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("malformed target URL:"))
	w.Write([]byte(req.Referer()))
	w.Write([]byte("\r\n"))
	w.Write([]byte(req.Header.Get(xURLErrorHeader)))
	w.Write([]byte("\r\n"))
}
