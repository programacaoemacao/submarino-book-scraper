package book

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/programacaoemacao/submarino-book-scraper/model"
	unicode "github.com/programacaoemacao/submarino-book-scraper/scraper/unicode"
	"go.uber.org/zap"
)

type bookScraper struct {
	collector *colly.Collector
	logger    *zap.SugaredLogger
}

func NewBookScraper(logger *zap.Logger) *bookScraper {
	collector := colly.NewCollector(
		colly.Async(false),
	)

	return &bookScraper{
		collector: collector,
		logger:    logger.Sugar(),
	}
}

func (bs *bookScraper) CollectDetailURLs(url string) (urls []string, totalItems uint, err error) {
	collector := bs.collector.Clone()

	urls = []string{}
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

	headers := http.Header{}
	headers.Add("authority", "www.submarino.com.br")
	headers.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Add("accept-language", "pt-BR,pt;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	headers.Add("cache-control", "max-age=0")
	headers.Add("if-none-match", "W/\"8d725-iqoXyTaY9BbpvmYyF1Fz/arZKws\"")
	headers.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"99\", \"Microsoft Edge\";v=\"115\", \"Chromium\";v=\"115\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"Windows\"")
	headers.Add("sec-fetch-dest", "document")
	headers.Add("sec-fetch-mode", "navigate")
	headers.Add("sec-fetch-site", "none")
	headers.Add("sec-fetch-user", "?1")
	headers.Add("upgrade-insecure-requests", "1")
	headers.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.188")

	functionError = collector.Request(http.MethodGet, url, nil, colly.NewContext(), headers)

	if functionError != nil {
		return nil, 0, functionError
	}

	return urls, totalOfItems, nil
}

func (bs *bookScraper) CollectDetail(detailURL string) (*model.Book, error) {
	collector := bs.collector.Clone()

	book := model.NewBook()
	var functionError error

	collector.OnXML(`//main//div[contains(@class,"image__WrapperImages")]//picture[contains(@class, "src__Picture")]/img`, func(x *colly.XMLElement) {
		book.CoverImageURL = x.Attr("src")
	})

	collector.OnXML(`//main//h1[contains(@class, "src__Title")]`, func(x *colly.XMLElement) {
		// The original title comes in lowercase
		bookPrefixRegex := regexp.MustCompile(`(?m)^livro[\s\-]+`)
		book.Title = unicode.Normalize(bookPrefixRegex.ReplaceAllString(x.Text, ""))
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
		book.PaymentCondition = unicode.Normalize(trimmed)
	})

	collector.OnXML(`//tr/td[text()="Autor"]/following-sibling::td`, func(x *colly.XMLElement) {
		authors := strings.Split(x.Text, ",")
		for i := 0; i < len(authors); i++ {
			authors[i] = unicode.Normalize(authors[i])
		}
		book.Authors = authors
	})

	collector.OnXML(`//div[contains(@class, "description__HTMLContent")]`, func(x *colly.XMLElement) {
		book.Description = unicode.Normalize(x.ChildText(`//*`))
	})

	collector.OnXML(`//tr/td[text()="Número de páginas"]/following-sibling::td`, func(x *colly.XMLElement) {
		pages, err := strconv.ParseUint(x.Text, 10, 64)
		if err == nil {
			book.Metadata.Pages = uint(pages)
		}
	})

	collector.OnXML(`//tr/td[text()="Idioma"]/following-sibling::td`, func(x *colly.XMLElement) {
		languages := strings.Split(x.Text, ",")
		for i := 0; i < len(languages); i++ {
			languages[i] = unicode.Normalize(languages[i])
		}
		book.Metadata.Languages = languages
	})

	collector.OnXML(`//tr/td[text()="Editora"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Publisher = unicode.Normalize(x.Text)
	})

	collector.OnXML(`//tr/td[text()="Data de Publicação"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.PublishDate = unicode.Normalize(x.Text)
	})

	collector.OnXML(`//tr/td[text()="ISBN-10"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.ISBN10 = unicode.Normalize(x.Text)
	})

	collector.OnXML(`//tr/td[text()="ISBN-13"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.ISBN13 = unicode.Normalize(x.Text)
	})

	collector.OnXML(`//tr/td[text()="Edição"]/following-sibling::td`, func(x *colly.XMLElement) {
		book.Metadata.Edition = unicode.Normalize(x.Text)
	})

	collector.OnError(func(r *colly.Response, err error) {
		functionError = err
		bs.logger.Errorln("error at scraping book:", err.Error())
	})

	collector.OnRequest(func(r *colly.Request) {
		bs.logger.Infof("scraping book on URL: %s", detailURL)
	})

	headers := http.Header{}
	headers.Add("authority", "www.submarino.com.br")
	headers.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Add("accept-language", "pt-BR,pt;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	headers.Add("cache-control", "max-age=0")
	headers.Add("if-none-match", "W/\"8d725-iqoXyTaY9BbpvmYyF1Fz/arZKws\"")
	headers.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"99\", \"Microsoft Edge\";v=\"115\", \"Chromium\";v=\"115\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"Windows\"")
	headers.Add("sec-fetch-dest", "document")
	headers.Add("sec-fetch-mode", "navigate")
	headers.Add("sec-fetch-site", "none")
	headers.Add("sec-fetch-user", "?1")
	headers.Add("upgrade-insecure-requests", "1")
	headers.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.188")

	functionError = collector.Request(http.MethodGet, detailURL, nil, colly.NewContext(), headers)

	if functionError != nil {
		return nil, functionError
	}

	return book, nil
}
