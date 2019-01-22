package github

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
)

const testResources = "test-resources"

var invalidJSONServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	rw.Write(readTestResource("invalid.json"))
}))

var noResponseServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
}))

func TestAPIService_GetBranches(t *testing.T) {
	// get test data
	jsonServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(readTestResource("get-branches/happy-path.json"))
	}))

	type args struct {
		prefix string
	}

	tests := []struct {
		name    string
		server  *httptest.Server
		service *APIService
		args    args
		want    []string
	}{
		{
			name:    "Test Happy Path with prefix",
			server:  jsonServer,
			service: &APIService{jsonServer.Client(), jsonServer.URL},
			args:    args{prefix: "master"},
			want:    []string{"master"},
		},
		{
			name:    "Test Happy Path without prefix",
			server:  jsonServer,
			service: &APIService{jsonServer.Client(), jsonServer.URL},
			args:    args{prefix: ""},
			want:    []string{"develop", "master"},
		},
		{
			name:    "Test Prefix doesn't match",
			server:  jsonServer,
			service: &APIService{jsonServer.Client(), jsonServer.URL},
			args:    args{prefix: "test"},
			want:    []string{},
		},
		{
			name:    "Test no response path",
			server:  noResponseServer,
			service: &APIService{noResponseServer.Client(), noResponseServer.URL},
			args:    args{prefix: "master"},
			want:    []string{},
		},
		{
			name:    "Test invalid JSON path",
			server:  invalidJSONServer,
			service: &APIService{invalidJSONServer.Client(), invalidJSONServer.URL},
			args:    args{prefix: "master"},
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
		want    int
	}{
		{
			name:    "Test Happy Path",
			server:  jsonServer,
			service: &APIService{jsonServer.Client(), jsonServer.URL},
			want:    1,
		},
		{
			name:    "Test invalid JSOn path",
			server:  invalidJSONServer,
			service: &APIService{invalidJSONServer.Client(), invalidJSONServer.URL},
			want:    0,
		},
		{
			name:    "Test no response path",
			server:  noResponseServer,
			service: &APIService{noResponseServer.Client(), noResponseServer.URL},
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.service.GetAheadBy("test", "test", "develop", "master"); got != tt.want {
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
