package prometheusconvert

import (
	"time"

	"github.com/grafana/agent/component/discovery"
	"github.com/grafana/agent/component/prometheus/scrape"
	promconfig "github.com/prometheus/prometheus/config"
	promdiscovery "github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/storage"
)

func toScrapeArguments(scrapeConfig *promconfig.ScrapeConfig) *scrape.Arguments {
	if scrapeConfig == nil {
		return nil
	}

	return &scrape.Arguments{
		Targets:               getTargets(scrapeConfig),
		ForwardTo:             []storage.Appendable{}, // TODO
		JobName:               scrapeConfig.JobName,
		HonorLabels:           scrapeConfig.HonorLabels,
		HonorTimestamps:       scrapeConfig.HonorTimestamps,
		Params:                scrapeConfig.Params,
		ScrapeInterval:        time.Duration(scrapeConfig.ScrapeInterval),
		ScrapeTimeout:         time.Duration(scrapeConfig.ScrapeTimeout),
		MetricsPath:           scrapeConfig.MetricsPath,
		Scheme:                scrapeConfig.Scheme,
		BodySizeLimit:         scrapeConfig.BodySizeLimit,
		SampleLimit:           scrapeConfig.SampleLimit,
		TargetLimit:           scrapeConfig.TargetLimit,
		LabelLimit:            scrapeConfig.LabelLimit,
		LabelNameLengthLimit:  scrapeConfig.LabelNameLengthLimit,
		LabelValueLengthLimit: scrapeConfig.LabelValueLengthLimit,
		HTTPClientConfig:      *toHttpClientConfig(&scrapeConfig.HTTPClientConfig),
		ExtraMetrics:          false,
		Clustering:            scrape.Clustering{Enabled: false},
	}
}

func getTargets(scrapeConfig *promconfig.ScrapeConfig) []discovery.Target {
	targets := []discovery.Target{}

	for _, serviceDiscoveryConfig := range scrapeConfig.ServiceDiscoveryConfigs {
		switch sdc := serviceDiscoveryConfig.(type) {
		case promdiscovery.StaticConfig:
			for _, target := range sdc {
				for _, labelSet := range target.Targets {
					for labelName, labelValue := range labelSet {
						targets = append(targets, map[string]string{string(labelName): string(labelValue)})
					}
				}
			}
		}
	}

	return targets
}
