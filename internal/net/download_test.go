package net

import (
	"crypto/sha256"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/everdrone/grab/internal/utils"
	tu "github.com/everdrone/grab/testutils"
	"github.com/spf13/afero"
	"golang.org/x/exp/slices"
)

func TestDownload(t *testing.T) {
	initialWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(initialWd)
	}()

	root := tu.GetOSRoot()

	defaultOptions := &FetchOptions{
		Retries: 1,
		Timeout: 3000,
		Headers: map[string]string{},
	}

	tests := []struct {
		Name      string
		Path      []string
		CustomURL string
		Dest      string
		Handler   func(fs http.FileSystem) http.Handler
		HasError  bool
		Options   *FetchOptions
	}{
		{
			Name: "simple file",
			Handler: func(fs http.FileSystem) http.Handler {
				ts := http.FileServer(fs)
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ts.ServeHTTP(w, r)
				})
			},
			Path:     []string{"net", "file.txt"},
			Dest:     filepath.Join(root, "net", "file.txt.dl"),
			HasError: false,
			Options:  defaultOptions,
		},
		{
			Name: "zero or negative retries",
			Handler: func(fs http.FileSystem) http.Handler {
				ts := http.FileServer(fs)
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ts.ServeHTTP(w, r)
				})
			},
			Path:     []string{"net", "file.txt"},
			Dest:     filepath.Join(root, "net", "file.txt.dl"),
			HasError: false,
			Options: &FetchOptions{
				Retries: 1,
				Timeout: -100,
				Headers: map[string]string{},
			},
		},
		{
			Name: "zero or negative retries",
			Handler: func(fs http.FileSystem) http.Handler {
				ts := http.FileServer(fs)
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ts.ServeHTTP(w, r)
				})
			},
			Path:     []string{"net", "file.txt"},
			Dest:     filepath.Join(root, "net", "file.txt.dl"),
			HasError: false,
			Options: &FetchOptions{
				Retries: -201,
				Timeout: 3000,
				Headers: map[string]string{},
			},
		},
		{
			Name: "check headers",
			Handler: func(fs http.FileSystem) http.Handler {
				ts := http.FileServer(fs)
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Header.Get("foo") == "bar" {
						ts.ServeHTTP(w, r)
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				})
			},
			Path:     []string{"net", "file.txt"},
			Dest:     filepath.Join(root, "net", "file.txt.dl"),
			HasError: false,
			Options: &FetchOptions{
				Retries: 1,
				Timeout: 3000,
				Headers: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			Name: "404",
			Handler: func(fs http.FileSystem) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				})
			},
			Path:     []string{"net", "file.txt"},
			Dest:     filepath.Join(root, "net", "file.txt.dl"),
			HasError: true,
			Options:  defaultOptions,
		},
		{
			Name: "times out",
			Handler: func(fs http.FileSystem) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(time.Millisecond * 400)
					w.WriteHeader(http.StatusOK)
				})
			},
			Path:     []string{"net", "file.txt"},
			Dest:     filepath.Join(root, "net", "file.txt.dl"),
			HasError: true,
			Options:  &FetchOptions{Retries: 1, Timeout: 300},
		},
		{
			Name: "read error",
			Handler: func(fs http.FileSystem) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// tells the client that the content has a length but it does not give any content back
					w.Header().Set("Content-Length", "1")
				})
			},
			Path:     []string{"net", "file.txt"},
			Dest:     filepath.Join(root, "net", "file.txt.dl"),
			HasError: true,
			Options:  defaultOptions,
		},
		{
			Name: "invalid url",
			Handler: func(fs http.FileSystem) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
			},
			CustomURL: "1http://example.com",
			Dest:      filepath.Join(root, "net", "file.txt.dl"),
			HasError:  true,
			Options:   defaultOptions,
		},
		// {
		// 	Name: "invalid url",
		// },
		// FIXME: afero's memMapFs.Create doesn't return an error!
		// see: https://github.com/spf13/afero/blob/2a70f2bb2db1524bf2aa3ca0cfebefa8d6367b7b/memmap_test.go#L65
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(tc *testing.T) {
			// start fresh for each test case
			utils.Fs, utils.AFS, utils.Wd = tu.SetupMemMapFs(root)
			utils.AFS.WriteFile(filepath.Join(root, "net", "file.txt"), []byte("binary"), os.ModePerm)
			httpFs := afero.NewHttpFs(utils.Fs)

			// start the server
			ts := httptest.NewServer(tt.Handler(httpFs.Dir(root)))
			defer ts.Close()

			fileURL := strings.Join(tt.Path, "/")
			if tt.CustomURL != "" {
				fileURL = tt.CustomURL
			} else {
				base, err := url.Parse(ts.URL)
				if err != nil {
					tc.Fatal(err)
				}

				path, err := url.Parse(fileURL)
				if err != nil {
					tc.Fatal(err)
				}

				resolved := base.ResolveReference(path)

				fileURL = resolved.String()
			}

			err := Download(fileURL, tt.Dest, tt.Options)
			if (err != nil) != tt.HasError {
				tc.Errorf("got: %v, want: %v", err, tt.HasError)
			}
			if err == nil && tt.HasError {
				tc.Errorf("expected error, but got none")
			}

			if !tt.HasError {
				// checksum the new file
				h1, err := getHash(tt.Dest)
				if err != nil {
					tc.Fatalf("unexpected error: %v", err)
				}

				// checksum the original file
				p := filepath.Join(tt.Path...)
				h2, err := getHash(filepath.Join(root, p))
				if err != nil {
					tc.Fatalf("unexpected error: %v", err)
				}

				if !slices.Equal(h1, h2) {
					tc.Errorf("got %v, want %v", string(h1), string(h2))
				}
			}
		})
	}
}

func getHash(filename string) ([]byte, error) {
	f, err := utils.AFS.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
