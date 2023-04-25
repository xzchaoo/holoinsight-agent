/*
 * Copyright 2022 Holoinsight Project Authors. Licensed under Apache-2.0.
 */

package executor

import "github.com/traas-stack/holoinsight-agent/pkg/util"

type (
	processTimeParser struct {
	}
)

func (a *processTimeParser) Parse(context *LogContext) (int64, error) {
	return util.CurrentMS(), nil
}
