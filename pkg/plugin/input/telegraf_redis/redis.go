package telegraf_redis

import (
	"encoding/json"
	"fmt"
	telegraf2 "github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs/redis"
	"github.com/traas-stack/holoinsight-agent/pkg/collecttask"
	"github.com/traas-stack/holoinsight-agent/pkg/logger"
	"github.com/traas-stack/holoinsight-agent/pkg/pipeline/telegraf/providers"
	"github.com/traas-stack/holoinsight-agent/pkg/plugin/api"
	"github.com/traas-stack/holoinsight-agent/pkg/telegraf"
)

type (
	Conf struct {
		Port     int    `json:"port,omitempty"`
		Password string `json:"password,omitempty"`
	}
)

func init() {
	providers.Register("telegraf_redis", func(task *collecttask.CollectTask) (interface{}, error) {
		return ParseInputPlugin(task)
	})
}

func ParseInputPlugin(task *collecttask.CollectTask) (api.Input, error) {
	conf := &Conf{}
	if err := json.Unmarshal(task.Config.Content, conf); err != nil {
		return nil, err
	}

	var plugin telegraf2.Input
	switch task.Target.Type {
	case collecttask.TargetPod:
		if conf.Port <= 0 {
			conf.Port = 6379
		}
		server := fmt.Sprintf("tcp://%s:%d", task.Target.GetIP(), conf.Port)
		plugin = &redis.Redis{
			Servers:  []string{server},
			Password: conf.Password,
			Log:      logger.ZapLogger.InfoS,
		}
	case collecttask.TargetNone:
		plugin = &redis.Redis{
			Log: logger.ZapLogger.InfoS,
		}
		if err := json.Unmarshal(task.Config.Content, plugin); err != nil {
			return nil, err
		}
	}
	return telegraf.NewInputAdapter(plugin), nil
}
