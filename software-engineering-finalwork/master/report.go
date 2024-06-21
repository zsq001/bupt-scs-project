package master

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	db "softengineering/database"
	"strconv"
	"time"
)

func GetReportByTime(syear int, smonth int, sday int, eyear int, emonth int, eday int) string {
	t := time.Now()
	csvFileName := "report" + t.Format("2006-01-02_15:04:05") + ".csv"
	file, err := os.Create(csvFileName)
	checkError(err)
	defer file.Close()

	startTime := time.Date(syear, time.Month(smonth), sday, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(eyear, time.Month(emonth), eday, 0, 0, 0, 0, time.UTC).Unix()
	//fmt.Printf("%d,%d", startTime, endTime)
	var result []db.MasterCounter
	db.DB.Where("start_time >= ? AND start_time < ?", startTime, endTime).Find(&result)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"Room", "IDCard", "StartTime", "EndTime", "NowSpeed", "StartReqTemp", "StopReqTemp", "Money"}
	writer.Write(header)
	for _, res := range result {
		err = writer.Write([]string{strconv.Itoa(res.Room), res.IDCard, strconv.FormatInt(res.StartTime, 10), strconv.Itoa(res.EndTime), strconv.Itoa(res.NowSpeed), fmt.Sprintf("%f", res.StartReqTemp), fmt.Sprintf("%f", res.StopReqTemp), fmt.Sprintf("%f", res.Money)})
		checkError(err)
	}
	return csvFileName
}

func GetReporrtByUID(uid string) string {
	t := time.Now()
	csvFileName := "report" + t.Format("2006-01-02_15:04:05") + ".csv"
	file, err := os.Create(csvFileName)
	checkError(err)
	defer file.Close()
	var result []db.MasterCounter
	db.DB.Where("uid = ?", uid).First(&result)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"Room", "IDCard", "StartTime", "EndTime", "NowSpeed", "StartReqTemp", "StopReqTemp", "Money"}
	writer.Write(header)
	for _, res := range result {
		err = writer.Write([]string{strconv.Itoa(res.Room), res.IDCard, strconv.FormatInt(res.StartTime, 10), strconv.Itoa(res.EndTime), strconv.Itoa(res.NowSpeed), fmt.Sprintf("%f", res.StartReqTemp), fmt.Sprintf("%f", res.StopReqTemp), fmt.Sprintf("%f", res.Money)})
		checkError(err)
	}
	checkError(err)
	return csvFileName

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
