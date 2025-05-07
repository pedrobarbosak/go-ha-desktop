package main

import (
	"fmt"
	"log"

	"wails-svelte/internal/ha"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) onSystrayReady() func() {
	return func() {
		systray.SetIcon(icon)
		systray.SetTitle("ha-desktop")
		systray.SetTooltip("HA Desktop")

		a.systrayDevices = [maxSystrayDevices]*systray.MenuItem{}
		for i := 0; i < maxSystrayDevices; i++ {
			s := fmt.Sprintf("Device %d", i+1)
			mDevice := systray.AddMenuItemCheckbox(s, s, false)
			mDevice.Hide()
			a.systrayDevices[i] = mDevice
		}

		systray.AddSeparator()

		mShow := systray.AddMenuItem("Show", "Show the app")
		mShow.SetIcon(icon)
		mShow.Click(func() {
			if runtime.WindowIsMinimised(a.ctx) {
				runtime.WindowUnminimise(a.ctx)
				return
			}

			runtime.Show(a.ctx)
		})

		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		mQuit.SetIcon(icon)
		mQuit.Click(func() {
			a.onSystrayExit()()
		})

		a.onSystrayUpdateDevices()
	}
}

func (a *App) onSystrayExit() func() {
	return func() {
		log.Println("--- Exiting...")
		systray.Quit()
		runtime.Quit(a.ctx)
	}
}

func (a *App) onSystrayUpdateDevices() {
	for i := 0; i < maxSystrayDevices; i++ {
		if a.systrayDevices[i] == nil {
			continue
		}

		a.systrayDevices[i].Hide()
		a.systrayDevices[i].Uncheck()
	}

	devices := a.state.GetPinnedDevices()
	for i, device := range devices {
		if device.Error != nil {
			continue
		}

		a.systrayDevices[i].SetTitle(device.Name)
		if device.State {
			a.systrayDevices[i].Check()
		}

		a.systrayDevices[i].Click(func() {
			var err error
			var updatedDevice *ha.Device

			if a.systrayDevices[i].Checked() {
				updatedDevice, err = a.TurnOff(device.ID)
			} else {
				updatedDevice, err = a.TurnOn(device.ID)
			}

			if err != nil {
				log.Println("Failed to toggle device:", err)
				return
			}

			if updatedDevice.State {
				a.systrayDevices[i].Check()
			} else {
				a.systrayDevices[i].Uncheck()
			}
		})

		a.systrayDevices[i].Show()
	}

	return
}
