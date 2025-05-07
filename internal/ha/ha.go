package ha

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/pedrobarbosak/go-errors"
	"github.com/ryanjohnsontv/go-homeassistant/logging"
	"github.com/ryanjohnsontv/go-homeassistant/rest"
	"github.com/ryanjohnsontv/go-homeassistant/shared/constants/domains"
	"github.com/ryanjohnsontv/go-homeassistant/shared/types"
)

type Client interface {
	GetLights() ([]*Device, error)
	TurnOn(entityID string) ([]*Device, error)
	TurnOff(entityID string) ([]*Device, error)
}

type service struct {
	cl          *rest.Client
	URL         string
	AccessToken string
}

type Device struct {
	ID    string
	Name  string
	Type  string
	State bool
	Error error
}

func NewDevice(entity types.Entity) *Device {
	d := &Device{ID: entity.EntityID.String(), Name: entity.EntityID.Name(), Type: string(entity.EntityID.Domain())}

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

func (s *service) TurnOn(entityID string) ([]*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data := map[string]interface{}{
		"entity_id": entityID,
	}

	updated, err := s.cl.CallService(ctx, domains.Light, "turn_on", data)
	if err != nil {
		return nil, errors.New("failed to turn on light:", err)
	}

	devices := make([]*Device, 0, len(updated))
	for _, ent := range updated {
		devices = append(devices, NewDevice(ent))
	}

	return devices, nil
}

func (s *service) TurnOff(entityID string) ([]*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data := map[string]interface{}{
		"entity_id": entityID,
	}

	updated, err := s.cl.CallService(ctx, domains.Light, "turn_off", data)
	if err != nil {
		return nil, errors.New("failed to turn off light:", err)
	}

	devices := make([]*Device, 0, len(updated))
	for _, ent := range updated {
		devices = append(devices, NewDevice(ent))
	}

	return devices, nil
}

func (s *service) GetLights() ([]*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	states, err := s.cl.GetStates(ctx)
	if err != nil {
		return nil, errors.New("failed to get states:", err)
	}

	lights := make([]*Device, 0, len(states))
	for _, ent := range states {
		if ent.EntityID.Domain() != domains.Light {
			continue
		}

		lights = append(lights, NewDevice(ent))
	}

	slices.SortFunc(lights, func(a, b *Device) int {
		return strings.Compare(a.ID, b.ID)
	})

	return lights, nil
}

func New(url, token string) (Client, error) {
	logger := &logging.DefaultLogger{}
	logger.SetLevel(logging.WarnLevel)

	client, err := rest.NewClient(url, token)
	if err != nil {
		return nil, errors.New("failed to create client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = client.GetHealth(ctx); err != nil {
		return nil, errors.New("failed to get health:", err)
	}

	return &service{
		cl:          client,
		URL:         url,
		AccessToken: token,
	}, nil
}
