help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## Download go dependencies
	@go mod tidy

run-example-best-sellers-books: ## Run scraper on best sellers books section - Quick execution
	@go run main.go -u=https://www.submarino.com.br/landingpage/trd-livros-mais-vendidos -o=best_sellers.json

run-example-economics-books: ## Run scraper on economics books section - Long execution - It can lead to a 403 - Forbidden error
	@go run main.go -u=https://www.submarino.com.br/categoria/livros/administracao-negocios-e-economia -o=economics.json