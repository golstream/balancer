package healthcheck

import (
	"balancer/internal/multiplexer"
	httputils "balancer/pkg/httputil"
	"balancer/pkg/utils"
	"context"
	"time"
)

type HealthCheck struct {
	initUrls   []string
	interSec   time.Duration
	timeoutSec time.Duration
}

func New(urls []string, intervalSec int, timeoutSec int) HealthCheck {
	return HealthCheck{
		initUrls:   urls,
		interSec:   time.Duration(intervalSec) * time.Second,
		timeoutSec: time.Duration(timeoutSec) * time.Second,
	}
}

func (hc HealthCheck) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.interSec)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			urls := multiplexer.GetSliceOfURLs()
			hc.checkURLs(ctx, utils.Ternary(urls == nil, hc.initUrls, urls))
		}
	}
}

func (hc HealthCheck) checkURLs(ctx context.Context, urls []string) {
	var (
		availableURLs = make([]string, 0, len(urls))
	)

	for _, url := range urls {
		resp, err := httputils.GetWithCtx(
			ctx,
			url,
			nil,
			nil,
			nil,
			hc.timeoutSec)
		if err != nil {
			continue
		}

		if !resp.GetStatusCode().Is2xxStatusCode() {
			continue
		}

		availableURLs = append(availableURLs, url)
	}

	multiplexer.Urls.Store(availableURLs)
}
