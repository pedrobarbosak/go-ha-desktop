package state

import (
	"log"
	"sync"

	"wails-svelte/internal/config"
	"wails-svelte/internal/ha"

	"github.com/pedrobarbosak/go-errors"
)

type State struct {
	m sync.RWMutex

	Err     error
	Ha      ha.Client
	Cfg     *config.Config
	Devices []*ha.Device
}

func (s *State) UpdateConfig(updated *config.Config) error {
	s.m.Lock()
	defer s.m.Unlock()

	if err := config.Save(updated); err != nil {
		return err
	}

	s.Cfg = updated
	s.connectClient()

	s.Devices, _ = s.Ha.GetDevices()

	return nil
}

func (s *State) TurnOn(entityID string) (*ha.Device, error) {
	s.m.Lock()
	defer s.m.Unlock()

	updated, err := s.Ha.TurnOn(entityID)
	if err != nil {
		return nil, err
	}

	if len(updated) == 0 {
		_, s.Err = s.getDevices()
		if s.Err != nil {
			return nil, nil
		}

		updated = append(updated, &ha.Device{ID: entityID, State: true})
	}

	for _, u := range updated {
		for _, d := range s.Devices {
			if d.ID == u.ID {
				d.State = u.State
				return d, nil
			}
		}
	}

	return nil, errors.New("device not found")
}

func (s *State) TurnOff(entityID string) (*ha.Device, error) {
	s.m.Lock()
	defer s.m.Unlock()

	updated, err := s.Ha.TurnOff(entityID)
	if err != nil {
		return nil, err
	}

	if len(updated) == 0 {
		_, s.Err = s.getDevices()
		if s.Err != nil {
			return nil, nil
		}

		updated = append(updated, &ha.Device{ID: entityID, State: false})
	}

	for _, u := range updated {
		for _, d := range s.Devices {
			if d.ID == u.ID {
				d.State = u.State
				return d, nil
			}
		}
	}

	return nil, errors.New("device not found")
}

func (s *State) GetDevices() ([]*ha.Device, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.getDevices()
}

func (s *State) getDevices() ([]*ha.Device, error) {
	s.Devices, s.Err = s.Ha.GetDevices()
	if s.Err != nil {
		return nil, s.Err
	}

	return s.Devices, nil
}

func (s *State) GetConfig() *config.Config {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.Cfg
}

func (s *State) GetPinnedDevices() []*ha.Device {
	s.m.RLock()
	defer s.m.RUnlock()

	if s.Cfg == nil {
		return []*ha.Device{}
	}

	if len(s.Cfg.PinnedDevices) == 0 {
		return []*ha.Device{}
	}

	pinned := make([]*ha.Device, 0, len(s.Cfg.PinnedDevices))
	for _, id := range s.Cfg.PinnedDevices {
		for _, d := range s.Devices {
			if d.ID == id {
				pinned = append(pinned, d)
				break
			}
		}
	}

	return pinned
}

func (s *State) connectClient() {
	s.Ha = nil

	for _, url := range s.Cfg.URls {
		s.Ha, s.Err = ha.New(url, s.Cfg.AccessToken)
		if s.Err != nil {
			log.Println("Invalid URL:", url, s.Err)
			continue
		}

		s.Devices, s.Err = s.Ha.GetDevices()
		if s.Err != nil {
			return
		}

		log.Println("Connected to:", url)
		return
	}
}

func (s *State) IsConnected() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.Ha != nil
}

func (s *State) SetBrightness(entityID string, brightness uint) (*ha.Device, error) {
	s.m.Lock()
	defer s.m.Unlock()

	updated, err := s.Ha.SetBrightness(entityID, brightness)
	if err != nil {
		return nil, err
	}

	if len(updated) == 0 {
		_, s.Err = s.getDevices()
		if s.Err != nil {
			return nil, nil
		}

		updated = append(updated, &ha.Device{ID: entityID, State: true, Brightness: &brightness})
	}

	for _, u := range updated {
		for _, d := range s.Devices {
			if d.ID == u.ID {
				d.State = u.State
				d.Brightness = u.Brightness
				return d, nil
			}
		}
	}

	return nil, errors.New("device not found")
}

func New() *State {
	s := &State{}

	s.Cfg, s.Err = config.Load()
	if s.Err != nil {
		return s
	}

	s.connectClient()

	return s
}
