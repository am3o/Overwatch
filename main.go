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
)

func main() {
	config, err := config.Read("./resources/configuration.yml")
	if err != nil {
		panic(err)
	}

	uri, _ := url.Parse("https://www.mindfactory.de")
	shop := client.NewMindfactoryClient(uri)

	collector := service.NewCollector()
	messanger, err := client.NewTelegram(config.Messanger.Token, config.Messanger.Ids)
	if err != nil {
		panic(err)
	}

	srv := service.New(collector, messanger)

	scraper := service.NewScraper(shop, srv)
	go scraper.Run(context.Background(), 5*time.Minute, config.Search)

	http.Handle("/api/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
