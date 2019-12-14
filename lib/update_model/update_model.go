package update_model

import (
	"os"
	"time"

	"fmt"

	mkr "github.com/mackerelio/mackerel-client-go"
)

func postEvaluatedMetricsToMackerel(metricNamePrefix string, accuracy float64, precision float64, recall float64, fvalue float64) error {
	apiKey := os.Getenv("MACKEREL_APIKEY")
	serviceName := os.Getenv("MACKEREL_SERVICE_NAME")
	if apiKey == "" || serviceName == "" {
		return nil
	}

	client := mkr.NewClient(apiKey)
	now := time.Now().Unix()
	err := client.PostServiceMetricValues(serviceName, []*mkr.MetricValue{
		{
			Name:  fmt.Sprintf("%s.accuracy", metricNamePrefix),
			Time:  now,
			Value: accuracy,
		},
		{
			Name:  fmt.Sprintf("%s.precision", metricNamePrefix),
			Time:  now,
			Value: precision,
		},
		{
			Name:  fmt.Sprintf("%s.recall", metricNamePrefix),
			Time:  now,
			Value: recall,
		},
		{
			Name:  fmt.Sprintf("%s.fvalue", metricNamePrefix),
			Time:  now,
			Value: fvalue,
		},
	})
	return err
}
