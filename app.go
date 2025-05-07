package main

import (
	"context"
	"fmt"
	"log"

	"wails-svelte/internal/config"
	"wails-svelte/internal/ha"
	"wails-svelte/internal/state"

	"github.com/energye/systray"
)

const maxSystrayDevices = 15

type App struct {
	ctx            context.Context
	state          *state.State
	systrayDevices [maxSystrayDevices]*systray.MenuItem
}

func NewApp(state *state.State) *App {
	return &App{state: state}
}

func (a *App) startup(ctx context.Context) {
	log.Println("startup")

	a.ctx = ctx
	a.systray()
}

func (a *App) GetError() error {
	return a.state.Err
}

func (a *App) GetDevices() ([]*ha.Device, error) {
	return a.state.GetDevices()
}

func (a *App) GetConfig() (*config.Config, error) {
	return a.state.GetConfig(), nil
}

func (a *App) TurnOn(entityID string) (*ha.Device, error) {
	device, err := a.state.TurnOn(entityID)
	if err != nil {
		return nil, err
	}

	go a.onSystrayUpdateDevices()

	return device, nil
}

func (a *App) TurnOff(entityID string) (*ha.Device, error) {
	device, err := a.state.TurnOff(entityID)
	if err != nil {
		return nil, err
	}

	go a.onSystrayUpdateDevices()

	return device, nil
}

func (a *App) UpdateConfig(updated *config.Config) error {
	if err := a.state.UpdateConfig(updated); err != nil {
		return err
	}

	go a.onSystrayUpdateDevices()

	return nil
}

func (a *App) TestConnection(url string, accessToken string) error {
	return nil
}

func (a *App) GetHA() ha.Client {
	return a.state.Ha
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) systray() {
	systray.Run(a.onSystrayReady(), a.onSystrayExit())
}
