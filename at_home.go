package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	GetMDHomeURLPath = "at-home/server/%s"
	MDHomeReportURL  = "https://api.mangadex.network/report"
	dataSaver = "data-saver"
	force443 = "forcePort443"
)

// AtHomeService : Provides MangaDex@Home services provided by the API.
type AtHomeService service

// MDHomeServerResponse : A response for getting a server URL to get chapters.
type MDHomeServerResponse struct {
	Result  string       `json:"result"`
	BaseURL string       `json:"baseUrl"`
	Chapter ChaptersData `json:"chapter"`
}

func (r *MDHomeServerResponse) GetResult() string {
	return r.Result
}

// ChaptersData : Struct containing data for the chapter's pages.
type ChaptersData struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

// MDHomeClient : Client for interfacing with MangaDex@Home.
type MDHomeClient struct {
	client  *http.Client
	baseURL string
	quality string
	hash    string
	report  *reportPayload
	start    time.Time
	Pages   []string
}

// NewMDHomeClient : Get MangaDex@Home client for a chapter.
// https://api.mangadex.org/docs.html#operation/get-at-home-server-chapterId
func (s *AtHomeService) NewMDHomeClient(chapterID string, quality string, forcePort443 bool) (*MDHomeClient, error) {
	return s.NewMDHomeClientContext(context.Background(), chapterID, quality, forcePort443)
}

// NewMDHomeClientContext : NewMDHomeClient with custom context.
func (s *AtHomeService) NewMDHomeClientContext(ctx context.Context, chapterID string, quality string, forcePort443 bool) (*MDHomeClient, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetMDHomeURLPath, chapterID)

	// Set query parameters
	q := u.Query()
	q.Set(force443, strconv.FormatBool(forcePort443))
	u.RawQuery = q.Encode()

	var r MDHomeServerResponse
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r)
	if err != nil {
		return nil, err
	}

	// Set the required pages data for required quality.
	pages := r.Chapter.Data
	if quality == dataSaver {
		pages = r.Chapter.DataSaver
	}

	return &MDHomeClient{
		client:  &http.Client{},
		baseURL: r.BaseURL,
		quality: quality,
		hash:    r.Chapter.Hash,
		Pages:   pages,
	}, nil
}

// GetChapterPage : Return page data for a chapter with the filename of that page.
func (c *MDHomeClient) GetChapterPage(filename string) ([]byte, error) {
	return c.GetChapterPageWithContext(context.Background(), filename)
}

// GetChapterPageWithContext : GetChapterPage with custom context.
func (c *MDHomeClient) GetChapterPageWithContext(ctx context.Context, filename string) (fileData []byte, err error) {
	// Send report in the background.
	defer func() { go c.reportContext(ctx) }()

	path := strings.Join([]string{c.baseURL, c.quality, c.hash, filename}, "/")
	c.report = newPayload(path)

	// time will be slightly inflated, but w/e
	c.start = time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Read file data.
	fileData, err = io.ReadAll(resp.Body)
	c.report.Bytes = len(fileData)
	c.report.Success = err == nil && resp.StatusCode == http.StatusOK 
	c.report.Cached = strings.HasPrefix(resp.Header.Get("X-Cache"), "HIT")
	return
}

func (c *MDHomeClient) FileUrl(filename string) string {
	return strings.Join([]string{c.baseURL, c.quality, c.hash, filename}, "/")
}

// reportPayload : Required fields for reporting page download result.
type reportPayload struct {
	URL      string
	Success  bool
	Bytes    int
	Duration int64
	Cached   bool
}

func newPayload(path string) *reportPayload {
	return &reportPayload{
		URL:     path,
	}
}

// reportContext : Report success of getting chapter page data.
func (c *MDHomeClient) reportContext(ctx context.Context) {
	c.report.Duration = time.Since(c.start).Milliseconds()
	rBytes, err := json.Marshal(c.report)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MDHomeReportURL, bytes.NewBuffer(rBytes))
	if err != nil {
		return
	}
	c.client.Do(req)
}
