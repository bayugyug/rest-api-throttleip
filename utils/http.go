package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type HttpCurl struct {
	HttpClient *http.Client
	Timeout    time.Duration
}

func NewHttpCurl() *HttpCurl {
	h := &HttpCurl{
		Timeout: time.Duration(30),
	}
	h.Init()
	return h

}

//Init start prep
func (g *HttpCurl) Init() {
	//web
	g.HttpClient = &http.Client{
		Timeout: time.Duration(g.Timeout * time.Second),
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   g.Timeout * time.Second,
				KeepAlive: 0,
			}).Dial,
			TLSHandshakeTimeout: g.Timeout * time.Second,
		},
	}
}

//Get request via get
func (g *HttpCurl) Get(url string) (string, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", -1, err
	}
	//settings
	req.Close = true
	req.Header.Set("Connection", "close")
	resp, err := g.HttpClient.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return "", -1, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	// read the body
	return strings.TrimSpace(string(contents)), resp.StatusCode, nil
}
