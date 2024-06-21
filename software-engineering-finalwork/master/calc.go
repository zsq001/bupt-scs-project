package master

import (
	"math"
	"softengineering/database"
	"softengineering/object"
)

func CalcPrice() {
	if object.MasterACInstance.Power == false {
		return
	}
	//var lowprice = 0.06666666666666667
	//var midprice = 0.08333333333333333
	//var highprice = 0.1
	var mastercounter []database.MasterCounter
	database.DB.Where("end_time = 0").Find(&mastercounter)
	for i := 0; i < len(mastercounter); i++ {
		if mastercounter[i].Serve == 0 || mastercounter[i].Power == 0 {
			continue
		}
		if mastercounter[i].NowSpeed == 1 {
			mastercounter[i].Energy += 0.013
			mastercounter[i].Energy = math.Round(mastercounter[i].Energy*100) / 100
			mastercounter[i].Money = 5.0 * mastercounter[i].Energy
			mastercounter[i].Money = math.Round(mastercounter[i].Money*100) / 100
		} else if mastercounter[i].NowSpeed == 2 {
			mastercounter[i].Energy += 0.016
			mastercounter[i].Energy = math.Round(mastercounter[i].Energy*100) / 100
			mastercounter[i].Money = 5.0 * mastercounter[i].Energy
			mastercounter[i].Money = math.Round(mastercounter[i].Money*100) / 100
		} else if mastercounter[i].NowSpeed == 3 {
			mastercounter[i].Energy += 0.02
			mastercounter[i].Energy = math.Round(mastercounter[i].Energy*100) / 100
			mastercounter[i].Money = 5.0 * mastercounter[i].Energy
			mastercounter[i].Money = math.Round(mastercounter[i].Money*100) / 100
		}
	}
	database.DB.Save(&mastercounter)
}
