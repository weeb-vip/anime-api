package datadog

import (
	"fmt"
	"github.com/DataDog/datadog-go/v5/statsd"
	metrics_lib "github.com/tempmee/go-metrics-lib"
	"log"
)

type DataDogClient struct {
	Client *statsd.Client
}

type DataDogConfig struct {
	DD_AGENT_HOST string `env:"DD_AGENT_HOST" default:"localhost"`
	DD_AGENT_PORT int    `env:"DD_AGENT_PORT" default:"8125"`
}

func NewDatadogClient(datadogConfig DataDogConfig) metrics_lib.Client {
	dogstatsd_client, err := statsd.New(fmt.Sprintf("%s:%d", datadogConfig.DD_AGENT_HOST, datadogConfig.DD_AGENT_PORT))
	if err != nil {
		log.Fatal(err)
	}

	if dogstatsd_client == nil {
		log.Fatal("dogstatsd_client is nil")
	}
	return &DataDogClient{
		dogstatsd_client,
	}
}

func (d *DataDogClient) Histogram(metric string, value float64, labels map[string]string, rate float64) error {
	tags := labelsToStringArray(labels)
	err := d.Client.Histogram(metric, value, tags, rate)
	if err != nil {
		return err
	}

	err = d.Client.Distribution(metric, value, tags, rate)
	if err != nil {
		return err
	}

	return nil
}

func (d *DataDogClient) Count(metric string, labels map[string]string, rate float64) error {
	tags := labelsToStringArray(labels)
	err := d.Client.Count(metric, 1, tags, rate)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataDogClient) Gauge(metric string, value float64, labels map[string]string, rate float64) error {
	tags := labelsToStringArray(labels)
	err := d.Client.Gauge(metric, value, tags, rate)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataDogClient) Summary(metric string, value float64, labels map[string]string, rate float64) error {
	log.Println("[Datadog] Summary is unsupported")
	return nil
}

func labelsToStringArray(labels map[string]string) []string {
	var tags []string
	for k, v := range labels {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}
	return tags
}
