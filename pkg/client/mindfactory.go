package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/am3o/overwatch/pkg/domain"
)

type MindfactoryClient struct {
	uri *url.URL
}

func NewMindfactoryClient (uri *url.URL) MindfactoryClient {
	return MindfactoryClient{
		uri: uri,
	}
}

func (client MindfactoryClient) Search(query string) (domain.Products, error) {
	client.uri.Path = fmt.Sprintf("/search_result.php/search_query/%s/article_per_page/5", url.QueryEscape(query))
	req, err := http.NewRequest(http.MethodGet, client.uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create search request: %W", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > http.StatusOK {
		return nil, fmt.Errorf("unvalid http status code response")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse response body: %w", err)
	}

	var products = make(domain.Products, 0)
	doc.Find(".pcontent").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".pname").Text()
		price := s.Find(".pprice").Text()
		available := strings.TrimSpace(s.Find(".pshipping .shipping1").Text())
		uri, _ := s.Find("a").Attr("href")

		product, err := domain.NewProduct(title, uri, price, available)
		if err != nil {
			fmt.Println(err)
			return
		}

		products = append(products, product)
	})

	return products, nil
}

