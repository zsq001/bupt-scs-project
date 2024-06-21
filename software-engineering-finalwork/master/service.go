package master

import (
	"errors"
	"softengineering/object"
)

type mAC struct {
	*object.MasterAC
}

func (m mAC) switchPower() {
	m.Power = !m.Power
}

func (m mAC) switchMode(mode int) {
	m.Mode = mode
}

func (m mAC) switchTemperature(temperature float64) error {
	if m.Mode == 0 && (temperature < 18 || temperature > 25) {
		return errors.New("cold mode temperature must be between 18 and 25")
	}
	if m.Mode == 1 && (temperature < 25 || temperature > 30) {
		return errors.New("hot mode temperature must be between 25 and 30")
	}
	m.Temperature = temperature
	return nil
}
