package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func initLeft(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	err := initAerospace(ctx, di)

	if err != nil {
		return fmt.Errorf("init aerospace %w", err)
	}

	return nil
}

func initAerospace(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	monitors, err := di.Aerospace.Monitors(ctx)

	if err != nil {
		return err
	}

	var aggregatedErr error
	for _, monitor := range monitors {
		di.Logger.InfoContext(ctx, monitor.Name)

		workspaces, err := di.Aerospace.Workspaces(ctx, monitor.Id)

		if err != nil {
			aggregatedErr = errors.Join(aggregatedErr, err)
		}

		for _, workspace := range workspaces {
			space := sketchybar.ItemOptions{
				Icon: sketchybar.ItemIconOptions{
					Value: workspace.Id,
					Font: sketchybar.FontOptions{
						Font: font,
						Kind: "Regular",
						Size: "12.0",
					},
					PaddingOptions: sketchybar.PaddingOptions{
						Right: 5,
						Left:  10,
					},
				},
				Label: sketchybar.ItemLabelOptions{
					PaddingOptions: sketchybar.PaddingOptions{
						Right: 5,
						Left:  10,
					},
				},
				Background: sketchybar.BackgroundOptions{
					BorderOptions: sketchybar.BorderOptions{
						Width: 2,
						Color: ColorBackground1,
					},
					ColorOptions: sketchybar.ColorOptions{
						Color: ColorBackground1,
					},
				},
				// UpdateFreq: 30,
				// Script:     "wentsketchy update calendar",
				// Click_script:            "$PLUGIN_DIR/zen.sh",
			}

			// spaceID := fmt.Sprintf("space.%s", workspace.Id)
			spaceID := "space"
			di.Sketchybar.Run(
				ctx,
				s("--add", "item", "space", "left"),
				m(s("--set", spaceID), space.ToArgs()),
			)
		}
		// space=(
		//   space="$workspace_id"
		//   icon="${workspace_icons[$workspace_id]}"
		//   icon.font="$ICON_FONT:Regular:14.0"
		//   icon.highlight_color="$RED"
		//   icon.padding_left=10
		//   icon.padding_right=10
		//   display="$m"
		//   padding_left=1
		//   padding_right=1
		//   label.padding_right=20
		//   label.color="$GREY"
		//   label.highlight_color="$WHITE"
		//   label.font="sketchybar-app-font:Regular:14.0"
		//   label.y_offset=-1
		//   background.drawing="on"
		//   background.color="$BACKGROUND_2"
		//   background.border_color="$BACKGROUND_1"
		//   script="$PLUGIN_DIR/space.sh"
		// )
		//
		// apps=$(aerospace list-windows --workspace "$workspace_id" | awk -F'|' '{gsub(/^ *| *$/, "", $2); print $2}')
		//
		// icon_strip=" "
		// if [ "${apps}" != "" ]; then
		//   while read -r app; do
		//     icon_strip+=" $("$CONFIG_DIR/plugins/icon_map.sh" "$app")"
		//   done <<<"${apps}"
		// else
		//   icon_strip=" —"
		// fi
		//
		// sketchybar --add space "space.$workspace_id" left \
		//   --set "space.$workspace_id" "${space[@]}" label="$icon_strip" \
		//   --subscribe "space.$workspace_id" mouse.clicked
	}

	return nil
}
