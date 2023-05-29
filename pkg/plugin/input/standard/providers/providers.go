/*
 * Copyright 2022 Holoinsight Project Authors. Licensed under Apache-2.0.
 */

package providers

import (
	"fmt"
	"github.com/traas-stack/holoinsight-agent/pkg/collecttask"
	"github.com/traas-stack/holoinsight-agent/pkg/plugin/api"
	"sync"
)

type (
	InputProvider func(task *collecttask.CollectTask) (api.Input, error)
)

var (
	providers = make(map[string]InputProvider)
	mutex     sync.RWMutex
)

func Register(configType string, p InputProvider) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := providers[configType]; ok {
		panic(fmt.Errorf("duplicated input provider: %s", configType))
	}

	providers[configType] = p
}

func Get(configType string) (InputProvider, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	p, ok := providers[configType]
	return p, ok
}