//nolint:testpackage // want to test internals
package items

import (
	"context"
	"testing"

	"github.com/lucax88x/wentsketchy/internal/command"
	"github.com/lucax88x/wentsketchy/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestUnitCpu(t *testing.T) {
	ctx := context.Background()
	logger := testutils.CreateTestLogger()
	command := command.NewCommand(logger)
	item := NewCPUItem(logger, command)

	t.Run("should max process load", func(t *testing.T) {
		// WHEN
		process, err := item.getTopProcess(ctx)

		// THEN
		require.NoError(t, err)
		require.NotNil(t, process)
		require.NotNil(t, process.name)
		require.NotNil(t, process.cpu)
		require.NotNil(t, process.pid)
	})

	t.Run("should get cpu load", func(t *testing.T) {
		// WHEN
		cpuLoad, err := item.getCPULoad()

		// THEN
		require.NoError(t, err)
		require.Greater(t, cpuLoad.sys, float32(0))
		require.Greater(t, cpuLoad.user, float32(0))
		require.Greater(t, cpuLoad.idle, float32(0))
	})
}
