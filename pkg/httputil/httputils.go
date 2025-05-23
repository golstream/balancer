package httputils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetWithCtx(
	ctx context.Context,
	url string,
	queries map[string][]string,
	headers map[string]string,
	cookies []*http.Cookie,
	timeout time.Duration,
) (*Response, error) {
	return DoRequestWithCtx(ctx, http.MethodGet, url, queries, nil, headers, cookies, timeout)
}

func PostWithCtx(
	ctx context.Context,
	url string,
	queries map[string][]string,
	body []byte,
	headers map[string]string,
	cookies []*http.Cookie,
	timeout time.Duration,
) (*Response, error) {
	return DoRequestWithCtx(ctx, http.MethodPost, url, queries, body, headers, cookies, timeout)
}

func PutWithCtx(
	ctx context.Context,
	url string,
	queries map[string][]string,
	body []byte,
	headers map[string]string,
	cookies []*http.Cookie,
	timeout time.Duration,
) (*Response, error) {
	return DoRequestWithCtx(ctx, http.MethodPut, url, queries, body, headers, cookies, timeout)
}

func DeleteWithCtx(
	ctx context.Context,
	url string,
	queries map[string][]string,
	body []byte,
	headers map[string]string,
	cookies []*http.Cookie,
	timeout time.Duration,
) (*Response, error) {
	return DoRequestWithCtx(ctx, http.MethodDelete, url, queries, body, headers, cookies, timeout)
}

func ReadBody(reader io.ReadCloser) ([]byte, error) {
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func DoRequestWithCtx(
	ctx context.Context,
	method string,
	url string,
	queries map[string][]string,
	body []byte,
	headers map[string]string,
	cookies []*http.Cookie,
	timeout time.Duration,
) (*Response, error) {
	var (
		err     error
		request *http.Request
	)

	request, err = http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("%s %w",
			statusFailedHttpResponse, err)
	}

	request.Header.Set(
		ContentTypeHeader,
		MIMEApplicationJSON,
	)

	for key, val := range headers {
		if key == ContentTypeHeader {
			continue
		}
		request.Header.Set(key, val)
	}

	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	query := request.URL.Query()
	for key, val := range queries {
		for _, v := range val {
			query.Add(key, v)
		}
	}
	request.URL.RawQuery = query.Encode()

	var (
		response *http.Response
	)

	client := http.Client{
		Timeout: timeout,
	}

	if response, err = client.Do(request); err != nil {
		return nil, fmt.Errorf("%s %w, %s, %s",
			statusFailedHttpRequest, err, statusTimeout, timeout.String())
	}

	if response == nil {
		return nil, ErrEmptyResponse
	}

	var (
		responseBody []byte
	)

	if responseBody, err = io.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("%s: %w",
			statusFailedHttpResponse, err)
	}

	resp := &Response{
		code: statusCode(response.StatusCode),
		body: responseBody,
	}

	return resp, nil
}
