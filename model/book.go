package model

type Book struct {
	CoverImageURL    string
	Title            string
	Authors          []string
	Rating           BookRating
	PriceInCents     uint64
	PaymentCondition string
	Description      string
	Metadata         BookMetadata
}

type BookRating struct {
	Average        float64
	TotalOfRatings uint
}

type BookMetadata struct {
	Pages       uint
	Languages   []string
	Publisher   string
	PublishDate string
	Dimension   BookDimension
	ISBN10      string
	ISBN13      string
}

type BookDimension struct {
	Height float64
	Width  float64
	Length float64
}
