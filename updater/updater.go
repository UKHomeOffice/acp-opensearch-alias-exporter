package updater

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.digital.homeoffice.gov.uk/acp/opensearch-reporter/models"
)

type PrometheusUpdater struct {
	gaugeVec *prometheus.GaugeVec
}

func NewPrometheusUpdater(namespace, name, help string) models.Updater {
	promAliasRate := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		},
		[]string{"domain"},
	)
	return &PrometheusUpdater{gaugeVec: promAliasRate}
}

func (p *PrometheusUpdater) Update(countRates []models.CountRate) {
	for _, countRate := range countRates {
		p.gaugeVec.WithLabelValues(countRate.Alias).Set(float64(countRate.Total))
	}
}
