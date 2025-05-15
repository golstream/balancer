package httputils

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetWithCtx_Success(t *testing.T) {
	// Поднимаем тестовый сервер
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", MIMEApplicationJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"ok"}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ctx := context.Background()
	resp, err := GetWithCtx(ctx, server.URL, nil, nil, nil, 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{"message":"ok"}`
	if !bytes.Equal(resp.body, []byte(expected)) {
		t.Errorf("unexpected body: got %s, want %s", string(resp.body), expected)
	}
}

func TestGetWithCtx_Timeout(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ctx := context.Background()
	_, err := GetWithCtx(ctx, server.URL, nil, nil, nil, 1*time.Second)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Errorf("expected context timeout error, got: %v", err)
	}
}

func TestDoRequestWithCtx_InvalidURL(t *testing.T) {
	ctx := context.Background()
	_, err := doRequestWithCtx(ctx, http.MethodGet, ":", nil, nil, nil, nil, time.Second)
	if err == nil {
		t.Fatal("expected error for malformed URL")
	}
}

func TestDoRequestWithCtx_NilResponse(t *testing.T) {
	// используем сервер, который сразу закрывается
	server := httptest.NewServer(nil)
	server.Close()

	ctx := context.Background()
	_, err := doRequestWithCtx(ctx, http.MethodGet, server.URL, nil, nil, nil, nil, time.Second)
	if err == nil {
		t.Fatal("expected error due to closed server")
	}
}

func TestDoRequestWithCtx_SetHeadersAndCookies(t *testing.T) {
	called := false

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Header") != "123" {
			t.Errorf("missing or incorrect custom header")
		}
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value != "abc" {
			t.Errorf("missing or incorrect cookie")
		}
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ctx := context.Background()
	headers := map[string]string{
		"X-Test-Header": "123",
	}
	cookies := []*http.Cookie{
		{Name: "session_id", Value: "abc"},
	}
	_, err := GetWithCtx(ctx, server.URL, nil, headers, cookies, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Errorf("handler was not called")
	}
}
