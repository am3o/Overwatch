package main

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/am3o/overwatch/pkg/client"
	"github.com/am3o/overwatch/pkg/config"
	"github.com/am3o/overwatch/pkg/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	config, err := config.Read("./resources/configuration.yml")
	if err != nil {
		panic(err)
	}

	uri, _ := url.Parse("https://www.mindfactory.de")
	shop := client.NewMindfactoryClient(uri, logger)

	collector := service.NewCollector()
	messanger, err := client.NewTelegram(config.Messanger.Token, config.Messanger.Ids)
	if err != nil {
		panic(err)
	}

	srv := service.New(collector, messanger, logger)

	scraper := service.NewScraper(shop, srv, logger)
	go scraper.Run(context.Background(), 5*time.Minute, config.Search)

	http.Handle("/api/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
