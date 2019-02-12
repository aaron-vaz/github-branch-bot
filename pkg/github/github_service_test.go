package github

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testResources = "test-resources"
	githubToken   = "token"
)

var (
	invalidJSONServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(readTestResource("invalid.json"))
	}))

	noResponseServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
)

func TestAPIService_GetBranches(t *testing.T) {
	// get test data
	jsonServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")

		if token != "token "+githubToken {
			t.Error("No github token was supplied")
		}

		rw.Write(readTestResource("get-branches/happy-path.json"))
	}))

	type args struct {
		prefix []string
	}

	tests := []struct {
		name    string
		server  *httptest.Server
		service *APIService
		args    args
		want    []string
	}{
		{
			name:    "Test Happy Path with 1 prefix",
			server:  jsonServer,
			service: &APIService{jsonServer.URL, githubToken, jsonServer.Client()},
			args:    args{prefix: []string{"master"}},
			want:    []string{"master"},
		},
		{
			name:    "Test Happy Path with 2 prefix",
			server:  jsonServer,
			service: &APIService{jsonServer.URL, githubToken, jsonServer.Client()},
			args:    args{prefix: []string{"master", "develop"}},
			want:    []string{"develop", "master"},
		},
		{
			name:    "Test Happy Path without prefix",
			server:  jsonServer,
			service: &APIService{jsonServer.URL, githubToken, jsonServer.Client()},
			args:    args{prefix: []string{""}},
			want:    []string{"develop", "master", "release"},
		},
		{
			name:    "Test Prefix doesn't match",
			server:  jsonServer,
			service: &APIService{jsonServer.URL, githubToken, jsonServer.Client()},
			args:    args{prefix: []string{"test"}},
			want:    []string{},
		},
		{
			name:    "Test no response path",
			server:  noResponseServer,
			service: &APIService{noResponseServer.URL, githubToken, noResponseServer.Client()},
			args:    args{prefix: []string{"master"}},
			want:    []string{},
		},
		{
			name:    "Test invalid JSON path",
			server:  invalidJSONServer,
			service: &APIService{invalidJSONServer.URL, githubToken, invalidJSONServer.Client()},
			args:    args{prefix: []string{"master"}},
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.service.GetBranches("test", "test", tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				if len(got) == 0 && len(tt.want) == 0 {
					return
				}

				t.Errorf("APIService.GetBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIService_GetAheadBy(t *testing.T) {
	jsonServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(readTestResource("get-ahead-by/happy-path.json"))
	}))

	tests := []struct {
		name    string
		server  *httptest.Server
		service *APIService
		want    map[string]int
	}{
		{
			name:    "Test Happy Path",
			server:  jsonServer,
			service: &APIService{jsonServer.URL, githubToken, jsonServer.Client()},
			want:    map[string]int{"master": 1},
		},
		{
			name:    "Test invalid JSOn path",
			server:  invalidJSONServer,
			service: &APIService{invalidJSONServer.URL, githubToken, invalidJSONServer.Client()},
			want:    map[string]int{"master": 0},
		},
		{
			name:    "Test no response path",
			server:  noResponseServer,
			service: &APIService{noResponseServer.URL, githubToken, noResponseServer.Client()},
			want:    map[string]int{"master": 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.service.GetAheadBy("test", "test", "develop", []string{"master"}); !cmp.Equal(got, tt.want) {
				t.Errorf("APIService.GetAheadBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func readTestResource(path string) []byte {
	content, err := ioutil.ReadFile(filepath.Join(testResources, path))
	if err != nil {
		log.Fatalln(err)
	}

	return content
}
