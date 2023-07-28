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
	urls, _ := c.scrapeBooksURLS("")
	for _, url := range urls {
		c.scrapeBook(url)
	}

	return []model.Book{}, nil
}

func (c *collyScraper) scrapeBooksURLS(booksPageURL string) ([]string, error) {
	urls := []string{}
	var functionError error

	c.collector.OnXML(`//a[contains(@class, "inStockCard__Link")]`, func(x *colly.XMLElement) {
		host := x.Request.URL.Host
		urls = append(urls, host+x.Attr("href"))
	})

	c.collector.OnError(func(r *colly.Response, err error) {
		functionError = err
	})

	c.collector.Visit(booksPageURL)

	if functionError != nil {
		return nil, functionError
	}

	return urls, nil
}

func (c *collyScraper) scrapeBook(url string) (*model.Book, error) {
	book := new(model.Book)
	var functionError error

	c.collector.OnXML(`//main//div[contains(@class,"image__WrapperImages")]//picture[contains(@class, "src__Picture")]/img`, func(x *colly.XMLElement) {
		book.CoverImageURL = x.Attr("src")
	})

	c.collector.OnXML(`//main//h1[contains(@class, "src__Title")]`, func(x *colly.XMLElement) {
		// The original title comes in lowercase
		bookPrefixRegex := regexp.MustCompile(`(?m)^livro[\s\-]+`)
		book.Title = bookPrefixRegex.ReplaceAllString(x.Text, "")
	})

	c.collector.OnXML(`//main//div[contains(@class, "src__BestPrice")]`, func(x *colly.XMLElement) {
		priceRegex := regexp.MustCompile(`(?m)\d+`)
		matches := priceRegex.FindAllString(x.Text, -1)
		fullPriceInCents := strings.Join(matches, "")
		priceInCents, err := strconv.Atoi(fullPriceInCents)
		if err == nil {
			book.PriceInCents = uint64(priceInCents)
		}
	})

	c.collector.OnXML(`//main//span[contains(@class, "src__RatingAverage")]`, func(x *colly.XMLElement) {
		ratingFloat, err := strconv.ParseFloat(x.Text, 64)
		if err == nil {
			book.Rating.Average = ratingFloat
		}
	})

	c.collector.OnXML(`//main//div[contains(@class, "src__ProductInfo")]//span[contains(@class, "src__Count")]`, func(x *colly.XMLElement) {
		totalOfRatingsRegex := regexp.MustCompile(`(?m)\d+`)
		match := totalOfRatingsRegex.FindString(x.Text)
		if match != "" {
			totalOfRatingUint, err := strconv.ParseUint(match, 10, 64)
			if err == nil {
				book.Rating.TotalOfRatings = uint(totalOfRatingUint)
			}
		}
	})

	c.collector.OnXML(`//main//p[contains(@class, "src__Text-")]`, func(x *colly.XMLElement) {
		nbSpaceRegex := regexp.MustCompile(`(?m)\p{Z}`)
		withoutNBSpace := nbSpaceRegex.ReplaceAllString(x.Text, " ")
		trimmed := strings.TrimSpace(withoutNBSpace)
		book.PaymentCondition = trimmed
	})

	c.collector.OnXML(`//tr/td[text()="Autor"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Authors = strings.Split(x.Text, ",")
	})

	c.collector.OnXML(`//div[contains(@class, "description__HTMLContent")]`, func(x *colly.XMLElement) {
		book.Description = x.ChildText(`//*`)
	})

	c.collector.OnXML(`//tr/td[text()="Número de páginas"]/following-sibling::td`, func(x *colly.XMLElement) {
		pages, err := strconv.ParseUint(x.Text, 10, 64)
		if err == nil {
			book.Metadata.Pages = uint(pages)
		}
	})

	c.collector.OnXML(`//tr/td[text()="Idioma"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Languages = strings.Split(x.Text, ",")
	})

	c.collector.OnXML(`//tr/td[text()="Editora"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Publisher = x.Text
	})

	c.collector.OnXML(`//tr/td[text()="Data de Publicação"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.PublishDate = x.Text
	})

	c.collector.OnXML(`//tr/td[text()="ISBN-10"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.ISBN10 = x.Text
	})

	c.collector.OnXML(`//tr/td[text()="ISBN-13"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.ISBN13 = x.Text
	})

	c.collector.OnXML(`//tr/td[text()="Edição"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Edition = x.Text
	})

	c.collector.OnError(func(r *colly.Response, err error) {
		functionError = err
	})

	c.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", url)
	})

	c.collector.Visit(url)

	if functionError != nil {
		return nil, functionError
	}

	return book, nil
}
