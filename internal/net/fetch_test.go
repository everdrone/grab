package net

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	retriesCount := 1

	tests := []struct {
		Name         string
		UseServerURL bool
		URL          string
		Handler      func(w http.ResponseWriter, r *http.Request)
		Options      *FetchOptions
		Want         string
		HasError     bool
	}{
		{
			Name:         "200",
			UseServerURL: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello, world!"))
			},
			Options: &FetchOptions{
				Retries: 1,
				Timeout: 10,
				Headers: map[string]string{},
			},
			Want:     "Hello, world!",
			HasError: false,
		},
		{
			Name:         "200 with headers",
			UseServerURL: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.Header.Get("X-Foo")))
			},
			Options: &FetchOptions{
				Retries: 1,
				Timeout: 10,
				Headers: map[string]string{
					"X-Foo": "bar",
				},
			},
			Want:     "bar",
			HasError: false,
		},
		{
			Name:         "500",
			UseServerURL: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			Options: &FetchOptions{
				Retries: 3,
				Timeout: 10,
				Headers: map[string]string{},
			},
			Want:     "",
			HasError: true,
		},
		{
			Name:         "invalid url",
			UseServerURL: false,
			URL:          "1http://example.com",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			Options: &FetchOptions{
				Retries: 3,
				Timeout: 10,
				Headers: map[string]string{},
			},
			Want:     "",
			HasError: true,
		},
		{
			Name:         "with retries",
			UseServerURL: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				if retriesCount <= 2 {
					retriesCount++
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprint(retriesCount)))
			},
			Options: &FetchOptions{
				Retries: 3,
				Timeout: 10,
				Headers: map[string]string{},
			},
			Want:     "3",
			HasError: false,
		},
		{
			Name:         "with negative retries",
			UseServerURL: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				retriesCount = 1
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprint(retriesCount)))
			},
			Options: &FetchOptions{
				Retries: -2,
				Timeout: 10,
				Headers: map[string]string{},
			},
			Want:     "1",
			HasError: false,
		},
		{
			Name:         "with negative timeout",
			UseServerURL: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			Options: &FetchOptions{
				Retries: 1,
				Timeout: -10,
				Headers: map[string]string{},
			},
			Want:     "",
			HasError: false,
		},
		{
			Name:         "with timeout",
			UseServerURL: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(300 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			},
			Options: &FetchOptions{
				Retries: 1,
				Timeout: 10,
			},
			Want:     "",
			HasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(test.Handler))
			defer server.Close()

			if test.UseServerURL {
				test.URL = server.URL
			}

			body, err := Fetch(test.URL, test.Options)
			if err != nil && !test.HasError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && test.HasError {
				t.Errorf("expected error, but got none")
			}

			if body != test.Want {
				t.Errorf("got: %v, want: %v", body, test.Want)
			}
		})
	}
}
