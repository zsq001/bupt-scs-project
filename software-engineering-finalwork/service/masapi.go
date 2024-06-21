package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	db "softengineering/database"
	"softengineering/master"
	mw "softengineering/object"
	"strconv"
	"time"
)

func GenRandUID() string {
	bytes := make([]byte, 4) // 4 bytes * 2 (for hex representation) = 8 hex digits
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes) // prints a random 8 digit hexadecimal number
}

func RegisterSlave(c *gin.Context) {
	var mc db.MasterCounter
	room, err := strconv.Atoi(c.Query("room"))
	IDCard := c.Query("IDCard")
	if IDCard == "" || err != nil {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "IDCard and Room is required",
		})
		return
	}
	var idroom db.IDRoom
	if res := db.DB.Where("room = ?", room).First(&idroom); errors.Is(res.Error, gorm.ErrRecordNotFound) || idroom.IDCard != IDCard {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "Room not found or incorrect IDCard",
		})
		return
	}
	if res := db.DB.Where("room = ?", room).Last(&mc); !errors.Is(res.Error, gorm.ErrRecordNotFound) && mc.EndTime == 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "Room already registered",
		})
		return
	}
	var mcc db.MasterCounter
	mcc.Power = 0
	mcc.PowerChangeTimes = 0
	mcc.UID = GenRandUID()
	mcc.Room = room
	mcc.Mode = mw.MasterACInstance.Mode
	mcc.NowTemp = mw.MasterACInstance.Temperature
	mcc.IDCard = IDCard
	mcc.StartTime = time.Now().Unix()
	mcc.Energy = 0
	mcc.StartReqTemp = mw.MasterACInstance.Temperature
	mcc.StopReqTemp = 0
	mcc.Money = 0
	db.DB.Create(&mcc)
	c.JSON(200, mw.WebResponse{
		Code: 200,
		Data: mcc.UID,
	})
}

func ModifyAC(c *gin.Context) {
	uid := c.PostForm("uid")
	temp := c.PostForm("temp")
	mode := c.PostForm("mode")    // cold hot
	speed := c.PostForm("speed")  //0 1 2 3
	power1 := c.PostForm("power") //on off
	power, _ := strconv.Atoi(power1)
	var mc db.MasterCounter
	if db.DB.Where("uid = ?", uid).First(&mc); mc.EndTime != 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "UID not found",
		})
		return
	}
	//fmt.Printf("", mc)
	//var mcc []db.MasterCounter
	//if db.DB.Where("end_time <> 0 AND now_speed <> 0 AND power = 1").Find(&mcc); len(mcc) == 3 {
	//	c.JSON(406, mw.WebResponse{
	//		Code:  http.StatusNotAcceptable,
	//		Error: "Please wait for other users to finish",
	//	})
	//	return
	//}
	if modeint, _ := strconv.Atoi(mode); modeint != mw.MasterACInstance.Mode {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "Mode is not match",
		})
		return
	}
	if tempint, _ := strconv.Atoi(temp); (mode == "cold" && (tempint > 25 || tempint < 18)) || (mode == "hot" && (tempint > 30 || tempint < 25)) {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "Temperature is out of range",
		})
		return
	}
	if speedint, err := strconv.Atoi(speed); speedint < 0 || speedint > 3 || err != nil {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "Speed is invalid",
		})
		return

	}
	var ml db.MasterStatusLogger
	if res := db.DB.Last(&ml, "uid = ?", uid); !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		ml.EndTime = int(time.Now().Unix())
		db.DB.Save(&ml)
	}
	var ml2 db.MasterStatusLogger
	ml2.Power = power == 1
	ml2.StartTime = time.Now().Unix()
	ml2.Speed, _ = strconv.Atoi(speed)
	ml2.ReqTemp, _ = strconv.ParseFloat(temp, 64)
	ml2.UID = uid
	db.DB.Create(&ml2)
	s2, _ := strconv.Atoi(speed)
	if power == 0 || s2 == 0 {
		var sys db.MasterSys
		db.DB.Find(&sys)
		if mc.Serve == 1 {
			sys.NowServeNum--
		}
		db.DB.Save(&sys)
		mc.Serve = 0
	}
	mc.PowerChangeTimes++
	mc.StartTime = time.Now().Unix()
	mc.NowTemp, _ = strconv.ParseFloat(temp, 64)
	fmt.Printf("", temp)
	mc.NowSpeed, _ = strconv.Atoi(speed)
	mc.Power = power
	db.DB.Save(&mc)
	c.JSON(200, mw.WebResponse{
		Code: 200,
		Data: mc,
	})
}

func ServeSlave(c *gin.Context) {
	uid := c.PostForm("uid")
	power1 := c.PostForm("power")    //on off
	power, _ := strconv.Atoi(power1) // 1 0
	var mc db.MasterCounter
	if db.DB.Where("uid = ?", uid).First(&mc); uid == "" || mc.EndTime != 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "UID not found or already log out",
		})
		return
	}

	var sys db.MasterSys
	db.DB.Find(&sys)
	if power == 0 {
		if mc.Serve == 0 {
			c.JSON(200, mw.WebResponse{
				Code: 200,
				Data: mc.Serve,
			})
			return
		}
		db.DB.Find(&sys)
		mc.Serve = 0
		db.DB.Save(&mc)
		c.JSON(200, mw.WebResponse{
			Code: 201,
			Data: mc.Serve,
		})
		sys.NowServeNum--
		db.DB.Save(&sys)
		return
	}
	if mc.Serve == 1 {
		c.JSON(200, mw.WebResponse{
			Code: 202,
			Data: mc.Serve,
		})
		db.DB.Save(&mc)
		return
	}
	if db.DB.Find(&sys); sys.NowServeNum >= 3 {
		if mc.Serve == 2 {
			c.JSON(200, mw.WebResponse{
				Code:  http.StatusOK,
				Error: "Still pending",
			})
			return
		}
		mc.Serve = 2
		mc.ServeReqTime = time.Now().Unix()
		c.JSON(200, mw.WebResponse{
			Code:  http.StatusOK,
			Error: "Please wait for other users to finish",
			Data:  mc.Serve,
		})
		return
	}
	sys.NowServeNum++
	db.DB.Save(&sys)
	mc.Serve = 1
	db.DB.Save(&mc)
	c.JSON(200, mw.WebResponse{
		Code: 200,
		Data: mc.Serve,
	})

}
func PendingSlaveServe(c *gin.Context) {
	uid := c.Query("uid")
	var mc db.MasterCounter
	if db.DB.Where("uid = ?", uid).First(&mc); uid == "" || mc.EndTime != 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "UID not found or already log out",
		})
		return
	}
	if mc.Serve != 2 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "Not pending",
		})
		return
	}
	var sys db.MasterSys
	db.DB.Find(&sys)
	if db.DB.Find(&sys); sys.NowServeNum >= 3 {
		c.JSON(200, mw.WebResponse{
			Code:  http.StatusOK,
			Error: "Please wait for other users to finish",
		})
	} else {
		sys.NowServeNum++
		db.DB.Save(&sys)
		mc.Serve = 1
		db.DB.Save(&mc)
		c.JSON(200, mw.WebResponse{
			Code: 200,
			Data: mc.Serve,
		})
	}
}

func StopSlaveServe(c *gin.Context) {
	uid := c.PostForm("uid")
	var mc db.MasterCounter
	if db.DB.Where("uid = ?", uid).First(&mc); uid == "" || mc.EndTime != 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "UID not found or already log out",
		})
		return
	}
	if mc.Serve == 1 {
		var sys db.MasterSys
		db.DB.Find(&sys)
		sys.NowServeNum--
		db.DB.Save(&sys)

	}
	mc.Serve = 0
	db.DB.Save(&mc)
	c.JSON(200, mw.WebResponse{
		Code: 200,
	})
}

func SlaveLogOut(c *gin.Context) {
	uid := c.Query("uid")
	var mc db.MasterCounter
	if db.DB.Where("uid = ?", uid).First(&mc); uid == "" || mc.EndTime != 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "UID not found or already log out",
		})
		return
	}
	if mc.Serve == 1 {
		var sys db.MasterSys
		db.DB.Find(&sys)
		sys.NowServeNum--
		db.DB.Save(&sys)
		mc.Serve = 0
	}
	mc.StopReqTemp = mc.NowTemp
	mc.EndTime = int(time.Now().Unix())
	db.DB.Save(&mc)
	c.JSON(200, mw.WebResponse{
		Code: 200,
	})
}

func Info(c *gin.Context) {
	c.JSON(200, mw.WebResponse{
		Code: 200,
		Data: mw.MasterACInstance,
	})
}

func SlaveStatus(c *gin.Context) {
	uid := c.Query("uid")
	var mc db.MasterCounter
	if res := db.DB.Where("uid = ?", uid).First(&mc); errors.Is(res.Error, gorm.ErrRecordNotFound) || uid == "" || mc.EndTime != 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "UID not found or already log out",
		})
		return
	}
	c.JSON(200, mw.WebResponse{
		Code: 200,
		Data: mc,
	})
}

func Report(c *gin.Context) {
	syear, _ := strconv.Atoi(c.Query("syear"))
	smonth, _ := strconv.Atoi(c.Query("smonth"))
	sday, _ := strconv.Atoi(c.Query("sday"))
	eyear, _ := strconv.Atoi(c.Query("eyear"))
	emonth, _ := strconv.Atoi(c.Query("emonth"))
	eday, _ := strconv.Atoi(c.Query("eday"))
	if syear == 0 || smonth == 0 || sday == 0 || eyear == 0 || emonth == 0 || eday == 0 {
		c.JSON(400, mw.WebResponse{
			Code:  400,
			Error: "Invalid date",
		})
		return
	}
	csvFileName := master.GetReportByTime(syear, smonth, sday, eyear, emonth, eday)
	c.Header("Content-Type", "application/octet-stream")
	c.FileAttachment(csvFileName, csvFileName)
	err := os.Remove(csvFileName)
	if err != nil {
		// 处理错误
		fmt.Println("Error while deleting file: ", err)
	}
}

func EditMaster(c *gin.Context) {
	mode := c.Query("mode")
	temp := c.Query("temp")
	power := c.Query("power")
	var Mode int
	var Power bool
	var Temperature float64
	if mode != "" {
		Mode, _ = strconv.Atoi(mode)
	}
	if temp != "" {
		Temperature, _ = strconv.ParseFloat(temp, 64)
		if (Mode == 0 && (Temperature < 18 || Temperature > 25)) || (Mode == 1 && (Temperature < 25 || Temperature > 30)) {
			c.JSON(400, mw.WebResponse{
				Code:  400,
				Error: "Temperature is out of range",
			})
			return
		}
	}
	if power != "" {
		Power = power == "on"
	}
	mw.MasterACInstance.Mode = Mode
	mw.MasterACInstance.Temperature = Temperature
	mw.MasterACInstance.Power = Power
	var mc []db.MasterCounter
	db.DB.Where("end_time <> 0").Find(&mc)
	for i := 0; i < len(mc); i++ {
		if Mode == 0 {
			if mc[i].NowTemp < 18 || mc[i].NowTemp > 25 {
				mc[i].NowTemp = 25
			}
		} else {
			if mc[i].NowTemp < 25 || mc[i].NowTemp > 30 {
				mc[i].NowTemp = 25
			}
		}
		mc[i].Mode = Mode
	}
	db.DB.Save(&mc)
	c.JSON(200, mw.WebResponse{
		Code: 200,
		Data: mw.MasterACInstance,
	})
}
