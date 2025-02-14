package consul

import (
	"net/url"
	"time"

	"github.com/grafana/agent/component"
	"github.com/grafana/agent/component/discovery"
	"github.com/grafana/agent/component/prometheus/exporter"
	"github.com/grafana/agent/pkg/integrations"
	"github.com/grafana/agent/pkg/integrations/consul_exporter"
)

func init() {
	component.Register(component.Registration{
		Name:    "prometheus.exporter.consul",
		Args:    Arguments{},
		Exports: exporter.Exports{},
		Build:   exporter.NewWithTargetBuilder(createExporter, "consul", customizeTarget),
	})
}

func createExporter(opts component.Options, args component.Arguments) (integrations.Integration, error) {
	a := args.(Arguments)
	return a.Convert().NewIntegration(opts.Logger)
}

func customizeTarget(baseTarget discovery.Target, args component.Arguments) []discovery.Target {
	a := args.(Arguments)
	target := baseTarget

	url, err := url.Parse(a.Server)
	if err != nil {
		return []discovery.Target{target}
	}

	target["instance"] = url.Host
	return []discovery.Target{target}
}

// DefaultArguments holds the default settings for the consul_exporter exporter.
var DefaultArguments = Arguments{
	Server:        "http://localhost:8500",
	Timeout:       500 * time.Millisecond,
	AllowStale:    true,
	KVFilter:      ".*",
	HealthSummary: true,
}

// Arguments controls the consul_exporter exporter.
type Arguments struct {
	Server             string        `river:"server,attr,optional"`
	CAFile             string        `river:"ca_file,attr,optional"`
	CertFile           string        `river:"cert_file,attr,optional"`
	KeyFile            string        `river:"key_file,attr,optional"`
	ServerName         string        `river:"server_name,attr,optional"`
	Timeout            time.Duration `river:"timeout,attr,optional"`
	InsecureSkipVerify bool          `river:"insecure_skip_verify,attr,optional"`
	RequestLimit       int           `river:"concurrent_request_limit,attr,optional"`
	AllowStale         bool          `river:"allow_stale,attr,optional"`
	RequireConsistent  bool          `river:"require_consistent,attr,optional"`

	KVPrefix      string `river:"kv_prefix,attr,optional"`
	KVFilter      string `river:"kv_filter,attr,optional"`
	HealthSummary bool   `river:"generate_health_summary,attr,optional"`
}

// UnmarshalRiver implements River unmarshalling for Arguments.
func (a *Arguments) UnmarshalRiver(f func(interface{}) error) error {
	*a = DefaultArguments

	type args Arguments
	return f((*args)(a))
}

func (a *Arguments) Convert() *consul_exporter.Config {
	return &consul_exporter.Config{
		Server:             a.Server,
		CAFile:             a.CAFile,
		CertFile:           a.CertFile,
		KeyFile:            a.KeyFile,
		ServerName:         a.ServerName,
		Timeout:            a.Timeout,
		InsecureSkipVerify: a.InsecureSkipVerify,
		RequestLimit:       a.RequestLimit,
		AllowStale:         a.AllowStale,
		RequireConsistent:  a.RequireConsistent,
		KVPrefix:           a.KVPrefix,
		KVFilter:           a.KVFilter,
		HealthSummary:      a.HealthSummary,
	}
}
