package agent

import (
	"fmt"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

func FormatURL(addr string, m metrics.Metric) string {
	return fmt.Sprintf("http://%v/update/%v/%v/%v", addr, m.GetType(), m.GetName(), m.GetValue())
}
