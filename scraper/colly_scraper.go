package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/programacaoemacao/submarino-book-scraper/model"
)

type bookScraper struct {
	collector *colly.Collector
}

func NewBookScraper() *bookScraper {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.183"),
		colly.Async(false),
	)

	return &bookScraper{
		collector: collector,
	}
}

func (c *bookScraper) CollectData(baseURL string) ([]model.Book, error) {
	limit := defaultLimit
	offset := uint(0)
	books := []model.Book{}
	hasMoreItems := true

	for hasMoreItems {
		bookListURL := mountURL(baseURL, limit, offset)
		urls, totalItems, err := c.scrapeBooksURLS(bookListURL)
		if err != nil {
			return nil, err
		}

		fmt.Printf("limit: %d | offset: %d | total: %d\n", limit, offset, totalItems)

		for _, url := range urls {
			book, err := c.scrapeBook(url)
			if err == nil {
				books = append(books, *book)
			}
		}
		offset += limit
		hasMoreItems = totalItems > (offset + limit)
	}

	return books, nil
}

func (c *bookScraper) scrapeBooksURLS(booksPageURL string) ([]string, uint, error) {
	collector := c.collector.Clone()

	urls := []string{}
	var totalOfItems uint
	var functionError error

	collector.OnXML(`//div[contains(@class, "inStockCard__Wrapper")]/a`, func(x *colly.XMLElement) {
		host := x.Request.URL.Host
		// It's necessary to set https protocol, otherwise, colly will use the http protocol, and it won't work
		url := "https://" + host + x.Attr("href")
		urls = append(urls, url)
	})

	collector.OnXML(`//span[contains(@class, "grid-area__TotalText")]`, func(x *colly.XMLElement) {
		totalItemsRegex := regexp.MustCompile(`(?m)\d+`)
		matches := totalItemsRegex.FindAllString(x.Text, -1)
		totalItems := strings.Join(matches, "")
		total, err := strconv.ParseUint(totalItems, 10, 64)
		if err != nil || total == 0 {
			functionError = fmt.Errorf("can't get total of pages")
		}
		totalOfItems = uint(total)
	})

	collector.OnError(func(r *colly.Response, err error) {
		functionError = err
	})

	collector.Visit(booksPageURL)

	if functionError != nil {
		return nil, 0, functionError
	}

	return urls, totalOfItems, nil
}

func (c *bookScraper) scrapeBook(url string) (*model.Book, error) {
	collector := c.collector.Clone()

	book := new(model.Book)
	var functionError error

	collector.OnXML(`//main//div[contains(@class,"image__WrapperImages")]//picture[contains(@class, "src__Picture")]/img`, func(x *colly.XMLElement) {
		book.CoverImageURL = x.Attr("src")
	})

	collector.OnXML(`//main//h1[contains(@class, "src__Title")]`, func(x *colly.XMLElement) {
		// The original title comes in lowercase
		bookPrefixRegex := regexp.MustCompile(`(?m)^livro[\s\-]+`)
		book.Title = bookPrefixRegex.ReplaceAllString(x.Text, "")
	})

	collector.OnXML(`//main//div[contains(@class, "src__BestPrice")]`, func(x *colly.XMLElement) {
		priceRegex := regexp.MustCompile(`(?m)\d+`)
		matches := priceRegex.FindAllString(x.Text, -1)
		fullPriceInCents := strings.Join(matches, "")
		priceInCents, err := strconv.Atoi(fullPriceInCents)
		if err == nil {
			book.PriceInCents = uint64(priceInCents)
		}
	})

	collector.OnXML(`//main//span[contains(@class, "src__RatingAverage")]`, func(x *colly.XMLElement) {
		ratingFloat, err := strconv.ParseFloat(x.Text, 64)
		if err == nil {
			book.Rating.Average = ratingFloat
		}
	})

	collector.OnXML(`//main//div[contains(@class, "src__ProductInfo")]//span[contains(@class, "src__Count")]`, func(x *colly.XMLElement) {
		totalOfRatingsRegex := regexp.MustCompile(`(?m)\d+`)
		match := totalOfRatingsRegex.FindString(x.Text)
		if match != "" {
			totalOfRatingUint, err := strconv.ParseUint(match, 10, 64)
			if err == nil {
				book.Rating.TotalOfRatings = uint(totalOfRatingUint)
			}
		}
	})

	collector.OnXML(`//main//p[contains(@class, "src__Text-")]`, func(x *colly.XMLElement) {
		nbSpaceRegex := regexp.MustCompile(`(?m)\p{Z}`)
		withoutNBSpace := nbSpaceRegex.ReplaceAllString(x.Text, " ")
		trimmed := strings.TrimSpace(withoutNBSpace)
		book.PaymentCondition = trimmed
	})

	collector.OnXML(`//tr/td[text()="Autor"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Authors = strings.Split(x.Text, ",")
	})

	collector.OnXML(`//div[contains(@class, "description__HTMLContent")]`, func(x *colly.XMLElement) {
		book.Description = x.ChildText(`//*`)
	})

	collector.OnXML(`//tr/td[text()="Número de páginas"]/following-sibling::td`, func(x *colly.XMLElement) {
		pages, err := strconv.ParseUint(x.Text, 10, 64)
		if err == nil {
			book.Metadata.Pages = uint(pages)
		}
	})

	collector.OnXML(`//tr/td[text()="Idioma"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Languages = strings.Split(x.Text, ",")
	})

	collector.OnXML(`//tr/td[text()="Editora"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Publisher = x.Text
	})

	collector.OnXML(`//tr/td[text()="Data de Publicação"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.PublishDate = x.Text
	})

	collector.OnXML(`//tr/td[text()="ISBN-10"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.ISBN10 = x.Text
	})

	collector.OnXML(`//tr/td[text()="ISBN-13"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.ISBN13 = x.Text
	})

	collector.OnXML(`//tr/td[text()="Edição"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Edition = x.Text
	})

	collector.OnError(func(r *colly.Response, err error) {
		functionError = err
		fmt.Println("error at scraping book on URL:", url)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("scraping book on URL: %s\n", url)
	})

	collector.Visit(url)

	if functionError != nil {
		return nil, functionError
	}

	return book, nil
}
