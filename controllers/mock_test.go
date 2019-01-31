package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bayugyug/rest-api-throttleip/config"
	"github.com/bayugyug/rest-api-throttleip/utils"
)

var thandler *ApiHandler
var tAuthToken string

//TestHandler default initializer
func TestHandler(t *testing.T) {
	var err error
	t.Log("Init test")
	var tcfg string
	if os.Getenv("REST_API_THROTTLEIP_DEV") != "" {
		tcfg = os.Getenv("REST_API_THROTTLEIP_DEV")
	} else {
		tcfg = `{"http_port":"8989","redis_host":"127.0.0.1:6379","showlog":true}`
	}
	//init
	thandler = &ApiHandler{}

	//init
	appcfg := config.NewAppSettings(config.WithSetupCmdParams(tcfg))

	//check
	if appcfg.Config == nil {
		t.Fatal("Oops! Config missing")
	}

	//init service
	if ApiInstance, err = NewApiService(
		WithSvcOptAddress(":"+appcfg.Config.HttpPort),
		WithSvcOptRedisHost(appcfg.Config.RedisHost),
	); err != nil {
		t.Fatal("Oops! config might be missing", err)
	}
	t.Log("OK")
}

//testRequest test for http req
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, auth string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	if auth != "" {
		req.Header.Add("Authorization", "Bearer "+auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

//TestHandlers
func TestHandlers(t *testing.T) {

	//setup
	ts := httptest.NewServer(ApiInstance.Router)
	defer ts.Close()

	mockLists := []struct {
		Method string
		URL    string
		Ctx    context.Context
		Body   string
	}{
		{
			Method: "GET",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "POST",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "DELETE",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "GET",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "POST",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "DELETE",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "GET",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "POST",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/request/{dummy}",
			Body:   ``,
		},
	}

	var how int
	for _, rec := range mockLists {

		how++
		formURL := strings.Replace(rec.URL, "{dummy}", utils.UHelper.UUID(), -1)

		ret, body := testRequest(t, ts, rec.Method, formURL, bytes.NewBufferString(rec.Body), "")
		if ret.StatusCode != http.StatusOK {
			t.Fatalf("Request status:%d", ret.StatusCode)
		}
		var reply APIResponse
		if err := json.Unmarshal([]byte(body), &reply); err != nil {
			t.Fatalf("Response failed")
		}
		if reply.Code <= 0 || reply.Status == "" {
			t.Fatalf("Response failed")
		}
		t.Log(how, "OKAY", reply.Code, reply.Status)
	}

	t.Log("OK")
}
