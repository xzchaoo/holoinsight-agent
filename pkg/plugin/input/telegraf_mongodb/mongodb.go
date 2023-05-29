package telegraf_mongodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/influxdata/telegraf/plugins/inputs/mongodb"
	"github.com/traas-stack/holoinsight-agent/pkg/collecttask"
	"github.com/traas-stack/holoinsight-agent/pkg/logger"
	"github.com/traas-stack/holoinsight-agent/pkg/pipeline/telegraf/providers"
	"github.com/traas-stack/holoinsight-agent/pkg/telegraf"
)

type (
	Conf struct {
		Port int `json:"port,omitempty"`
	}
)

func init() {
	providers.Register("telegraf_mongodb", func(task *collecttask.CollectTask) (interface{}, error) {
		conf := &Conf{}
		if err := json.Unmarshal(task.Config.Content, conf); err != nil {
			return nil, err
		}

		var telegrafInput *mongodb.MongoDB
		switch task.Target.Type {
		case collecttask.TargetPod:
			ip := task.Target.GetIP()
			if conf.Port <= 0 {
				conf.Port = 27017
			}
			telegrafInput = &mongodb.MongoDB{
				Servers: []string{fmt.Sprintf("mongodb://%s:%d/?connect=direct", ip, conf.Port)},
				Log:     logger.ZapLogger.InfoS,
			}
		case collecttask.TargetNone:
			telegrafInput = &mongodb.MongoDB{
				Log: logger.ZapLogger.InfoS,
			}
			if err := json.Unmarshal(task.Config.Content, telegrafInput); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unsupported target type: " + task.Target.Type)
		}
		if err := telegrafInput.Init(); err != nil {
			return nil, err
		}
		return telegraf.NewInputAdapter(telegrafInput), nil
	})
}
