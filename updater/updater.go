package updater

import (
	"github.com/UKHomeOffice/acp-opensearch-alias-exporter/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
		[]string{"namespace"},
	)
	return &PrometheusUpdater{gaugeVec: promAliasRate}
}

func (p *PrometheusUpdater) UpdateRolloverAttemptFailures(aliasStatuses models.AliasStatuses) {
	for _, aliasStatus := range aliasStatuses {
		health := 0
		if aliasStatus.RolloverAttemptFailed {
			health = 1
		}
		p.gaugeVec.WithLabelValues(aliasStatus.Name).Set(float64(health))
	}
}

func (p *PrometheusUpdater) UpdateDocsAddedRate(statusChanges []models.StatusChange) {
	for _, statusChange := range statusChanges {
		p.gaugeVec.WithLabelValues(statusChange.Alias).Set(float64(statusChange.DocsAdded))
	}
}

func (p *PrometheusUpdater) UpdateIndexOperationFailures(statusChanges []models.StatusChange) {
	for _, statusChange := range statusChanges {
		p.gaugeVec.WithLabelValues(statusChange.Alias).Set(float64(statusChange.NewIndexOperationFailures))
	}
}
