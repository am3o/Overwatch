package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/am3o/overwatch/pkg/domain"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type Messanger interface {
	Message(text string) error
}

type Service struct {
	collector ProductCollector
	messanger Messanger
	logger    *zap.Logger
	products  *cache.Cache
}

func New(collector ProductCollector, messanger Messanger, logger *zap.Logger) Service {
	return Service{
		collector: collector,
		messanger: messanger,
		logger:    logger,
		products:  cache.New(10*time.Minute, 15*time.Minute),
	}
}

func (s Service) Notify(products domain.Products) {
	s.logger.With(zap.Int("Products", len(products))).Info("Notified about new received products")
	for _, product := range products {
		// s.collector.TrackProduct(product.Name, product.Price)

		if !strings.Contains(strings.ToLower(product.Name), "geforce") {
			continue
		}

		key := strconv.FormatUint(product.Hash(), 32)
		p, ok := s.products.Get(key)
		if !ok {
			s.products.Set(key, product, cache.DefaultExpiration)
			s.messanger.Message(fmt.Sprintf("Neue Karte steht zur Verfügung: %v € \n %v",
				product.Price, product.URI))
			continue
		}

		if obj, ok := p.(domain.Product); ok && obj.Price > product.Price {
			s.messanger.Message(fmt.Sprintf("Ist im Preis gesunken von %v € auf %v €: \n %v",
				obj.Price, product.Price, product.URI))
		}
	}
}

func (s Service) Check() []string {
	var products = make([]string, 0)
	for _, item := range s.products.Items() {
		product := item.Object.(domain.Product)
		products = append(products, fmt.Sprintf("%v: %v € \n %v", product.Name, product.Price, product.URI))
	}
	return products
}
