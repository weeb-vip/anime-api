package metrics

import (
	"fmt"
)

func ResolverCountMetricSuccess(ResolverName string) {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "success"),
	}

	client.GetClient().Incr("graphql.resolver.count", tags, 1)
}

func ResolverCountMetricError(ResolverName string) {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "error"),
	}

	client.GetClient().Incr("graphql.resolver.count", tags, 1)
}
