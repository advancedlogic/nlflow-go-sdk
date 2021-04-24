package cache

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*
https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-get-save-state/#step-2-save-and-retrieve-a-single-state
dapr run --app-id myapp --dapr-http-port 3500
*/
type nlflowCacheDapr struct {
	baseUrl string
	client  *http.Client
}

const statestore_path = "/v1.0/state/statestore"

func NewNLFlowCacheDapr() NLFlowCache {
	return &nlflowCacheDapr{
		baseUrl: "http://localhost:3500",
		client:  &http.Client{},
	}
}

func (c *nlflowCacheDapr) Close() error {
	return nil
}

func (c *nlflowCacheDapr) Read(k string) (string, error) {
	resp, err := c.httpRequest(c.baseUrl+statestore_path+"/"+k, "GET", map[string]string{}, "")

	// return "", KeyNotFound

	if err != nil {
		log.Println("Http GET Error")
		return "", err
	}
	return resp, nil
}

func (c *nlflowCacheDapr) Write(k string, v string) error {
	headers := map[string]string{"Content-Type": "application/json"}
	body := fmt.Sprintf("[ { \"key\": \"%s\", \"value\": %s } ]", k, v)
	_, err := c.httpRequest(c.baseUrl+statestore_path, http.MethodPost, headers, body)
	if err != nil {
		log.Println("Http POST Error")
		return err
	}
	return err
}

func (c *nlflowCacheDapr) httpRequest(url string, requestType string, headers map[string]string, body string) (string, error) {
	var req *http.Request
	var err error

	if requestType == "GET" {
		req, err = http.NewRequest(requestType, url, nil)
	} else {
		var bodyStr = []byte(body)
		req, err = http.NewRequest(requestType, url, bytes.NewBuffer(bodyStr))
	}

	if err != nil {
		return "", err
	}

	// Custom Headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	// If successful HTTP call, but Client/Server error, we return error
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return "", fmt.Errorf("%d Http Client Error for url: %s", resp.StatusCode, url)
	}
	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return "", fmt.Errorf("%d Http Server Error for url: %s", resp.StatusCode, url)
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	responseString := string(response)
	if err != nil {
		log.Println("ERROR reading response for URL: ", url)
		return "", err
	}

	return responseString, nil
}
