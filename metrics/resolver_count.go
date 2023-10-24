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

	client.GetClient().Count("graphql.resolver.count", 1, tags, 1)
}

func ResolverCountMetricError(ResolverName string) {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "error"),
	}

	client.GetClient().Count("graphql.resolver.count", 1, tags, 1)
}
