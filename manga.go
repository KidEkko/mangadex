package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	MangaPath                = "manga/%s"
	MangaAggregatePath       = "manga/%s/aggregate"
	MangaListPath            = "manga"
	CheckIfMangaFollowedPath = "user/follows/manga/%s"
	ToggleMangaFollowPath    = "manga/%s/follow"
)

// MangaService : Provides Manga services provided by the API.
type MangaService service

// MangaList : A response for getting a list of manga.
type MangaList struct {
	CommonResponse
	Data []Manga `json:"data"`
}

func (ml *MangaList) GetResult() string {
	return ml.Result
}

// Manga : Struct containing information on a Manga.
type Manga struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    MangaAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
}

// GetTitle : Get title of the Manga.
func (m *Manga) GetTitle(langCode string) string {
	if title := m.Attributes.Title.GetLocalString(langCode); title != "" {
		return title
	}
	return m.Attributes.AltTitles.GetLocalString(langCode)
}

// GetDescription : Get description of the Manga.
func (m *Manga) GetDescription(langCode string) string {
	return m.Attributes.Description.GetLocalString(langCode)
}

type ListMangaParams struct {
	Limit                int        `json:"limit" url:"limit,omitempty"`
	Offset               int        `json:"offset" url:"offset,omitempty"`
	Title                string     `json:"title" url:"title,omitempty"`
	AuthorArtist         string     `json:"authorOrArtist" url:"authorOrArtist,omitempty"`
	Authors              []string   `json:"authors" url:"authors[],omitempty"`
	Artists              []string   `json:"artists" url:"artists[],omitempty"`
	Year                 string     `json:"year" url:"year,omitempty"` // can also be an int if you want to make your own struct, don't care enough to make both
	Tags                 []string   `json:"includedTags" url:"includedTags[],omitempty"`
	TagMode              string     `json:"includedTagsMode" url:"includedTagsMode,omitempty"` // default "AND"
	ExcludedTags         []string   `json:"excludedTags" url:"excludedTags[],omitempty"`
	ExcludedTagsMode     string     `json:"excludedTagsMode" url:"excludedTagsMode,omitempty"` // default "AND"
	Status               []string   `json:"status" url:"status[],omitempty"`                   // "ongoing" "completed" "hiatus" "cancelled"
	OGLanguage           []string   `json:"originalLanguage" url:"originalLanguage[],omitempty"`
	EXOGLanguage         []string   `json:"excludedOriginalLanguage" url:"excludedOriginalLanguage[],omitempty"`
	AvailableLanguages   []string   `json:"availableTranslatedLanguage" url:"availableTranslatedLanguage[],omitempty"` // filter by available languages
	Demographic          []string   `json:"publicationDemographic" url:"publicationDemographic[],omitempty"`           // "shounen" "shoujo" "josei" "seinen" "none"
	Ids                  []string   `json:"ids" url:"ids[],omitempty"`
	ContentRating        []string   `json:"contentRating" url:"contentRating[],omitempty"`
	CreatedSince         string     `json:"createdAtSince" url:"createdAtSince,omitempty"`
	UpdatedSince         string     `json:"updatedAtSince" url:"updatedAtSince,omitempty"`
	Order                MangaOrder `json:"order" url:"order,omitempty"`
	Includes             []string   `json:"includes" url:"includes[],omitempty"`                       // "manga" "cover_art" "author" "artist" "tag" "creator"
	HasAvailableChapters string     `json:"hasAvailableChapters" url:"hasAvailableChapters,omitempty"` // "0" "1" "true" "false"
	Group                string     `json:"group" url:"group,omitempty"`
}

// Control the ordering of the output of an api call
// All values must either be asc or desc
type MangaOrder struct {
	Title        string `json:"title" url:"title,omitempty"`
	Year         string `json:"year" url:"year,omitempty"`
	Created      string `json:"createdAt" url:"createdAt,omitempty"`
	Updated      string `json:"updatedAt" url:"updatedAt,omitempty"`
	LatestUpload string `json:"latestUploadedChapter" url:"latestUploadedChapter,omitempty"`
	FollowCount  string `json:"followedCount" url:"followedCount,omitempty"`
	Relevance    string `json:"relevance" url:"relevance,omitempty"`
	Rating       string `json:"rating" url:"rating,omitempty"`
}

// MangaAttributes : Attributes for a Manga.
type MangaAttributes struct {
	Title                  LocalisedStrings `json:"title"`
	AltTitles              LocalisedStrings `json:"altTitles"`
	Description            LocalisedStrings `json:"description"`
	IsLocked               bool             `json:"isLocked"`
	Links                  LocalisedStrings `json:"links"`
	OriginalLanguage       string           `json:"originalLanguage"`
	LastVolume             *string          `json:"lastVolume"`
	LastChapter            *string          `json:"lastChapter"`
	PublicationDemographic *string          `json:"publicationDemographic"`
	Status                 *string          `json:"status"`
	Year                   *int             `json:"year"`
	ContentRating          *string          `json:"contentRating"`
	Tags                   []Tag            `json:"tags"`
	State                  string           `json:"state"`
	Version                int              `json:"version"`
	CreatedAt              string           `json:"createdAt"`
	UpdatedAt              string           `json:"updatedAt"`
}

// GetMangaList : Get a list of Manga.
// https://api.mangadex.org/docs.html#operation/get-search-manga
func (s *MangaService) GetMangaList(params *ListMangaParams) (*MangaList, error) {
	return s.GetMangaListContext(context.Background(), params)
}

// GetMangaListContext : GetMangaList with custom context.
func (s *MangaService) GetMangaListContext(ctx context.Context, params *ListMangaParams) (*MangaList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = MangaListPath

	// Set query parameters
	u.RawQuery = EncodeParams(params)

	var l MangaList
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

type GetMangaParams struct {
	Includes []string `json:"includes" url:"includes,omitempty"`
}

type SingleManga struct {
	CommonResponse
	Manga Manga `json:"data"`
}

func (c *SingleManga) GetResult() string {
	return c.Result
}

func (s *MangaService) GetManga(id string, params *GetMangaParams) (*SingleManga, error) {
	return s.GetMangaWithContext(context.Background(), id, params)
}

func (s *MangaService) GetMangaWithContext(ctx context.Context, id string, params *GetMangaParams) (*SingleManga, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaPath, id)

	u.RawQuery = EncodeParams(params)

	var l SingleManga
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

type MangaAggregateParams struct {
	Language []string `json:"translatedLanguage" url:"translatedLanguage[],omitempty"`
	Groups   string   `json:"groups" url:"groups[],omitempty"`
}

// Volumes: map of string to Volume
// Volumes tend to group up to 4 chapters together on MangaDex
// "none" will be returned for the latest chapters of an unfinished manga (probably, unless the volume is finished)
type MangaAggregate struct {
	Result  string                  `json:"result"`
	Volumes map[string]MangaVolumes `json:"volumes"`
}

// Volume: String version of valume. ex: "1". can be "none"
// Count: Total number of uploads for this volume
// Chapters: map of string Chapter number to aggregate
type MangaVolumes struct {
	Volume   string                      `json:"volume"`
	Count    int                         `json:"count"`
	Chapters map[string]ChapterAggregate `json:"chapters"`
}

// LatestId: Id for the most recently uploaded chapter for this Chapter
// Chapter: String version of chapter number. ex: "1". can maybe be "none"
// AdditionalChapters: string array of all other uploaded chapter Ids for this chapter
// Count: Total number of uploads for this chapter
type ChapterAggregate struct {
	LatestId           string   `json:"id"`
	Chapter            string   `json:"chapter"`
	AdditionalChapters []string `json:"others"`
	Count              int      `json:"count"`
}

func (ma *MangaAggregate) GetResult() string {
	return ma.Result
}

// Get the aggregate for a manga.
// Optimal if you want to get one Chapter Id per chapter Since ChapterList doesn't have this option
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-manga-aggregate
func (s *MangaService) GetMangaAggregate(mangaId string, params *MangaAggregateParams) (*MangaAggregate, error) {
	return s.GetMangaAggregateContext(context.Background(), mangaId, params)
}

// GetMangaListContext : GetMangaList with custom context.
func (s *MangaService) GetMangaAggregateContext(ctx context.Context, mangaId string, params *MangaAggregateParams) (*MangaAggregate, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaAggregatePath, mangaId)

	// Set query parameters
	u.RawQuery = EncodeParams(params)

	var l MangaAggregate
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// CheckIfMangaFollowed : Check if a user follows a manga.
func (s *MangaService) CheckIfMangaFollowed(id string) (bool, error) {
	return s.CheckIfMangaFollowedContext(context.Background(), id)
}

// CheckIfMangaFollowedContext : CheckIfMangaFollowed with custom context.
func (s *MangaService) CheckIfMangaFollowedContext(ctx context.Context, id string) (bool, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(CheckIfMangaFollowedPath, id)

	var r Response
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ToggleMangaFollowStatus :Toggle follow status for a manga.
func (s *MangaService) ToggleMangaFollowStatus(id string, toFollow bool) (*Response, error) {
	return s.ToggleMangaFollowStatusContext(context.Background(), id, toFollow)
}

// ToggleMangaFollowStatusContext  ToggleMangaFollowStatus with custom context.
func (s *MangaService) ToggleMangaFollowStatusContext(ctx context.Context, id string, toFollow bool) (*Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(ToggleMangaFollowPath, id)

	method := http.MethodPost // To follow
	if !toFollow {
		method = http.MethodDelete // To unfollow
	}

	var r Response
	err := s.client.RequestAndDecode(ctx, method, u.String(), nil, &r)
	return &r, err
}
