package object

type MasterAC struct {
	Power       bool    // Power on or off
	Mode        int     // cold0 or hot1
	Temperature float64 // temperature
}

var MasterACInstance = MasterAC{}

func InitMAC(m *MasterAC) {
	m.Power = true
	m.Temperature = 22
	m.Mode = 0
}
