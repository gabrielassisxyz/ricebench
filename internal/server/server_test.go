package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func builtAssets() fstest.MapFS {
	return fstest.MapFS{"index.html": {Data: []byte("<!doctype html><title>RiceBench</title>")}}
}

func TestNewRejectsUnbuiltFrontend(t *testing.T) {
	_, err := New(fstest.MapFS{})
	if !errors.Is(err, ErrAssetsMissing) {
		t.Fatalf("New with no index.html: got %v, want ErrAssetsMissing", err)
	}
}

func TestNewServesIndex(t *testing.T) {
	handler, err := New(builtAssets())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, http.StatusOK)
	}
	body, _ := io.ReadAll(recorder.Body)
	if !strings.Contains(string(body), "RiceBench") {
		t.Fatalf("body did not come from the embedded index: %q", body)
	}
}

func TestStartupNoticeWarnsOnlyWhenExposed(t *testing.T) {
	cases := []struct {
		addr     string
		exposed  bool
		wantHost string
	}{
		{addr: "127.0.0.1:7391", exposed: false, wantHost: "127.0.0.1:7391"},
		{addr: "localhost:7391", exposed: false, wantHost: "localhost:7391"},
		{addr: "[::1]:7391", exposed: false, wantHost: "[::1]:7391"},
		{addr: "0.0.0.0:7391", exposed: true, wantHost: "localhost:7391"},
		{addr: ":7391", exposed: true, wantHost: "localhost:7391"},
		{addr: "192.168.1.20:7391", exposed: true, wantHost: "192.168.1.20:7391"},
	}

	for _, tc := range cases {
		t.Run(tc.addr, func(t *testing.T) {
			notice := StartupNotice(tc.addr)
			if !strings.Contains(notice, "http://"+tc.wantHost+"/") {
				t.Errorf("notice did not advertise %s: %q", tc.wantHost, notice)
			}
			if warned := strings.Contains(notice, "WARNING"); warned != tc.exposed {
				t.Errorf("exposure warning: got %v, want %v in %q", warned, tc.exposed, notice)
			}
		})
	}
}
