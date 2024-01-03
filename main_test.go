package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testType struct {
	name   string
	url    string
	expect any
}

var testCases = []testType{
	{
		name: "no query doesn't have query field.",
		url:  "/",
		expect: JsonResponse{
			Status:  "success",
			Message: "root handler",
		},
	},
	{
		name: "query is like int",
		url:  "/?query=123",
		expect: JsonResponse{
			Status:  "success",
			Message: "root handler",
			Query:   123,
		},
	},
	{
		name: "query is like string.",
		url:  "/?query=foobar",
		expect: JsonResponse{
			Status:  "success",
			Message: "root handler",
			Query:   0,
		},
	},
}

var test404Cases = []testType{
	{
		name:   "not found status code is 404.",
		url:    "/not-found",
		expect: "Not Found.",
	},
}

func TestHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(func() { ts.Close() })

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res := JsonResponse{}
			code, b := testHelper(t, ts, tt.url)
			if err := json.Unmarshal(b, &res); err != nil {
				t.Errorf("error: %s", err)
			}

			if code != 200 {
				t.Errorf("status code is not 200, but %v", code)
			}
			if diff := cmp.Diff(res, tt.expect); diff != "" {
				t.Errorf("diff: %s", diff)
			}
		})
	}
}

func TestNotFoundHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(func() { ts.Close() })

	for _, tt := range test404Cases {
		t.Run(tt.name, func(t *testing.T) {
			var res JsonResponse
			code, b := testHelper(t, ts, tt.url)

			if code != 404 {
				t.Errorf("status code is not 404, but %v", code)
			}
			if diff := cmp.Diff(string(b), tt.expect); diff != "" {
				t.Errorf("p: %v, %v", res, tt.expect)
				t.Errorf("diff: %s", diff)
			}
		})
	}
}

func testHelper(t *testing.T, ts *httptest.Server, u string) (int, []byte) {
	r, err := http.Get(ts.URL + u)
	if err != nil {
		t.Errorf("error: %s", err)
		return 0, nil
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("error: %s", err)
		return 0, nil
	}
	return r.StatusCode, body
}
