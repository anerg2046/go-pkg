package httpclient

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var HttpClient = resty.New()

func init() {
	// HttpClient.SetDebug(true)
	HttpClient.SetRetryCount(3)
	HttpClient.SetRetryWaitTime(5 * time.Second)
	HttpClient.SetRetryMaxWaitTime(20 * time.Second)
	HttpClient.SetTimeout(15 * time.Second)
	HttpClient.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			if r.StatusCode() != http.StatusOK || r == nil {
				return true
			}
			return false
		},
	)
}
