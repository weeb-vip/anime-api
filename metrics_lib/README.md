# Go Metrics Library

The purpose of this library is to make it easier to write metrics and provide a standard for usage of metrics. This
library currently supports these standard metrics:

## Standard Metrics

| Metric                                           | Labels                                                                                                                                               | Description                                                                                                                                 | |
|:-------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------|:-|
| resolver_request_duration_histogram_milliseconds | resolver= function name of the resolver<br/>result= success \|fail<br/>service= name of the service<br/>protocol= http \|grpc \|graphql              | This metric gives an overview of success/failures of resolvers, the duration of resolvers, and the distribution of the duration of requests | |
| http_request_duration_histogram_milliseconds     | result= success \|fail<br/>service= name of the service<br/>method= POST\|GET\|PATCH…                                                                | all http requests to our service (datadog gives to us for free).                                                                            | |
| api_request_duration_histogram_milliseconds      | service= current service<br/>vendor= internal or external vendor<br/>call= name of the query being called (function name)<br/>result=success \| fail | Calculating communication between services or vendors, where they came from, where they are meant to go, duration of request.               | |
| database_query_duration_histogram_milliseconds   | service= service name<br/>result= success \|fail<br/>table= table name<br/>method= insert \|delete \|find<br/>database= mongodb \|postgres           | Getting duration of queries in respect to the service they are in.                                                                          | |
| call_duration_histogram_milliseconds             | service= service name<br/>result= success \|fail<br/>function= function name                                                                         | Looking at the duration of a call for a function (not for every function, used for things we want to watch)                                 | |


See examples in examples folder.

Example Usage:

```go
package main

import (
	MetricsLib "github.com/tempmee/go-metrics-lib"
	"github.com/tempmee/go-metrics-lib/clients/datadog"
	"log"
)

type Result string

const (
	ResultSuccess Result = "success"
	ResultError   Result = "error"
)

type Labels struct {
	Name    string
	Service string
	Result  Result
}

func main() {
	datadogClient := datadog.NewDatadogClient(datadog.DataDogConfig{
		DD_AGENT_HOST: "localhost",
		DD_AGENT_PORT: 8125,
	})
	metrics := MetricsLib.NewMetrics(
		datadogClient,
		1,
	)
	err := metrics.HistogramMetric("graphql.resolver.millisecond", 100,
		map[string]string{
			"resolver": "resolver",
			"service":  "graphql",
			"result":   "success",
		},
	)

	if err != nil {
		log.Println("BORKED!")
		panic(err)
	}

	err = metrics.SummaryMetric("graphql.resolver.millisecond", 100, map[string]string{
		"resolver": "resolver",
		"service":  "graphql",
		"result":   "success",
	})

	if err != nil {
		log.Println("BORKED!")

	}

	err = metrics.ResolverMetric(100, MetricsLib.ResolverMetricLabels{
		Resolver: "resolver",
		Result:   MetricsLib.Success,
	})

	if err != nil {
		log.Println("BORKED!")
		panic(err)
	}

}

```