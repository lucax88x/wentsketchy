package config

import (
	"fmt"
)

func (cfg *Config) right(
	batches [][]string,
) ([][]string, error) {
	batches, err := cfg.items.Calendar.Init(batches, cfg.fifoPath)

	if err != nil {
		return batches, fmt.Errorf("init: calendar. %w", err)
	}

	batches, err = cfg.items.Battery.Init(batches, cfg.fifoPath)

	if err != nil {
		return batches, fmt.Errorf("init: battery. %w", err)
	}
	batches, err = cfg.items.FrontApp.Init(batches, cfg.fifoPath)

	if err != nil {
		return batches, fmt.Errorf("init: front-app. %w", err)
	}

	return batches, nil
}
