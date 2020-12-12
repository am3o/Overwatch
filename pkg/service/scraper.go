package service

import (
	"context"
	"fmt"
	"time"

	"github.com/am3o/overwatch/pkg/domain"
)

type ShopClient interface {
	Search(query string) (domain.Products,error)
}

type Notifier interface {
	Notify(products domain.Products)
}

type ShopScraper struct {
	client    ShopClient
	publisher Notifier
}

func NewScraper(client ShopClient, publisher Notifier) ShopScraper {
	return ShopScraper{
		client:    client,
		publisher: publisher,
	}
}

func (scraper *ShopScraper) Run(ctx context.Context, interval time.Duration, query string) {
	ticker := time.NewTicker(interval)
	for ;;<-ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			products, err := scraper.client.Search(query)
			if err != nil {
				fmt.Println(err)
			}

			scraper.publisher.Notify(products)
		}
	}
}
