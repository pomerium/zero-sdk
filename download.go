package zerosdk

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pomerium/zero-sdk/apierror"
	cluster_api "github.com/pomerium/zero-sdk/cluster"
)

const (
	maxErrorResponseBodySize = 2 << 14 // 32kb
)

// DownloadClusterResourceBundle downloads given cluster resource bundle to given writer.
func (api *API) DownloadClusterResourceBundle(
	ctx context.Context,
	dst io.Writer,
	id string,
	current *DownloadConditional,
) (*DownloadResult, error) {
	req, err := api.getDownloadRequest(ctx, id, current)
	if err != nil {
		return nil, fmt.Errorf("get download request: %w", err)
	}

	resp, err := api.cfg.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return &DownloadResult{NotModified: true}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpDownloadError(ctx, resp)
	}

	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("write body: %w", err)
	}

	updated, err := newConditionalFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain cache conditions from response: %w", err)
	}

	return &DownloadResult{
		DownloadConditional: updated,
	}, nil
}

func (api *API) getDownloadRequest(ctx context.Context, id string, current *DownloadConditional) (*http.Request, error) {
	url, err := api.getDownloadURL(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get download URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	err = current.SetHeaders(req)
	if err != nil {
		return nil, fmt.Errorf("set conditional download headers: %w", err)
	}

	return req, nil
}

func (api *API) getDownloadURL(ctx context.Context, id string) (*url.URL, error) {
	u, ok := api.downloadURLCache.Get(id, api.cfg.downloadURLCacheTTL)
	if ok {
		return &u, nil
	}

	return api.updateBundleDownloadURL(ctx, id)
}

func (api *API) updateBundleDownloadURL(ctx context.Context, id string) (*url.URL, error) {
	now := time.Now()

	resp, err := apierror.CheckResponse[cluster_api.DownloadBundleResponse](
		api.cluster.DownloadClusterResourceBundleWithResponse(ctx, id),
	)
	if err != nil {
		return nil, fmt.Errorf("get bundle download URL: %w", err)
	}

	expiresSeconds, err := strconv.ParseInt(resp.ExpiresInSeconds, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse expiration: %w", err)
	}

	u, err := url.Parse(resp.Url)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	api.downloadURLCache.Set(id, *u, now.Add(time.Duration(expiresSeconds)*time.Second))
	return u, nil
}

// DownloadResult contains the result of a download operation
type DownloadResult struct {
	// NotModified is true if the bundle has not been modified
	NotModified bool
	// DownloadConditional contains the new conditional
	*DownloadConditional
}

type DownloadConditional struct {
	ETag         string
	LastModified string
}

func (c *DownloadConditional) Validate() error {
	if c.ETag == "" && c.LastModified == "" {
		return fmt.Errorf("either ETag or LastModified must be set")
	}
	return nil
}

func (c *DownloadConditional) SetHeaders(req *http.Request) error {
	if c == nil {
		return nil
	}
	if err := c.Validate(); err != nil {
		return err
	}
	req.Header.Set("If-None-Match", c.ETag)
	req.Header.Set("If-Modified-Since", c.LastModified)
	return nil
}

func newConditionalFromResponse(resp *http.Response) (*DownloadConditional, error) {
	c := &DownloadConditional{
		ETag:         resp.Header.Get("ETag"),
		LastModified: resp.Header.Get("Last-Modified"),
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

type xmlError struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
	Details string   `xml:"Details"`
}

func (e xmlError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func tryXMLError(body []byte) (bool, error) {
	var xmlErr xmlError
	err := xml.Unmarshal(body, &xmlErr)
	if err != nil {
		return false, fmt.Errorf("unmarshal xml error: %w", err)
	}

	return true, xmlErr
}

func httpDownloadError(ctx context.Context, resp *http.Response) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, io.LimitReader(resp.Body, maxErrorResponseBodySize))

	if isXML(resp.Header.Get("Content-Type")) {
		ok, err := tryXMLError(buf.Bytes())
		if ok {
			return err
		}
	}

	log.Ctx(ctx).Debug().Err(err).
		Str("error", resp.Status).
		Str("body", buf.String()).Msg("bundle download error")

	return fmt.Errorf("download error: %s", resp.Status)
}

// isXML parses content-type for application/xml
func isXML(ct string) bool {
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return false
	}
	return mediaType == "application/xml"
}
