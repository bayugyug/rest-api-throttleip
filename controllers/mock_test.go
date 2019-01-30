package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bayugyug/benjerry-icecream/config"
	"github.com/bayugyug/benjerry-icecream/utils"
)

var thandler *ApiHandler
var tAuthToken string

//TestHandler default initializer
func TestHandler(t *testing.T) {
	var err error
	t.Log("Init test")
	var tcfg string
	if os.Getenv("BENJERRY_ICECREAM_CONFIG_DEV") != "" {
		tcfg = os.Getenv("BENJERRY_ICECREAM_CONFIG_DEV")
	} else {
		tcfg = `{"http_port":"8989","driver":{"user":"benjerry_dev","pass":"icecream","port":"3306","name":"benjerry_dev","host":"127.0.0.1"},"showlog":true}`
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
		WithSvcOptDbConf(&appcfg.Config.Driver),
		WithSvcOptDumpFile(appcfg.Config.DumpFile),
	); err != nil {
		t.Fatal("Oops! config might be missing", err)
	}
	t.Log("OK")
}

func clearTable(prodID int64) {
	ApiInstance.DB.Exec("DELETE FROM users     WHERE user='ben@jerry.com' ")
	ApiInstance.DB.Exec("DELETE FROM icecreams WHERE name='test-01-Vanilla Toffee Bar Crunch' ")
	if prodID > 0 {
		ApiInstance.DB.Exec(fmt.Sprintf("DELETE FROM ingredients     WHERE icecream_id='%d' ", prodID))
		ApiInstance.DB.Exec(fmt.Sprintf("DELETE FROM sourcing_values WHERE icecream_id='%d' ", prodID))
	}
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

	clearTable(0)

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
			Method: "POST",
			URL:    "/v1/api/user",
			Body:   `{"user":"ben@jerry.com","pass":"8888"}`,
		},
		{
			Method: "POST",
			URL:    "/v1/api/otp",
			Body:   `{"user":"ben@jerry.com","otp":"{OTP}"}`,
		},
		{
			Method: "POST",
			URL:    "/v1/api/login",
			Body:   `{"user":"ben@jerry.com","pass":"8888"}`,
		},
		{
			Method: "POST",
			URL:    "/v1/api/icecream",
			Body: `{"name": "test-01-Vanilla Toffee Bar Crunch",
						"image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
						"image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
						"description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
						"story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
						"sourcing_values": [
						"Fairtrade",
						"Responsibly Sourced Packaging",
						"Caring Dairy"
						],
						"ingredients": [
						"vegetable oil (canola",
						"safflower",
						"and/or sunflower oil)",
						"guar gum",
						"carrageenan"
						],
						"allergy_info": "may contain wheat, peanuts and other tree nuts",
						"dietary_certifications": "Kosher"}`,
		},
		{
			Method: "PUT",
			URL:    "/v1/api/icecream/{icecream_id}",
			Body: `{"name": "test-01-Vanilla Toffee Bar Crunch",
						"image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
						"image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
						"description": "UPDATED::Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
						"story": "UPDATED::Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
						"sourcing_values": [
							"1-Fairtrade",
							"2-Responsibly Sourced Packaging",
							"3-Caring Dairy"
						],
						"ingredients": [
							"a-vegetable oil (canola",
							"b-safflower",
							"c-and/or sunflower oil)",
							"d-guar gum",
							"e-carrageenan"
						],
						"allergy_info": "UPDATED::may contain wheat, peanuts and other tree nuts",
						"dietary_certifications": "UPDATED::Kosher"}`,
		},
		{
			Method: "GET",
			URL:    "/v1/api/icecream/{icecream_id}",
			Body:   ``,
		},
		{
			Method: "POST",
			URL:    "/v1/api/ingredient/{icecream_id}",
			Body: `{"ingredients": [
					"a1 vegetable oil (canola",
					"b2 safflower",
					"c3 and/or sunflower oil)",
					"d4 guar gum",
					"e5 carrageenan"
					]}`,
		},
		{
			Method: "POST",
			URL:    "/v1/api/sourcing/{icecream_id}",
			Body: `{"sourcing_values": [
					"y1 hehehe Fairtrade",
					"z1 responsibly Sourced Packaging",
					"w1 yez-Caring Dairy"
					]}`,
		},
		{
			Method: "DELETE",
			URL:    "/v1/api/ingredient/{icecream_id}",
			Body:   ``,
		},
		{
			Method: "DELETE",
			URL:    "/v1/api/sourcing/{icecream_id}",
			Body:   ``,
		},
	}

	var otp, icecreamID string

	for _, rec := range mockLists {

		formURL := rec.URL
		formdata := rec.Body

		if rec.URL == "/v1/api/otp" {
			formdata = strings.Replace(formdata, "{OTP}", otp, -1)
			t.Log("OTP::PARAMS::", rec.URL, formdata)
		}

		if strings.Contains(formURL, "{icecream_id}") {
			formURL = strings.Replace(formURL, "{icecream_id}", icecreamID, -1)
		}

		ret, body := testRequest(t, ts, rec.Method, formURL, bytes.NewBufferString(formdata), tAuthToken)
		if ret.StatusCode != http.StatusOK {
			t.Fatalf("Request status:%d", ret.StatusCode)
		}

		t.Log(rec.URL, formURL)

		switch rec.URL {
		case "/v1/api/user":
			var respOtp OtpResponse
			if err := json.Unmarshal([]byte(body), &respOtp); err != nil {
				t.Fatalf("Response failed")
			}
			otp = respOtp.Otp
			t.Log("OTP::", otp)

		case "/v1/api/login":
			var respLog TokenResponse
			if err := json.Unmarshal([]byte(body), &respLog); err != nil {
				t.Fatalf("Response failed")
			}
			tAuthToken = respLog.Token
			utils.Dumper("TOK::", tAuthToken)
		case "/v1/api/icecream":
			var iceRes IcereamResponse
			if err := json.Unmarshal([]byte(body), &iceRes); err != nil {
				t.Fatalf("Response failed")
			}
			icecreamID = iceRes.ProductID
			t.Log("ICECREAM_ID::", icecreamID)
		}
	}

	prodID, _ := strconv.ParseInt(icecreamID, 10, 64)

	clearTable(prodID)
	t.Log("OK")
}
