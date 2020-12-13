package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
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

	path, exists := os.LookupEnv("CONFIGURATION")
	if !exists {
		path = "./resources/configuration.yml"
	}

	config, err := config.Read(path)
	if err != nil {
		panic(err)
	}

	uri, _ := url.Parse("https://www.mindfactory.de")
	shop := client.NewMindfactoryClient(uri, logger)

	collector := service.NewCollector()
	messanger, err := client.NewTelegram(config.Messanger.Token, config.Messanger.Ids, logger)
	if err != nil {
		panic(err)
	}

	srv := service.New(collector, messanger, logger)

	scraper := service.NewScraper(shop, srv, logger)
	go scraper.Run(context.Background(), time.Minute, config.Search)
	go messanger.Update(srv.Check)

	http.Handle("/api/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
