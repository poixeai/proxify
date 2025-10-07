package controller

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/infra/config"
	"github.com/poixeai/proxify/infra/ctx"
	"github.com/poixeai/proxify/infra/logger"
	"github.com/poixeai/proxify/infra/response"
)

func ProxyHandler(c *gin.Context) {
	// if reserved route, return 404
	topRoute := c.GetString(ctx.TopRoute)
	if config.ReservedTopRoutes[topRoute] {
		logger.Warnf("404 Not Found: %s", topRoute)
		response.RespondSystemRouteNotFoundError(c)
		return
	}

	// build target URL
	targetEndpoint := c.GetString(ctx.TargetEndpoint)
	subPath := c.GetString(ctx.SubPath)
	targetURL := targetEndpoint + subPath

	// construct new request
	ctx := c.Request.Context()
	req, err := http.NewRequestWithContext(ctx, c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		logger.Errorf("Failed to create new request: %v", err)
		response.RespondInternalError(c)
		return
	}

	// copy headers
	for k, v := range c.Request.Header {
		req.Header[k] = v
	}

	// create client
	client := &http.Client{
		Timeout: 0, // no timeout, let ctx control it
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DisableCompression:  true, // disable gzip, avoid stream cache
			MaxIdleConnsPerHost: 50,
		},
	}

	// do request
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Failed to do request to target: %v", err)
		response.RespondInternalError(c)
		return
	}
	defer resp.Body.Close()

	// copy response headers
	for k, v := range resp.Header {
		c.Writer.Header()[k] = v
	}

	// set status code
	c.Status(resp.StatusCode)

	// determine if response is a stream
	ct := resp.Header.Get("Content-Type")
	te := resp.Header.Get("Transfer-Encoding")
	isStream := strings.Contains(te, "chunked") || strings.Contains(ct, "text/event-stream")

	if isStream {
		streamCopy(c, resp)
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}

// stream support SSE / chunked
func streamCopy(c *gin.Context, resp *http.Response) {
	buf := make([]byte, 4096)
	writer := c.Writer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := writer.Write(buf[:n])
			if writeErr != nil {
				return // client disconnected
			}
			writer.Flush() // keep flushing to client
		}
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
	}
}
