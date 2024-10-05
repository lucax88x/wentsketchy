package items

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/colors"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/icons"
	"github.com/lucax88x/wentsketchy/internal/command"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

type CPUItem struct {
	logger  *slog.Logger
	command *command.Command
}

func NewCPUItem(logger *slog.Logger, command *command.Command) CPUItem {
	return CPUItem{
		logger,
		command,
	}
}

const cpuBracketName = "cpu.bracket"
const cpuItemIconName = "cpu.icon"
const cpuItemTopName = "cpu.top"
const cpuItemPercentName = "cpu.percent"
const cpuItemSysName = "cpu.sys"
const cpuItemUserName = "cpu.user"
const cpuItemSpacerName = "cpu.spacer"

func (i CPUItem) Init(
	_ context.Context,
	position sketchybar.Position,
	batches Batches,
) (Batches, error) {
	updateEvent, err := args.BuildEvent()

	if err != nil {
		return batches, errors.New("cpu: could not generate update event")
	}

	cpuIconItem := sketchybar.ItemOptions{
		Display: "active",
		Icon: sketchybar.ItemIconOptions{
			Value: icons.CPU,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.Sketchybar.IconPadding,
				Right: pointer(*settings.Sketchybar.IconPadding / 2),
			},
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	cpuTopItem := sketchybar.ItemOptions{
		Display: "active",
		Label: sketchybar.ItemLabelOptions{
			Value: "",
			Font: sketchybar.FontOptions{
				Size: "8.0",
			},
		},
		Icon: sketchybar.ItemIconOptions{
			Drawing: "off",
		},
		Padding: sketchybar.PaddingOptions{
			Right: settings.Sketchybar.ItemSpacing,
		},
		YOffset: pointer(4),
		Width:   pointer(0),
	}
	cpuPercentItem := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  pointer(10),
			Right: settings.Sketchybar.ItemSpacing,
		},
		Label: sketchybar.ItemLabelOptions{
			Value: "",
			Font: sketchybar.FontOptions{
				Size: "8.0",
			},
		},
		Icon: sketchybar.ItemIconOptions{
			Drawing: "off",
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
		YOffset: pointer(-6),
		// Width:      pointer(0),
		UpdateFreq: pointer(4),
		Updates:    "on",
		Script:     updateEvent,
	}
	cpuSysItem := sketchybar.GraphOptions{
		Display: "active",
		Width:   pointer(75),
		Graph: sketchybar.ItemGraphOptions{
			Color:     colors.Red,
			FillColor: colors.Red,
		},
		Icon: sketchybar.ItemIconOptions{
			Drawing: "off",
		},
		Label: sketchybar.ItemLabelOptions{
			Drawing: "off",
		},
		YOffset: pointer(6),
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
			Height:  pointer(0),
		},
	}
	cpuUserItem := sketchybar.GraphOptions{
		Display: "active",
		Width:   pointer(0),
		Graph: sketchybar.ItemGraphOptions{
			Color: settings.Sketchybar.ItemBackgroundColor,
		},
		Icon: sketchybar.ItemIconOptions{
			Drawing: "off",
		},
		Label: sketchybar.ItemLabelOptions{
			Drawing: "off",
		},
		YOffset: pointer(10),
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
			Height:  pointer(0),
		},
	}
	cpuBracketItem := sketchybar.BracketOptions{
		Background: sketchybar.BackgroundOptions{
			Drawing: "on",
			Color: sketchybar.ColorOptions{
				Color: settings.Sketchybar.ItemBackgroundColor,
			},
		},
	}
	cpuSpacerItem := sketchybar.ItemOptions{
		Display: "active",
		Label: sketchybar.ItemLabelOptions{
			Value: "",
		},
		Padding: sketchybar.PaddingOptions{
			Right: settings.Sketchybar.ItemSpacing,
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	batches = batch(batches, s("--add", "item", cpuItemSpacerName, position))
	batches = batch(batches, m(s("--set", cpuItemSpacerName), cpuSpacerItem.ToArgs()))

	batches = batch(batches, s("--add", "item", cpuItemTopName, position))
	batches = batch(batches, m(s("--set", cpuItemTopName), cpuTopItem.ToArgs()))

	batches = batch(batches, s("--add", "item", cpuItemPercentName, position))
	batches = batch(batches, m(s("--set", cpuItemPercentName), cpuPercentItem.ToArgs()))

	batches = batch(batches, s("--add", "graph", cpuItemUserName, position, "75"))
	batches = batch(batches, m(s("--set", cpuItemUserName), cpuUserItem.ToArgs()))

	batches = batch(batches, s("--add", "graph", cpuItemSysName, position, "75"))
	batches = batch(batches, m(s("--set", cpuItemSysName), cpuSysItem.ToArgs()))

	batches = batch(batches, s("--add", "item", cpuItemIconName, position))
	batches = batch(batches, m(s("--set", cpuItemIconName), cpuIconItem.ToArgs()))

	batches = batch(batches, s(
		"--add",
		"bracket",
		cpuBracketName,
		cpuItemIconName,
		cpuItemTopName,
		cpuItemPercentName,
		cpuItemSysName,
		cpuItemUserName,
	))
	batches = batch(batches, m(s("--set", cpuBracketName), cpuBracketItem.ToArgs()))

	return batches, nil
}

func (i CPUItem) Update(
	ctx context.Context,
	batches Batches,
	_ sketchybar.Position,
	args *args.In,
) (Batches, error) {
	if !isCPU(args.Name) {
		return batches, nil
	}

	if args.Event == events.Routine || args.Event == events.Forced {
		topProcess, err := i.getTopProcess(ctx)

		if err != nil {
			return batches, err
		}

		cpuLoad, err := i.getCPULoad()

		if err != nil {
			return batches, err
		}

		cpuTopItem := sketchybar.ItemOptions{
			Label: sketchybar.ItemLabelOptions{
				Value: fmt.Sprintf("%.2f%% %s", topProcess.cpu, truncateString(topProcess.name, 8)),
			},
		}
		cpuPercentItem := sketchybar.ItemOptions{
			Label: sketchybar.ItemLabelOptions{
				Value: fmt.Sprintf("%.2f%%", cpuLoad.sys+cpuLoad.user),
			},
		}

		batches = batch(batches, s("--push", cpuItemSysName, fmt.Sprintf("%.2f", cpuLoad.sys/100)))
		batches = batch(batches, s("--push", cpuItemUserName, fmt.Sprintf("%.2f", cpuLoad.user/100)))

		batches = batch(batches, m(s("--set", cpuItemPercentName), cpuPercentItem.ToArgs()))
		batches = batch(batches, m(s("--set", cpuItemTopName), cpuTopItem.ToArgs()))
	}

	return batches, nil
}

func isCPU(name string) bool {
	return name == cpuItemPercentName
}

type process struct {
	pid  int
	cpu  float64
	name string
}

func (item CPUItem) getProcesses(ctx context.Context) ([]*process, error) {
	out, err := item.command.RunBufferized(ctx, "ps", "auxcr")

	if err != nil {
		return make([]*process, 0), errors.New("cpu: could not get processes")
	}

	processes := make([]*process, 0)
	for {
		line, err := out.ReadString('\n')

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return make([]*process, 0), fmt.Errorf("cpu: could not read line. %w", err)
		}

		if strings.HasPrefix(line, "USER") {
			continue
		}

		tokens := strings.Split(line, " ")
		ft := make([]string, 0)

		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}

		pid, err := strconv.Atoi(ft[1])

		if err != nil {
			return make([]*process, 0), fmt.Errorf("cpu: could not parse pid from %s", line)
		}

		cpu, err := strconv.ParseFloat(ft[2], 64)

		if err != nil {
			return make([]*process, 0), fmt.Errorf("cpu: could not parse cpu, %s. %w", line, err)
		}

		name := ft[10]

		if cpu > 0 {
			processes = append(processes, &process{pid, cpu, name})
		}
	}

	return processes, nil
}

func (i CPUItem) getTopProcess(ctx context.Context) (*process, error) {
	processes, err := i.getProcesses(ctx)

	if err != nil {
		return nil, fmt.Errorf("cpu: could get processes. %w", err)
	}

	if len(processes) == 0 {
		return nil, fmt.Errorf("cpu: no processes found. %w", err)
	}

	return processes[0], nil
}

type cpuLoad struct {
	user float32
	sys  float32
	idle float32
}

func (item CPUItem) getCPULoad() (*cpuLoad, error) {
	cmd := "top | head -n 4"
	output, err := exec.Command("zsh", "-c", cmd).Output()
	if err != nil {
		return &cpuLoad{}, fmt.Errorf("cpu: could not get top")
	}

	cpuRegexp := regexp.MustCompile(`(\d+\.\d+)% user, (\d+\.\d+)% sys, (\d+\.\d+)% idle`)

	matches := cpuRegexp.FindStringSubmatch(string(output))

	if len(matches) != 4 {
		return nil, errors.New("cpu: could not get cpu load, not enough matches")
	}

	user, err1 := strconv.ParseFloat(matches[1], 32)
	sys, err2 := strconv.ParseFloat(matches[2], 32)
	idle, err3 := strconv.ParseFloat(matches[3], 32)

	// Check for parsing errors
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("cpu: error parsing CPU usage values")
	}

	// Convert to float32
	userCPU := float32(user)
	sysCPU := float32(sys)
	idleCPU := float32(idle)

	return &cpuLoad{
		userCPU,
		sysCPU,
		idleCPU,
	}, nil
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return fmt.Sprintf("%s...", s[:maxLen])
	}
	return s
}

var _ WentsketchyItem = (*CPUItem)(nil)
