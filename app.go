package main

import (
	"context"
	"time"

	"wails-svelte/internal/state"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	maxSystrayDevices = 15
	eventDevicesList  = "devices:list"
	eventDevicesError = "devices:error"
)

type App struct {
	ctx            context.Context
	state          *state.State
	systrayDevices [maxSystrayDevices]*systray.MenuItem
	ticker         *time.Ticker
}

func NewApp(state *state.State) *App {
	a := &App{state: state}

	if state.IsConnected() {
		a.ticker = time.NewTicker(state.Cfg.ScanInterval * time.Second)
	}

	return a
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogWarning(a.ctx, "startup")

	go a.systray()
	a.updateDeviceList()

	if a.ticker != nil {
		go a.backgroundUpdate()
	}

}

func (a *App) systray() {
	runtime.LogWarning(a.ctx, "starting systray")
	systray.Run(a.onSystrayReady(), a.onSystrayExit())
	runtime.LogWarning(a.ctx, "starting systray done")
}

func (a *App) updateDeviceList() {
	a.eventDevicesList()
	a.onSystrayUpdateDevices()
}

func (a *App) eventDevicesList() {
	runtime.LogWarning(a.ctx, "emitting devices:list")

	devices, err := a.state.GetDevices()
	if err != nil {
		runtime.LogErrorf(a.ctx, "failed to get devices: %v", err)
		runtime.EventsEmit(a.ctx, eventDevicesError, err.Error())
		return
	}

	runtime.LogWarning(a.ctx, "emitting devices:list successfully")
	runtime.EventsEmit(a.ctx, eventDevicesList, devices)
}

func (a *App) backgroundUpdate() {
	if a.ticker == nil {
		runtime.LogWarning(a.ctx, "not running background update")
		return
	}

	runtime.LogWarning(a.ctx, "running background update")
	for range a.ticker.C {
		runtime.LogWarning(a.ctx, "running background update tick")
		a.updateDeviceList()
	}
}
