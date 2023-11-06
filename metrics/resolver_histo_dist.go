package metrics

func ResolverHistoDistMetricSuccess(ResolverName string, duration float64) error {
	err := ResolverDistributionMetricSuccess(ResolverName, duration)
	if err != nil {
		return err
	}

	err = ResolverHistrogramMetricSuccess(ResolverName, duration)
	if err != nil {
		return err
	}

	return nil
}

func ResolverHistoDistMetricError(ResolverName string, duration float64) error {
	err := ResolverDistributionMetricError(ResolverName, duration)
	if err != nil {
		return err
	}

	err = ResolverHistrogramMetricError(ResolverName, duration)
	if err != nil {
		return err
	}

	return nil
}
