package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

func TestGetRepository(t *testing.T) {
	tests := []struct {
		name    string
		cfg     configs.RepositoryConfig
		want    Repository
		wantErr bool
	}{
		{
			name: "OK RAM default",
			cfg:  configs.RepositoryConfig{},
			want: &ramRepository{
				counterStorage: map[string]metrics.Counter{},
				gaugeStorage:   map[string]metrics.Gauge{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRepository(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepository(%v) error = %v, wantErr %v", tt.cfg, err, tt.wantErr)
			}
			assert.EqualValuesf(t, tt.want, got, "GetRepository(%v)", tt.cfg)
		})
	}
}
