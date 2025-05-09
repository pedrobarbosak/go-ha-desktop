package ha

type mock struct {
}

func (m mock) TurnOn(entityID string) ([]*Device, error) {
	return []*Device{}, nil
}

func (m mock) TurnOff(entityID string) ([]*Device, error) {
	return []*Device{}, nil
}

func (m mock) GetDevices() ([]*Device, error) {
	return []*Device{}, nil
}

func (m mock) SetBrightness(entityID string, brightness uint) ([]*Device, error) {
	return []*Device{}, nil
}

func Mock() Client {
	return &mock{}
}
