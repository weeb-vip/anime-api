package metrics

import (
	"fmt"
)

func ResolverDistributionMetricSuccess(ResolverName string, duration float64) error {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "success"),
	}

	fmt.Println(tags)

	err := client.GetClient().Distribution("graphql.resolver.dist.millisecond", duration, tags, 1)
	if err != nil {
		return err
	}

	return nil
}

func ResolverDistributionMetricError(ResolverName string, duration float64) error {
	client := NewMetricPusher()

	tags := []string{
		fmt.Sprintf("resolver:%s", ResolverName),
		fmt.Sprintf("result:%s", "error"),
	}

	fmt.Println(tags)

	err := client.GetClient().Distribution("graphql.resolver.dist.millisecond", duration, tags, 1)
	if err != nil {
		return err
	}

	return nil
}
