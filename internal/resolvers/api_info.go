package resolvers

import (
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/graph/model"
)

func APIInfo(cfg config.Config) (*model.APIInfo, error) {
	return &model.APIInfo{
		Name: cfg.AppConfig.APPName,
		AnimeAPI: &model.AnimeAPI{
			Version: cfg.AppConfig.Version,
		},
	}, nil
}

func AnimeAPI(cfg config.Config) (*model.AnimeAPI, error) {
	return &model.AnimeAPI{
		Version: cfg.AppConfig.Version,
	}, nil
}
