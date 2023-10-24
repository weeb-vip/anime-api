package metrics

import (
	"fmt"
	"github.com/DataDog/datadog-go/v5/statsd"
	"log"
	"os"
)

var (
	// MetricsPusherInstance is a singleton of MetricsPusher
	MetricsPusherInstance MetricsPusher
)

type MetricsPusher interface {
	GetClient() *statsd.Client
}

type MetricsPusherImpl struct {
	Client *statsd.Client
}

// NewMetricPusher singleton
func NewMetricPusher() MetricsPusher {
	// get DD_AGENT_HOST and DD_DOGSTATSD_PORT from env
	DD_AGENT_HOST := os.Getenv("DD_AGENT_HOST")
	dogstatsd_client, err := statsd.New(fmt.Sprintf("%s:%d", DD_AGENT_HOST, 8125))
	if err != nil {
		log.Fatal(err)
	}

	if dogstatsd_client == nil {
		log.Fatal("dogstatsd_client is nil")
	}

	if MetricsPusherInstance == nil {
		MetricsPusherInstance = &MetricsPusherImpl{
			Client: dogstatsd_client,
		}
	}

	return MetricsPusherInstance
}

func (m *MetricsPusherImpl) GetClient() *statsd.Client {
	return m.Client
}
