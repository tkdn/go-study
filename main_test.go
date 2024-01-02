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
	url string
	expect any
}

var testCases = []testType{
	{
		url: "/",
		expect: JsonResponse{
			Status: "success",
			Message: "root handler",	
		},
	},
	{
		url: "/?query=123",
		expect: JsonResponse{
			Status: "success",
			Message: "root handler",
			Query: 123,
		},
	},
	{
		url: "/?query=foobar",
		expect: JsonResponse{
			Status: "success",
			Message: "root handler",
			Query: 0,
		},
	},
}

var test404Cases = []testType{
	{
		url: "/not-found",
		expect: "Not Found.",
	},
}

func TestHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(func() { ts.Close() })

	for _, tc := range testCases {
		var res JsonResponse
		code, b := testHelper(t, ts, tc.url)
		if err := json.Unmarshal(b, &res); err != nil {
			t.Errorf("error: %s", err)
		}

		if code != 200 {
			t.Errorf("status code is not 200, but %v", code)
		}
		if diff := cmp.Diff(res, tc.expect); diff != "" {
			t.Errorf("diff: %s", diff)
		}
	}
}

func TestNotFoundHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(func() { ts.Close() })

	for _, tc := range test404Cases {
		var res JsonResponse
		code, b := testHelper(t, ts, tc.url)

		if code != 404 {
			t.Errorf("status code is not 404, but %v", code)
		}
		if diff := cmp.Diff(string(b), tc.expect); diff != "" {
			t.Errorf("p: %v, %v", res, tc.expect)
			t.Errorf("diff: %s", diff)
		}
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
