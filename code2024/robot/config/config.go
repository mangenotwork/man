package config

import (
	commonConfig "common/config"
	"github.com/jinzhu/configor"
	"sync"
)

var (
	cfg Config
	mu  sync.RWMutex
)

type (
	LoginServer struct {
		Scheme        string `json:"scheme"`
		Host          string `json:"host"`
		ConferenceId  int32  `json:"conferenceId"`
		SceneId       int32  `json:"sceneId"`
		SceneUniqueNO string `json:"sceneUniqueNO"`
		MinId         int32  `json:"minId"`
		MaxId         int32  `json:"maxId"`
	}

	Configuration struct {
		CommonPath  string      `json:"common_path"`
		LoginServer LoginServer `json:"login_server"`
	}

	Config struct {
		Configuration
		commonConfig.CommonConfig
	}
)

func Init(file string) (Config, error) {
	mu.Lock()
	defer mu.Unlock()

	var conf Configuration
	err := configor.Load(&conf, file)
	if err != nil {
		return Config{}, err
	}

	var commonCfg commonConfig.CommonConfig
	err = configor.Load(&commonCfg, conf.CommonPath)
	if err != nil {
		return Config{}, err
	}

	cfg = Config{
		Configuration: conf,
		CommonConfig:  commonCfg,
	}
	return cfg, err
}

func GetConfig() Config {
	mu.Lock()
	defer mu.Unlock()
	return cfg
}
