package metrics

import metricsLib "github.com/weeb-vip/go-metrics-lib"

// Metric result constants for easy usage
const (
	Success = metricsLib.Success
	Error   = metricsLib.Error
)

// Database method constants
const (
	MethodSelect = metricsLib.DatabaseMetricMethodSelect
	MethodInsert = metricsLib.DatabaseMetricMethodInsert
	MethodUpdate = metricsLib.DatabaseMetricMethodUpdate
	MethodDelete = metricsLib.DatabaseMetricMethodDelete
)

// Common table/component names
const (
	TableAnime        = "anime"
	TableAnimeSeason  = "anime_seasons"
	TableEpisodes     = "episodes"
	TableCharacters   = "characters"
	TableStaff        = "staff"

	ComponentResolver   = "resolver"
	ComponentService    = "service"
	ComponentRepository = "repository"
)