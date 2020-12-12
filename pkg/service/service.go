package service

import (
	"fmt"
	"strings"

	"github.com/am3o/overwatch/pkg/domain"
)

type Messanger interface {
	Message(text string) error
}

type Service struct {
	collector ProductCollector
	messanger Messanger
	products map[uint32]domain.Product
}

func New(collector ProductCollector, messanger Messanger) Service {
	return Service{
		collector: collector,
		messanger: messanger,
		products: make(map[uint32]domain.Product),
	}
}

func (s Service) Notify(products domain.Products) {
	for _, product := range products {
		//s.collector.TrackProduct(product.Name, product.Price)

		if strings.Contains(product.Name, "Aqua") {
			continue
		}

		p, ok := s.products[product.Hash()]
		if !ok {
			s.products[product.Hash()] = product
			s.messanger.Message(fmt.Sprintf("Neue Karte steht zur Verfügung: %v € \n %v",
				product.Price, product.URI))
			continue
		}

		if p.Price > product.Price {
			s.messanger.Message(fmt.Sprintf("Ist im Preis gesunken von %v € auf %v €: \n %v",
				p.Price, product.Price, product.URI))
		}
	}
}
