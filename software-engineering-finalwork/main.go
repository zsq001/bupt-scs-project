package main

import (
	"github.com/robfig/cron/v3"
	"softengineering/database"
	"softengineering/master"
	"softengineering/object"
	"softengineering/router"
)

func main() {

	// master start
	database.InitDB()

	object.InitMAC(&object.MasterACInstance)
	c := cron.New(cron.WithSeconds())
	c.AddFunc("*/1 * * * * *", master.CalcPrice)
	c.Start()
	router.InitRouter()

	// master end

}
