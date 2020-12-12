package service

import "github.com/prometheus/client_golang/prometheus"

type ProductCollector struct {
	products *prometheus.GaugeVec
}

func NewCollector() ProductCollector {
	return ProductCollector{
		products: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "",
			Name:      "products",
		}, []string{"name"}),
	}
}

func (p *ProductCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(p, descs)
}

func (p *ProductCollector) Collect(ch chan<- prometheus.Metric) {
	p.products.Collect(ch)
}

func (pc *ProductCollector) TrackProduct(name string, price float64) {
	pc.products.With(prometheus.Labels{
		"name": name,
	}).Add(price)
}
