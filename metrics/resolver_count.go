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

	fmt.Println(tags)

	err := client.GetClient().Count("graphql.resolver.count", 1, tags, 1)
	if err != nil {
		fmt.Println(err)
	}
}

func ResolverCountMetricError(ResolverName string) {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "error"),
	}

	fmt.Println(tags)

	err := client.GetClient().Count("graphql.resolver.count", 1, tags, 1)
	if err != nil {
		fmt.Println(err)
	}
}
