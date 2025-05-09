package ha

import (
	"encoding/json"
	"slices"

	"github.com/pedrobarbosak/go-errors"
	"github.com/ryanjohnsontv/go-homeassistant/shared/constants/domains"
	"github.com/ryanjohnsontv/go-homeassistant/shared/types"
)

type Device struct {
	ID         string
	Name       string
	Type       string
	State      bool
	Error      error
	Brightness *uint
}

func NewDevice(entity types.Entity) *Device {
	d := &Device{ID: entity.EntityID.String(), Name: entity.EntityID.Name(), Type: string(entity.EntityID.Domain())}

	parseDeviceBrightness(d, entity)

	if entity.State.IsUnavailable() {
		d.Error = errors.New("device is unavailable")
		return d
	}

	if entity.State.IsUnknown() {
		d.Error = errors.New("device is unknown")
		return d
	}

	status, err := entity.State.AsBool()
	if err != nil {
		d.Error = err
		return d
	}

	d.State = *status

	return d
}

func parseDeviceBrightness(d *Device, entity types.Entity) {
	if entity.EntityID.Domain() != domains.Light {
		return
	}

	if entity.Attributes == nil {
		return
	}

	var attributes map[string]interface{}
	if err := json.Unmarshal(entity.Attributes, &attributes); err != nil {
		d.Error = err
		return
	}

	if v, ok := attributes["brightness"]; ok {
		if v != nil {
			if brightness, ok := v.(float64); ok {
				b := uint(brightness)
				d.Brightness = &b
				return
			}
		}
	}

	supportedColorModes, ok := attributes["supported_color_modes"].([]string)
	if !ok {
		return
	}

	if slices.Contains(supportedColorModes, "brightness") {
		var brightness uint
		d.Brightness = &brightness
	}
}
