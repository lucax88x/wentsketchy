//nolint:lll // tests
package args_test

import (
	"testing"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/stretchr/testify/require"
)

func TestUnitArgs(t *testing.T) {
	t.Run("should build event correclty", func(t *testing.T) {
		// WHEN
		event, err := args.BuildEvent()

		// THEN
		require.NoError(t, err)
		require.Equal(t, `echo "update args: {\"name\":\"$NAME\",\"event\":\"$SENDER\",\"button\":\"$BUTTON\",\"modifier\":\"$MODIFIER\"} info: $INFO ¬" > /path`, event)
	})

	t.Run("should extract args from event", func(t *testing.T) {
		// GIVEN
		event := `update args: {"name":"some-name","event":"some-sender","button":"some-button","modifier":"some-modifier"} info: { "key": "value" } `

		// WHEN
		argsIn, err := args.FromEvent(event)

		// THEN
		require.NoError(t, err)
		require.Equal(t, "some-name", argsIn.Name)
		require.Equal(t, "some-sender", argsIn.Event)
		require.Equal(t, `{ "key": "value" } `, argsIn.Info)
		require.Equal(t, "some-button", argsIn.Button)
		require.Equal(t, "some-modifier", argsIn.Modifier)
	})

	t.Run("should extract args from event where info is multiline json", func(t *testing.T) {
		// GIVEN
		event := `update args: {"name":"aerospace-checker","event":"space_change","button":"","modifier":""} info: {
	"display-1": 1
}`

		// WHEN
		argsIn, err := args.FromEvent(event)

		// THEN
		require.NoError(t, err)
		require.Equal(t, "aerospace-checker", argsIn.Name)
		require.Equal(t, "space_change", argsIn.Event)
		require.Equal(t, `{
	"display-1": 1
}`, argsIn.Info)
	})
}
