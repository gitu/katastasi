package config_test

import (
	"github.com/gitu/katastasi/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseServiceComponents(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    []*config.ServiceComponent
		wantErr bool
	}{
		{name: "basic complex component", args: args{data: `- name: "Availability Pods"
  query: "unhealthy_pods"
  description: "There are unhealthy pods in the environment"
  conditions:
    - severity: "warning"
      condition: "gt"
      threshold: "0"
    - severity: "critical"
      condition: "gt"
      threshold: "0"
      duration: "5m"`}, want: []*config.ServiceComponent{
			{
				Name:        "Availability Pods",
				Query:       "unhealthy_pods",
				Description: "There are unhealthy pods in the environment",
				Conditions: []*config.Condition{
					{
						Severity:  config.Warning,
						Condition: config.GreaterThan,
						Threshold: "0",
					}, {
						Severity:  config.Critical,
						Condition: config.GreaterThan,
						Threshold: "0",
						Duration:  "5m",
					},
				},
			},
		}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.ParseServiceComponents([]byte(tt.args.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseServiceComponents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}
