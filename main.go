package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wails-svelte/internal/state"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
	const name = "ha-desktop"
	const width = 1024
	const height = 768

	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC:", r)
		}
	}()

	state := state.New()
	app := NewApp(state)

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGSEGV, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL)

		<-signalChan

		log.Println("--- Exiting...")
	}()

	err := wails.Run(&options.App{
		Title:             name,
		Width:             width,
		Height:            height,
		MinWidth:          width,
		MinHeight:         height,
		StartHidden:       true,
		HideWindowOnClose: true,
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 255},
		AssetServer:       &assetserver.Options{Assets: assets},
		LogLevel:          logger.WARNING, // use warning because default info logs a lot of stuff
		OnStartup:         app.startup,
		OnDomReady:        func(ctx context.Context) { log.Println("--- Dom is ready") },
		Bind:              []interface{}{app},
		Linux:             &linux.Options{Icon: icon, ProgramName: name},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	})

	if err != nil {
		log.Panicln("E:", err)
	}
}
