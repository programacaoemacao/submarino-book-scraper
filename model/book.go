package model

type Book struct {
	CoverImageURL    string       `json:"cover_image_url"`
	Title            string       `json:"title"`
	Authors          []string     `json:"authors"`
	Rating           BookRating   `json:"rating"`
	PriceInCents     uint64       `json:"price_in_cents"`
	PaymentCondition string       `json:"payment_condition"`
	Description      string       `json:"description"`
	Metadata         BookMetadata `json:"metadata"`
}

type BookRating struct {
	Average        float64 `json:"average"`
	TotalOfRatings uint    `json:"total_of_ratings"`
}

type BookMetadata struct {
	Pages       uint     `json:"pages"`
	Languages   []string `json:"languages"`
	Publisher   string   `json:"publisher"`
	PublishDate string   `json:"publish_date"`
	ISBN10      string   `json:"isbn_10"`
	ISBN13      string   `json:"isbn_13"`
	Edition     string   `json:"edition"`
}
