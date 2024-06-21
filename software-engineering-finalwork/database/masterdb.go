package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type MasterCounter struct {
	gorm.Model
	UID              string
	IDCard           string
	Room             int
	Mode             int
	Power            int
	PowerChangeTimes int
	Serve            int //0:stop 1:serve 2:pending
	ServeReqTime     int64
	StartTime        int64
	EndTime          int
	NowSpeed         int //0:stop 1:low 2:mid 3:high
	StartReqTemp     float64
	NowTemp          float64
	StopReqTemp      float64
	Money            float64
	Energy           float64
}

type MasterStatusLogger struct {
	gorm.Model
	Power     bool
	UID       string
	StartTime int64
	EndTime   int
	Speed     int //0:stop 1:low 2:mid 3:high
	ReqTemp   float64
}

type MasterSys struct {
	ID          uint `gorm:"primarykey"`
	NowServeNum int
}

type IDRoom struct {
	IDCard string
	Room   int
}

var DB *gorm.DB

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
		},
	)
	var err error
	DB, err = gorm.Open(sqlite.Open("master.DB"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&MasterCounter{})
	DB.AutoMigrate(&IDRoom{})
	var ir []IDRoom
	if DB.Find(&ir); len(ir) == 0 {
		DB.Create(&IDRoom{IDCard: "123", Room: 101})
		DB.Create(&IDRoom{IDCard: "124", Room: 102})
		DB.Create(&IDRoom{IDCard: "125", Room: 103})
	}
	DB.AutoMigrate(&MasterStatusLogger{})
	DB.AutoMigrate(&MasterSys{})
	var sys []MasterSys
	if DB.Find(&sys); len(sys) == 0 {
		DB.Create(&MasterSys{NowServeNum: 0})
	}
}
