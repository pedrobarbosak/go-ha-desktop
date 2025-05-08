package main

import (
	"time"

	"wails-svelte/internal/config"
	"wails-svelte/internal/ha"
)

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

	a.updateDeviceList()

	return device, nil
}

func (a *App) TurnOff(entityID string) (*ha.Device, error) {
	device, err := a.state.TurnOff(entityID)
	if err != nil {
		return nil, err
	}

	a.updateDeviceList()

	return device, nil
}

func (a *App) UpdateConfig(updated *config.Config) error {
	if err := a.state.UpdateConfig(updated); err != nil {
		return err
	}

	if a.ticker != nil {
		a.ticker.Stop()
		a.ticker = nil
	}

	if updated.ScanInterval > 0 {
		a.ticker = time.NewTicker(updated.ScanInterval)
		go a.backgroundUpdate()
	}

	a.updateDeviceList()

	return nil
}

func (a *App) TestConnection(url string, accessToken string) error {
	return nil
}
