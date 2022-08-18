package net

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/everdrone/grab/internal/utils"
)

func Download(url, dest string, options *FetchOptions) error {
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
		return err
	}

	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

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
		return err
	} else {
		if statusOK {
			// Create a empty file
			file, err := utils.Fs.Create(dest)
			if err != nil {
				return err
			}

			defer file.Close()

			// Write the bytes to the file
			_, err = io.Copy(file, res.Body)
			if err != nil {
				return err
			}

			return nil
		} else {
			return fmt.Errorf(res.Status)
		}
	}
}
