package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	MangaChapterPath     = "chapter/%s"
	MangaChaptersPath    = "manga/%s/feed"
	MangaReadMarkersPath = "manga/%s/read"
)

// ChapterService : Provides Chapter services provided by the API.
type ChapterService service

// ChapterList : A response for getting a list of chapters.
type ChapterList struct {
	CommonResponse
	Data []Chapter `json:"data"`
}

func (cl *ChapterList) GetResult() string {
	return cl.Result
}

// Chapter : Struct containing information on a manga.
type Chapter struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Attributes    ChapterAttributes `json:"attributes"`
	Relationships []Relationship    `json:"relationships"`
}

// GetTitle : Get a title for the chapter.
func (c *Chapter) GetTitle() string {
	return c.Attributes.Title
}

// GetChapterNum : Get the chapter's chapter number.
func (c *Chapter) GetChapterNum() string {
	if num := c.Attributes.Chapter; num != nil {
		return *num
	}
	return "-"
}

// All parameters that are accepted when making a Chapter Feed Call
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-manga-id-feed
type ListChapterParams struct {
	Limit                int          `json:"limit" url:"limit,omitempty"`
	Offset               int          `json:"offset" url:"offset,omitempty"`
	Language             []string     `json:"translatedLanguage" url:"translatedLanguage[],omitempty"`
	OGLanguage           []string     `json:"originalLanguage" url:"originalLanguage[],omitempty"`
	EXOGLanguage         []string     `json:"excludedOriginalLanguage" url:"excludedOriginalLanguage[],omitempty"`
	ContentRating        []string     `json:"contentRating" url:"contentRating[],omitempty"`
	ExcludedGroups       []string     `json:"excludedGroups" url:"excludedGroups[],omitempty"`
	ExcludedUploaders    []string     `json:"excludedUploaders" url:"excludedUploaders[],omitempty"`
	IncludeFutureUpdates string       `json:"includeFutureUpdates" url:"includeFutureUpdates,omitempty"`
	CreatedSince         string       `json:"createdAtSince" url:"createdAtSince,omitempty"`
	UpdatedSince         string       `json:"updatedAtSince" url:"updatedAtSince,omitempty"`
	PublishedSince       string       `json:"publishAtSince" url:"publishAtSince,omitempty"`
	Order                ChapterOrder `json:"order" url:"order,omitempty"`
	Includes             []string     `json:"includes" url:"includes[],omitempty"`
	IncludeEmptyPages    int          `json:"includeEmptyPages" url:"includeEmptyPages,omitempty"`
	IncludeFuturePublish int          `json:"includeFuturePublishAt" url:"includeFuturePublishAt,omitempty"`
	IncludeExternalUrl   int          `json:"includeExternalUrl" url:"includeExternalUrl,omitempty"`
}

// Control the ordering of the output of an api call
// All values must either be asc or desc
type ChapterOrder struct {
	Created  string `json:"createdAt" url:"createdAt,omitempty"`
	Updated  string `json:"updatedAt" url:"updatedAt,omitempty"`
	Publish  string `json:"publishAt" url:"publishAt,omitempty"`
	Readable string `json:"readableAt" url:"readableAt,omitempty"`
	Volume   string `json:"volume" url:"volume,omitempty"`
	Chapter  string `json:"chapter" url:"chapter,omitempty"`
}

// ChapterAttributes : Attributes for a Chapter.
type ChapterAttributes struct {
	Title              string  `json:"title"`
	Volume             *string `json:"volume"`
	Chapter            *string `json:"chapter"`
	TranslatedLanguage string  `json:"translatedLanguage"`
	Uploader           string  `json:"uploader"`
	ExternalURL        *string `json:"externalUrl"`
	Version            int     `json:"version"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
	PublishAt          string  `json:"publishAt"`
}

// GetMangaChapters : Get a list of chapters for a manga.
// https://api.mangadex.org/docs.html#operation/get-manga-id-feed
func (s *ChapterService) GetMangaChapters(id string, params *ListChapterParams) (*ChapterList, error) {
	return s.GetMangaChaptersContext(context.Background(), id, params)
}

// GetMangaChaptersContext : GetMangaChapters with custom context.
func (s *ChapterService) GetMangaChaptersContext(ctx context.Context, id string, params *ListChapterParams) (*ChapterList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaChaptersPath, id)

	// Set request parameters
	u.RawQuery = EncodeParams(params)

	var l ChapterList
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

type GetChapterParams struct {
	Includes []string `json:"includes" url:"includes,omitempty"`
}

type SingleChapter struct {
	CommonResponse
	Chapter Chapter `json:"data"`
}

func (c *SingleChapter) GetResult() string {
	return c.Result
}

func (s *ChapterService) GetMangaChapter(id string, params *GetChapterParams) (*SingleChapter, error) {
	return s.GetMangaChapterWithContext(context.Background(), id, params)
}

func (s *ChapterService) GetMangaChapterWithContext(ctx context.Context, id string, params *GetChapterParams) (*SingleChapter, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaChapterPath, id)

	u.RawQuery = EncodeParams(params)

	var l SingleChapter
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// ChapterReadMarkers : A response for getting a list of read chapters.
type ChapterReadMarkers struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

func (rmr *ChapterReadMarkers) GetResult() string {
	return rmr.Result
}

// GetReadMangaChapters : Get list of Chapter IDs that are marked as read for a specified manga ID.
// https://api.mangadex.org/docs.html#operation/get-manga-chapter-readmarkers
func (s *ChapterService) GetReadMangaChapters(id string) (*ChapterReadMarkers, error) {
	return s.GetReadMangaChaptersContext(context.Background(), id)
}

// GetReadMangaChaptersContext : GetReadMangaChapters with custom context.
func (s *ChapterService) GetReadMangaChaptersContext(ctx context.Context, id string) (*ChapterReadMarkers, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaReadMarkersPath, id)

	var rmr ChapterReadMarkers
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &rmr)
	return &rmr, err
}

// SetReadUnreadMangaChapters : Set read/unread manga chapters.
func (s *ChapterService) SetReadUnreadMangaChapters(id string, read, unRead []string) (*Response, error) {
	return s.SetReadUnreadMangaChaptersContext(context.Background(), id, read, unRead)
}

// SetReadUnreadMangaChaptersContext : SetReadUnreadMangaChapters with custom context.
func (s *ChapterService) SetReadUnreadMangaChaptersContext(ctx context.Context, id string, read, unRead []string) (*Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaReadMarkersPath, id)

	// Set request body.
	req := map[string][]string{
		"chapterIdsRead":   read,
		"chapterIdsUnread": unRead,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r Response
	err = s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewBuffer(rBytes), &r)
	return &r, err
}
