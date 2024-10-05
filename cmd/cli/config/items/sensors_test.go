//nolint:testpackage // want to test internals
package items

import (
	"context"
	"testing"

	"github.com/lucax88x/wentsketchy/internal/command"
	"github.com/lucax88x/wentsketchy/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestUnitSensors(t *testing.T) {
	ctx := context.Background()
	logger := testutils.CreateTestLogger()
	command := command.NewCommand(logger)
	item := NewSensorsItem(logger, command)

	t.Run("should get fan speeds", func(t *testing.T) {
		// WHEN
		result, err := item.getFanSpeeds(ctx)

		// THEN
		require.NoError(t, err)
		require.Len(t, result, 2)
		// fmt.Println(result[0])
		// fmt.Println(result[1])
	})

	t.Run("should get temperatures", func(t *testing.T) {
		// WHEN
		result, err := item.getTemperatures(ctx)

		// THEN
		require.NoError(t, err)
		require.NotZero(t, result.averageCPUs)
		require.NotZero(t, result.highest)
		// fmt.Println(result.highest)
		// fmt.Println(result.averageCPUs)
	})
}
