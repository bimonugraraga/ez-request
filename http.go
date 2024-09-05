package ezrequest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type RequestParams struct {
	Ctx                  context.Context
	Method               string
	URL                  string
	BodyPayload          string
	Attempts             int
	BackoffMs            int64
	TimeoutMs            int64
	Headers              map[string]string
	StatusCodeConstraint []int
}

func (r *RequestParams) createRequest() (*http.Request, error) {
	// Create Request
	req, err := http.NewRequestWithContext(r.Ctx, r.Method, r.URL, bytes.NewBuffer([]byte(r.BodyPayload)))
	if err != nil {
		return nil, err
	}

	for key, val := range r.Headers {
		req.Header.Set(key, val)
	}

	return req, nil
}

func (r *RequestParams) EzRequest() (*http.Response, error) {
	// Create Request
	req, err := r.createRequest()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Send  Request
	httpClient := &http.Client{Timeout: time.Duration(r.TimeoutMs) * time.Millisecond}
	res, err := r.EzDoIt(req, httpClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *RequestParams) EzRetriableRequest() (*http.Response, error) {
	// Create Request
	req, err := r.createRequest()
	if err != nil {
		return nil, err
	}

	// Send  Request
	httpClient := &http.Client{Timeout: time.Duration(r.TimeoutMs) * time.Millisecond}
	res, err := r.EzDoIt(req, httpClient)
	if err == nil {
		return res, nil
	}

	if r.BackoffMs == 0 {
		r.BackoffMs = 10
	}

	// Retry when error
	for i := 0; i < r.Attempts; i++ {
		time.Sleep(time.Millisecond * time.Duration(r.BackoffMs))
		res, err = r.EzDoIt(req, httpClient)
		if err == nil {
			return res, nil
		}
	}

	return nil, err
}

func (r *RequestParams) EzDoIt(req *http.Request, httpClient *http.Client) (*http.Response, error) {
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if len(r.StatusCodeConstraint) != 0 && !r.checkStatusCode(res.StatusCode) {
		return nil, errors.New(fmt.Sprintf("error constraint with status code: %d", res.StatusCode))
	}
	return res, nil
}

func (r *RequestParams) EzDoItRetriable(req *http.Request, httpClient *http.Client) (*http.Response, error) {
	res, err := httpClient.Do(req)
	if (err == nil && len(r.StatusCodeConstraint) == 0) || (err == nil && r.checkStatusCode(res.StatusCode)) {
		return res, nil
	}

	for i := 0; i < r.Attempts; i++ {
		time.Sleep(time.Millisecond * time.Duration(r.BackoffMs))
		res, err := httpClient.Do(req)
		if (err == nil && len(r.StatusCodeConstraint) == 0) || (err == nil && r.checkStatusCode(res.StatusCode)) {
			return res, nil
		}
	}

	if err != nil {
		return nil, err
	}

	if len(r.StatusCodeConstraint) != 0 && !r.checkStatusCode(res.StatusCode) {
		return nil, errors.New(fmt.Sprintf("error constraint with status code: %d", res.StatusCode))
	}

	return res, nil
}

func (r *RequestParams) checkStatusCode(targetStatus int) bool {
	for _, val := range r.StatusCodeConstraint {
		if val == targetStatus {
			return true
		}
	}
	return false
}
