package model

// Books represents the top-level structure of the JSON response.
type GoogleBookResponse struct {
	Kind         string           `json:"kind"`
	TotalItems   int              `json:"totalItems"`
	HasMorePages bool             `json:"hasMorePages"`
	Items        []GoogleBookItem `json:"items"`
}

// GoogleBookItem represents individual items in the Items array.
type GoogleBookItem struct {
	Kind       string               `json:"kind"`
	ID         string               `json:"id"`
	Etag       string               `json:"etag"`
	SelfLink   string               `json:"selfLink"`
	VolumeInfo GoogleBookVolumeInfo `json:"volumeInfo"`
	SaleInfo   GoogleBookSaleInfo   `json:"saleInfo"`
	AccessInfo GoogleBookAccessInfo `json:"accessInfo"`
	SearchInfo GoogleBookSearchInfo `json:"searchInfo"`
}

// GoogleBookVolumeInfo contains detailed information about the volume.
type GoogleBookVolumeInfo struct {
	Title               string                         `json:"title"`
	Authors             []string                       `json:"authors"`
	Publisher           string                         `json:"publisher"`
	PublishedDate       string                         `json:"publishedDate"`
	Description         string                         `json:"description"`
	IndustryIdentifiers []GoogleBookIndustryIdentifier `json:"industryIdentifiers"`
	ReadingModes        GoogleBookReadingModes         `json:"readingModes"`
	PageCount           int                            `json:"pageCount"`
	PrintType           string                         `json:"printType"`
	Categories          []string                       `json:"categories"`
	MaturityRating      string                         `json:"maturityRating"`
	AllowAnonLogging    bool                           `json:"allowAnonLogging"`
	ContentVersion      string                         `json:"contentVersion"`
	PanelizationSummary GoogleBookPanelizationSummary  `json:"panelizationSummary"`
	ImageLinks          GoogleBookImageLinks           `json:"imageLinks"`
	Language            string                         `json:"language"`
	PreviewLink         string                         `json:"previewLink"`
	InfoLink            string                         `json:"infoLink"`
	CanonicalVolumeLink string                         `json:"canonicalVolumeLink"`
}

// GoogleBookIndustryIdentifier represents industry identifiers for the book.
type GoogleBookIndustryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// GoogleBookReadingModes represents reading modes for the book.
type GoogleBookReadingModes struct {
	Text  bool `json:"text"`
	Image bool `json:"image"`
}

// GoogleBookPanelizationSummary provides summary information about panelization.
type GoogleBookPanelizationSummary struct {
	ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
	ContainsImageBubbles bool `json:"containsImageBubbles"`
}

// GoogleBookImageLinks provides URLs for images related to the book.
type GoogleBookImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

// GoogleBookSaleInfo contains sale information about the book.
type GoogleBookSaleInfo struct {
	Country     string `json:"country"`
	Saleability string `json:"saleability"`
	IsEbook     bool   `json:"isEbook"`
}

// GoogleBookAccessInfo contains access information about the book.
type GoogleBookAccessInfo struct {
	Country                string             `json:"country"`
	Viewability            string             `json:"viewability"`
	Embeddable             bool               `json:"embeddable"`
	PublicDomain           bool               `json:"publicDomain"`
	TextToSpeechPermission string             `json:"textToSpeechPermission"`
	Epub                   GoogleBookEpubInfo `json:"epub"`
	Pdf                    GoogleBookPdfInfo  `json:"pdf"`
	WebReaderLink          string             `json:"webReaderLink"`
	AccessViewStatus       string             `json:"accessViewStatus"`
	QuoteSharingAllowed    bool               `json:"quoteSharingAllowed"`
}

// GoogleBookEpubInfo contains epub-specific information.
type GoogleBookEpubInfo struct {
	IsAvailable bool `json:"isAvailable"`
}

// GoogleBookPdfInfo contains pdf-specific information.
type GoogleBookPdfInfo struct {
	IsAvailable  bool   `json:"isAvailable"`
	AcsTokenLink string `json:"acsTokenLink"`
}

// GoogleBookSearchInfo contains search information about the book.
type GoogleBookSearchInfo struct {
	TextSnippet string `json:"textSnippet"`
}
