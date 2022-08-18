package net

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Fetch(url string, options *FetchOptions) (string, error) {
	retriesLeft := options.Retries

	if options.Retries < 1 {
		retriesLeft = 1
	}

	if options.Timeout < 1 {
		options.Timeout = 10000
	}

	client := &http.Client{
		Timeout: time.Duration(options.Timeout) * time.Millisecond,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	var body []byte
	var res *http.Response

	statusOK := false
	for retriesLeft > 0 {
		res, err = client.Do(req)
		// we get an error or no response, so retry
		if err != nil || res == nil {
			retriesLeft -= 1
			continue
		}

		defer res.Body.Close()

		// we get an "ok" response, so we can stop retrying
		statusOK = res.StatusCode >= 200 && res.StatusCode < 300
		if statusOK {
			break
		}

		retriesLeft -= 1
	}

	if res == nil {
		return "", err
	} else {
		if statusOK {
			body, err = ioutil.ReadAll(res.Body)
			return string(body), err
		} else {
			return "", fmt.Errorf(res.Status)
		}
	}
}
