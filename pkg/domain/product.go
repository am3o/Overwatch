package domain

import (
	"hash/fnv"
	"net/url"
	"strconv"
	"strings"
)

type Products []Product

type Product struct {
	Name      string
	URI       *url.URL
	Price     float64
	Available bool
}

func NewProduct(name string, uri string, price string, available string) (Product, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return Product{}, err
	}

	price = strings.ReplaceAll(strings.ReplaceAll(price, ".", ""), ",", ".")
	p, err := strconv.ParseFloat(price[5:len(price)-1], 64)
	if err != nil {
		return Product{}, err
	}

	return Product{
		Name:      strings.TrimSpace(name),
		URI:       u,
		Price:     p,
		Available: available == "Lagernd",
	}, nil
}

func (p Product) Hash() uint64 {
	h := fnv.New32a()
	_, err := h.Write([]byte(p.Name))
	if err != nil {
		return 0
	}
	return uint64(h.Sum32())
}
