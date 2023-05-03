package resolvers

import (
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/graph/model"
)

func APIInfo(cfg config.Config) (*model.APIInfo, error) {
	return &model.APIInfo{
		Name:    cfg.AppConfig.APPName,
		Version: cfg.AppConfig.Version,
	}, nil
}
