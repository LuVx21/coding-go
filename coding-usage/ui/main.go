package main

import (
	_ "github.com/gogpu/gg/gpu"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/ui/app"
	"github.com/gogpu/ui/primitives"
	"github.com/gogpu/ui/widget"
)

func main() {
	gogpuApp := gogpu.NewApp(
		gogpu.DefaultConfig().
			WithTitle("gogpu/ui — Widget Demo").
			WithSize(800, 900).
			WithContinuousRender(false), // Event-driven: 0% CPU when idle
	)

	uiApp := app.New(
		app.WithWindowProvider(gogpuApp),
		app.WithPlatformProvider(gogpuApp),
		app.WithEventSource(gogpuApp.EventSource()),
	)
	uiApp.SetRoot(buildUI())
}

func buildUI() *primitives.BoxWidget {
	card := primitives.Box().
		Padding(32).
		Gap(12).
		Background(widget.RGBA8(255, 255, 255, 255)).
		Rounded(12).
		ShadowLevel(2)

	return primitives.Box(card).Padding(24)
}
