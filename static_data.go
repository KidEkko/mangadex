package mangodex

// For Param ordering
const (
	AscendingOrder  = "asc"
	DescendingOrder = "desc"
)

// Publication demographic
const (
	ShonenDemographic = "shounen"
	ShoujoDemographic = "shoujo"
	JoseiDemographic  = "josei"
	SeinenDemograpic  = "seinen"
)

// Manga publication status
const (
	OngoingStatus   = "ongoing"
	CompletedStatus = "completed"
	HiatusStatus    = "hiatus"
	CancelledStatus = "cancelled"
	NoStatus        = "none"
)

// Manga reading status
const (
	Reading    = "reading"
	OnHold     = "on_hold"
	PlanToRead = "plan_to_read"
	Dropped    = "dropped"
	ReReading  = "re_reading"
	Completed  = "completed"
)

// Manga content rating
const (
	Safe       = "safe"
	Suggestive = "suggestive"
	Erotica    = "erotica"
	Porn       = "pornographic"
)

// Relationship types. Useful for reference expansions
const (
	MangaRel           = "manga"
	ChapterRel         = "chapter"
	CoverArtRel        = "cover_art"
	AuthorRel          = "author"
	ArtistRel          = "artist"
	ScanlationGroupRel = "scanlation_group"
	TagRel             = "tag"
	UserRel            = "user"
	CustomListRel      = "custom_list"
)

// Includes enums, use in your arrays
const (
	IncManga   = "manga"
	IncCover   = "cover_art"
	IncAuthor  = "author"
	IncArtist  = "artist"
	IncTag     = "tag"
	IncCreator = "creator"
)
