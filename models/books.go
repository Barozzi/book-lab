package model

// Books represents the top-level structure of the JSON response.
type Books struct {
	Kind       string `json:"kind"`
	TotalItems int    `json:"totalItems"`
	Items      []Item `json:"items"`
}

// Item represents individual items in the Items array.
type Item struct {
	Kind       string     `json:"kind"`
	ID         string     `json:"id"`
	Etag       string     `json:"etag"`
	SelfLink   string     `json:"selfLink"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
	SaleInfo   SaleInfo   `json:"saleInfo"`
	AccessInfo AccessInfo `json:"accessInfo"`
	SearchInfo SearchInfo `json:"searchInfo"`
}

// VolumeInfo contains detailed information about the volume.
type VolumeInfo struct {
	Title               string               `json:"title"`
	Authors             []string             `json:"authors"`
	Publisher           string               `json:"publisher"`
	PublishedDate       string               `json:"publishedDate"`
	Description         string               `json:"description"`
	IndustryIdentifiers []IndustryIdentifier `json:"industryIdentifiers"`
	ReadingModes        ReadingModes         `json:"readingModes"`
	PageCount           int                  `json:"pageCount"`
	PrintType           string               `json:"printType"`
	Categories          []string             `json:"categories"`
	MaturityRating      string               `json:"maturityRating"`
	AllowAnonLogging    bool                 `json:"allowAnonLogging"`
	ContentVersion      string               `json:"contentVersion"`
	PanelizationSummary PanelizationSummary  `json:"panelizationSummary"`
	ImageLinks          ImageLinks           `json:"imageLinks"`
	Language            string               `json:"language"`
	PreviewLink         string               `json:"previewLink"`
	InfoLink            string               `json:"infoLink"`
	CanonicalVolumeLink string               `json:"canonicalVolumeLink"`
}

// IndustryIdentifier represents industry identifiers for the book.
type IndustryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// ReadingModes represents reading modes for the book.
type ReadingModes struct {
	Text  bool `json:"text"`
	Image bool `json:"image"`
}

// PanelizationSummary provides summary information about panelization.
type PanelizationSummary struct {
	ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
	ContainsImageBubbles bool `json:"containsImageBubbles"`
}

// ImageLinks provides URLs for images related to the book.
type ImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

// SaleInfo contains sale information about the book.
type SaleInfo struct {
	Country     string `json:"country"`
	Saleability string `json:"saleability"`
	IsEbook     bool   `json:"isEbook"`
}

// AccessInfo contains access information about the book.
type AccessInfo struct {
	Country                string   `json:"country"`
	Viewability            string   `json:"viewability"`
	Embeddable             bool     `json:"embeddable"`
	PublicDomain           bool     `json:"publicDomain"`
	TextToSpeechPermission string   `json:"textToSpeechPermission"`
	Epub                   EpubInfo `json:"epub"`
	Pdf                    PdfInfo  `json:"pdf"`
	WebReaderLink          string   `json:"webReaderLink"`
	AccessViewStatus       string   `json:"accessViewStatus"`
	QuoteSharingAllowed    bool     `json:"quoteSharingAllowed"`
}

// EpubInfo contains epub-specific information.
type EpubInfo struct {
	IsAvailable bool `json:"isAvailable"`
}

// PdfInfo contains pdf-specific information.
type PdfInfo struct {
	IsAvailable  bool   `json:"isAvailable"`
	AcsTokenLink string `json:"acsTokenLink"`
}

// SearchInfo contains search information about the book.
type SearchInfo struct {
	TextSnippet string `json:"textSnippet"`
}
