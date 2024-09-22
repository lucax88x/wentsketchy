package config

import (
	"fmt"
)

func (cfg *Config) left(
	batches [][]string,
) ([][]string, error) {
	batches, err := cfg.items.MainIcon.Init(batches)

	if err != nil {
		return batches, fmt.Errorf("init aerospace %w", err)
	}

	batches, err = cfg.items.Aerospace.Init(batches, cfg.fifoPath)

	if err != nil {
		return batches, fmt.Errorf("init aerospace %w", err)
	}

	return batches, nil
}
