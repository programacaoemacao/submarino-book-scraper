# submarino-products-scraper

Submarino products scraper made in go.

This scraper was made using [Colly](https://github.com/gocolly/colly). It was implemented using a streamming approach (with [Observer pattern](https://refactoring.guru/design-patterns/observer)), so, you can stop the execution when you want, and the data won't be loss. 

The scraper has only 1 product type implemented: `Books`, but you can build your own scraper [strategy](https://refactoring.guru/design-patterns/strategy) to scrape another product category (games for example).

> **Important**: If you run the scraper again after a stopped/an error execution, the another execution data, will be appended to the same file creating some duplicates, so, if you wanna run it again, choose a new output file.

## Requirements

* Go >= 1.20
* Make

## Setup

Just run this command to download Go dependencies:

```shell
$ make setup
```

## Running the scraper

Execution flags:

| Short option | Long option | Description                     | Required |
|--------------|-------------|---------------------------------|----------|
| -u           | --url       | submarino url to scrape         | ✓        |
| -o           | --output    | output file (supported: [json]) | ✓        |

You can run the examples using make:

* Run example with best sellers books
    ```shell
    $ make run-example-best-sellers-books
    ```

* run-example-economics-books (Caution - It can lead to a `403 - Forbidden error`)
    ```shell
    $ make run-example-economics-books
    ```


### Choosing the URL

You can choose the books url to collect.

But the **URL can't have any query params**, for example:

* Best sellers: https://www.submarino.com.br/landingpage/trd-livros-mais-vendidos
* Economics: https://www.submarino.com.br/categoria/livros/administracao-negocios-e-economia
* Self Help: https://www.submarino.com.br/categoria/livros/autoajuda

You can choose any URL that respect this products grid layout:

![Submarino page with books in a grid layout](./docs/img/grid_layout_1.png "Submarino page")

## Known errors

* Sometimes we'll get a `403` from scraped website and the execution will stop.
