package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/programacaoemacao/submarino-book-scraper/model"
)

type collyScraper struct {
	collector *colly.Collector
}

func NewScraper() *collyScraper {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.183"),
	)

	return &collyScraper{
		collector: collector,
	}
}

func (c *collyScraper) CollectData() ([]model.Book, error) {
	urls := c.collectURLS()
	for _, url := range urls {
		c.scrapeBook(url)
	}

	return []model.Book{}, nil
}

func (c *collyScraper) collectURLS() []string {
	urls := []string{}
	c.collector.OnXML(`//a[contains(@class, "inStockCard__Link")]`, func(x *colly.XMLElement) {
		link := x.Attr("href")
		urls = append(urls, link)
	})

	c.collector.Visit("https://www.submarino.com.br/categoria/livros/ciencias-exatas")

	return urls
}

func (c *collyScraper) scrapeBook(url string) (*model.Book, error) {
	book := new(model.Book)
	var functionError error

	c.collector.OnXML(`//main//div[contains(@class,"image__WrapperImages")]//picture[contains(@class, "src__Picture")]/img`, func(x *colly.XMLElement) {
		imageURL := x.Attr("src")
		book.CoverImageURL = imageURL
	})

	c.collector.OnXML(`//main//h1[contains(@class, "src__Title")]`, func(x *colly.XMLElement) {
		title := x.Text
		book.Title = strings.Replace(title, "livro - ", "", 1)
	})

	c.collector.OnXML(`//main//div[contains(@class, "src__BestPrice")]`, func(x *colly.XMLElement) {
		price := x.Text
		priceRegex := regexp.MustCompile(`(?m)\d+`)
		matches := priceRegex.FindAllString(price, -1)
		fullPriceInCents := strings.Join(matches, "")
		priceInteger, err := strconv.Atoi(fullPriceInCents)
		if err == nil {
			book.PriceInCents = uint64(priceInteger)
		}
	})

	c.collector.OnXML(`//main//span[contains(@class, "src__RatingAverage")]`, func(x *colly.XMLElement) {
		rating := x.Text
		ratingFloat, err := strconv.ParseFloat(rating, 64)
		if err == nil {
			book.Rating.Average = ratingFloat
		}
	})

	c.collector.OnXML(`//main//div[contains(@class, "src__ProductInfo")]//span[contains(@class, "src__Count")]`, func(x *colly.XMLElement) {
		TotalOfRatings := x.Text
		totalOfRatingsRegex := regexp.MustCompile(`(?m)\d+`)
		match := totalOfRatingsRegex.FindString(TotalOfRatings)
		if match != "" {
			totalOfRatingUint, err := strconv.ParseUint(match, 10, 64)
			if err == nil {
				book.Rating.TotalOfRatings = uint(totalOfRatingUint)
			}
		}
	})

	c.collector.OnXML(`//main//p[contains(@class, "src__Text")]/strong`, func(x *colly.XMLElement) {
		paymentCondition := x.Text
		book.PaymentCondition = paymentCondition
	})

	c.collector.OnXML(`//tr/td[text()="Autor"]/following-sibling::td`, func(x *colly.XMLElement) {
		authors := strings.Split(x.Text, ",")
		book.Authors = authors
	})

	c.collector.OnXML(`//div[contains(@class, "description__HTMLContent")]`, func(x *colly.XMLElement) {
		descriptionParts := x.ChildTexts(`//*`)
		description := strings.Join(descriptionParts, "\n")
		book.Description = description
	})

	c.collector.OnError(func(r *colly.Response, err error) {
		functionError = err
	})

	c.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", url)
	})

	c.collector.Visit(url)

	return book, functionError
}
