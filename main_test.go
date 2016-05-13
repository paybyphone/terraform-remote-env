package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const remoteStateResponseText = `
{
	"version": 1,
	"serial": 1,
	"remote": {
		"type": "http",
		"config": {
			"address": "http://127.0.0.1:12345/",
			"skip_cert_verification": "0"
		}
	},
	"modules": [{
		"path": [
			"root"
		],
		"outputs": {
			"foo": "bar",
			"baz": "qux"
		},
		"resources": {}
	}]
}
`

func remoteStateOutputExpected() map[string]string {
	return map[string]string{
		"foo": "bar",
		"baz": "qux",
	}
}

const varsOutputExpectedPlain = "TF_VAR_foo=bar TF_VAR_baz=qux"
const varsOutputExpectedPrefixed = "TF_VAR_foobar_foo=bar TF_VAR_foobar_baz=qux"

func newHTTPTestServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(f))
	return ts
}

func httpRemoteStateTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, remoteStateResponseText, http.StatusOK)
	})
}

func testConfig() programConfig {
	return programConfig{
		backend:       "http",
		backendConfig: map[string]string{},
	}
}

// TestParseArgs tests the parseArgs function.
func TestParseArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"cmd",
		"-prefix=foobar",
		"-backend=s3",
		"-backend-config=bucket=test-bucket",
		"-backend-config=key=foo",
		"-backend-config=region=bar",
	}

	cfg := parseArgs()
	if cfg.backend != "s3" {
		t.Fatalf("Expected backend to be s3, got %v", cfg.backend)
	}
	if cfg.prefix != "foobar" {
		t.Fatalf("Expected prefix to be foobar, got %v", cfg.prefix)
	}

	backendConfigExpected := map[string]string{
		"bucket": "test-bucket",
		"key":    "foo",
		"region": "bar",
	}

	if reflect.DeepEqual(cfg.backendConfig, backendConfigExpected) != true {
		t.Fatalf("Expected backend to be %v, got %v", backendConfigExpected, cfg.backendConfig)
	}
}

// testGetMockState is a helper function to TestGetState and TestOutputState_*
// that pulls the test state from the mock HTTP server.
func testGetMockState(t *testing.T) map[string]string {
	ts := httpRemoteStateTestServer()
	defer ts.Close()
	cfg := testConfig()
	cfg.backendConfig["address"] = ts.URL

	out, err := getState(cfg)

	if err != nil {
		t.Fatalf("Unexpected request error: %s", err)
	}

	return out
}

// TestGetState tests the GetState function.
func TestGetState(t *testing.T) {
	out := testGetMockState(t)

	expected := remoteStateOutputExpected()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

// TestOutputState_plain tests the outputState function, without tags.
func TestOutputState_plain(t *testing.T) {
	in := testGetMockState(t)
	out := strings.Split(" ", outputState(testConfig(), in))

	expected := strings.Split(" ", varsOutputExpectedPlain)

	sort.Strings(out)
	sort.Strings(expected)

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

// TestOutputState_prefixed tests the outputState function, with a prefix.
func TestOutputState_prefixed(t *testing.T) {
	in := testGetMockState(t)
	cfg := testConfig()
	cfg.prefix = "foobar"
	out := strings.Split(" ", outputState(cfg, in))

	expected := strings.Split(" ", varsOutputExpectedPrefixed)

	sort.Strings(out)
	sort.Strings(expected)

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}
