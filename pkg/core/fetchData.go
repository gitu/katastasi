package core

import (
	"bytes"
	"context"
	"github.com/gitu/katastasi/pkg/config"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func (k *Katastasi) GetPageStatus(env string, page string) config.PageStatus {
	p := config.PageStatus{
		Status:     config.Unknown,
		LastUpdate: time.Now(),
		Services:   map[string]*config.ServiceStatus{},
	}
	if _, f := k.Config.Environments[env]; !f {
		return p
	}
	if _, f := k.Config.Environments[env].StatusPages[page]; !f {
		return p
	}
	for _, service := range k.Config.Environments[env].StatusPages[page].Services {
		p.Services[service] = k.GetStatusOfService(env, service)
		if p.Services[service].Status.IsHigherThan(p.Status) {
			p.Status = p.Services[service].Status
		}
		if p.Services[service].LastUpdate.Before(p.LastUpdate) {
			p.LastUpdate = p.Services[service].LastUpdate
		}
	}
	return p
}

func (k *Katastasi) GetStatusOfService(env string, service string) *config.ServiceStatus {
	s := k.statusCache.Get(ToCacheKey(env, service))
	return s.Value()
}

func ToCacheKey(env string, service string) string {
	return env + "|" + service
}

func (k *Katastasi) loadServiceStatus(key string) *config.ServiceStatus {
	keyval := strings.SplitN(key, "|", 2)
	env := keyval[0]
	service := keyval[1]

	ret := &config.ServiceStatus{
		ID:         service,
		Status:     config.Unknown,
		LastUpdate: time.Now(),
		Components: make([]*config.ComponentStatus, 0),
	}

	if _, f := k.Config.Environments[env]; !f {
		return ret
	}
	if _, f := k.Config.Environments[env].Services[service]; !f {
		return ret
	}

	for _, component := range k.Config.Environments[env].Services[service].Components {
		nc := k.loadComponentStatus(component)
		if nc.Status.IsHigherThan(ret.Status) {
			ret.Status = nc.Status
		}
		ret.Components = append(ret.Components, nc)
	}

	return ret
}

func (k *Katastasi) loadComponentStatus(component *config.ServiceComponent) *config.ComponentStatus {
	nc := &config.ComponentStatus{
		Name:         component.Name,
		StatusString: "",
		Status:       config.OK,
	}
	q, foundQuery := k.Config.Queries[component.Query]
	if !foundQuery {
		nc.Status = config.Unknown
		nc.StatusString = "Query not found"
		slog.Warn("Query not found",
			"query", component.Query,
			"component", component.Name,
			"parameters", component.Parameters)
		return nc
	}

	var tpl bytes.Buffer
	if err := q.Execute(&tpl, component.Parameters); err != nil {
		nc.Status = config.Unknown
		nc.StatusString = "Error executing template: " + err.Error()
		slog.Error("Error executing template", "error", err)
		return nc
	}
	query := tpl.String()
	slog.Debug("built query",
		"query", query,
		"component", component.Name,
		"params", component.Parameters)

	v1api := v1.NewAPI(k.prometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	result, warnings, err := v1api.QueryRange(ctx, query, v1.Range{
		Start: time.Now().Add(-time.Hour * 6),
		End:   time.Now(),
		Step:  time.Second * 60,
	},
	)
	if err != nil {
		nc.StatusString = "Error querying Prometheus: " + err.Error()
		nc.Status = config.Unknown
		slog.Error("Error querying Prometheus", "error", err)
		return nc
	}
	if len(warnings) > 0 {
		slog.Warn("Warnings", "warnings", warnings)
	}

	if result.Type() == model.ValMatrix {
		matrix := result.(model.Matrix)
		for _, s := range matrix {

			if len(s.Values) == 0 {
				nc.StatusString = "No data returned"
				nc.Status = config.Warning
				slog.Info("No data returned", query, component.Name, component.Parameters)
				return nc
			}

			for _, c := range component.Conditions {
				match := buildComparator(c)
				if c.Duration != "" {
					duration, _ := time.ParseDuration(c.Duration)
					from := model.TimeFromUnix(time.Now().Add(-duration).Unix())
					for _, v := range s.Values {
						if v.Timestamp.Add(duration).After(from) {
							if match(v.Value) {
								nc.StatusString = component.Description
								nc.Status = c.Severity
							}
						}
					}
				} else {
					lastPair := s.Values[len(s.Values)-1]
					if match(lastPair.Value) {
						nc.StatusString = component.Description
						nc.Status = c.Severity
					}
				}

			}

		}

		if matrix.Len() == 0 {
			nc.StatusString = "No data returned"
			nc.Status = config.Warning
			slog.Info("No data returned", query, component.Name, component.Parameters)
			return nc
		}
	} else {
		nc.StatusString = "Query did return a wrong type"
		nc.Status = config.Unknown
		slog.Warn("Query did return a wrong type",
			"type", result.Type(),
			"query", query,
			"component", component.Name,
			"params", component.Parameters)
		return nc
	}

	return nc
}

func buildComparator(c *config.Condition) func(value model.SampleValue) bool {
	threshold, _ := strconv.ParseFloat(c.Threshold, 64)

	switch c.Condition {
	case config.Equal:
		return func(value model.SampleValue) bool {
			return float64(value) == threshold
		}
	case config.NotEqual:
		return func(value model.SampleValue) bool {
			return float64(value) != threshold
		}
	case config.GreaterThan:
		return func(value model.SampleValue) bool {
			return float64(value) > threshold
		}
	case config.LessThan:
		return func(value model.SampleValue) bool {
			return float64(value) < threshold
		}
	}
	return func(value model.SampleValue) bool {
		return false
	}
}
