package service

import (
	"context"
	"time"

	"github.com/am3o/overwatch/pkg/domain"
	"go.uber.org/zap"
)

type ShopClient interface {
	Search(query string) (domain.Products, error)
}

type Notifier interface {
	Notify(products domain.Products)
}

type ShopScraper struct {
	client    ShopClient
	logger    *zap.Logger
	publisher Notifier
}

func NewScraper(client ShopClient, publisher Notifier, logger *zap.Logger) ShopScraper {
	return ShopScraper{
		client:    client,
		logger:    logger,
		publisher: publisher,
	}
}

func (scraper *ShopScraper) Run(ctx context.Context, interval time.Duration, query string) {
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		scraper.logger.Info("scrape the shop site, again")
		select {
		case <-ctx.Done():
			return
		default:
			products, err := scraper.client.Search(query)
			if err != nil {
				scraper.logger.Error("could not receive search response of the shop", zap.Error(err))
				continue
			}

			scraper.publisher.Notify(products)
		}
	}
}
