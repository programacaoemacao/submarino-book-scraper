package book

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/programacaoemacao/submarino-book-scraper/model"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/consts"
	"github.com/programacaoemacao/submarino-book-scraper/scraper/utils"
	"go.uber.org/zap"
)

type bookScraper struct {
	collector *colly.Collector
	cookies   []*http.Cookie
	logger    *zap.SugaredLogger
}

func NewBookScraper(logger *zap.Logger) *bookScraper {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.183"),
		colly.Async(false),
	)

	sugaredLogger := logger.Sugar()

	return &bookScraper{
		collector: collector,
		logger:    sugaredLogger,
	}
}

func (bs *bookScraper) CollectData(baseURL string) ([]model.Book, error) {
	limit := consts.DefaultLimit
	offset := uint(0)
	books := []model.Book{}
	hasMoreItems := true
	currentBook := 0

	for hasMoreItems {
		bookListURL := utils.MountURL(baseURL, limit, offset)
		urls, totalItems, err := bs.scrapeBooksURLS(bookListURL)
		if err != nil {
			return nil, err
		}

		bs.logger.Debugf("limit: %d | offset: %d | total: %d\n", limit, offset, totalItems)

		for _, url := range urls {
			currentBook += 1
			bs.logger.Debugf("current progress: book %d of %d\n", currentBook, totalItems)
			book, err := bs.scrapeBook(url)
			if err == nil {
				books = append(books, *book)
			}
			bs.randomDelay()
		}

		hasMoreItems = totalItems > (offset + limit)
		offset += limit
	}

	return books, nil
}

func (bs *bookScraper) scrapeBooksURLS(booksPageURL string) ([]string, uint, error) {
	collector := bs.collector.Clone()

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
		bs.logger.Errorln("error at scraping books page: ", err.Error())
	})

	collector.Visit(booksPageURL)
	bs.cookies = collector.Cookies(booksPageURL)

	if functionError != nil {
		return nil, 0, functionError
	}

	return urls, totalOfItems, nil
}

func (bs *bookScraper) scrapeBook(url string) (*model.Book, error) {
	collector := bs.collector.Clone()

	book := model.NewBook()
	var functionError error

	// Clearing the cookies
	collector.SetCookies(url, []*http.Cookie{})

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
		bs.logger.Errorln("error at scraping book:", err.Error())
	})

	collector.OnRequest(func(r *colly.Request) {
		bs.logger.Infof("scraping book on URL: %s\n", url)
	})

	collector.SetCookies(url, bs.cookies)

	collector.Visit(url)

	if functionError != nil {
		return nil, functionError
	}

	return book, nil
}

func (bs *bookScraper) randomDelay() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 1 + rand.Intn(5) // n will be between 1 and 5
	for i := n; i > 0; i-- {
		bs.logger.Debugf("sleeping %d seconds ...\n", i)
		time.Sleep(time.Duration(1) * time.Second)
	}
}
