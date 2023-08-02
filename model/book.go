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

/*
Initializes a new book instance with default slices values.
*/
func NewBook() *Book {
	/*
		By default, slices are initalized as nil and when it was marshalled to json, the value of a
		initialized slice will be `null`.

		I prefer to initalize the values as empty string arrays to get a empty array in a json representation
	*/
	book := &Book{
		Authors: []string{},
		Metadata: BookMetadata{
			Languages: []string{},
		},
	}
	return book
}
