package metrics

import (
	"fmt"
)

func ResolverHistrogramMetricSuccess(ResolverName string, duration float64) {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "success"),
	}

	fmt.Println(tags)

	err := client.GetClient().Histogram("graphql.resolver.millisecond", duration, tags, 1)
	if err != nil {
		fmt.Println(err)
	}
}

func ResolverHistrogramMetricError(ResolverName string, duration float64) {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "error"),
	}

	fmt.Println(tags)

	err := client.GetClient().Histogram("graphql.resolver.millisecond", duration, tags, 1)
	if err != nil {
		fmt.Println(err)
	}
}
