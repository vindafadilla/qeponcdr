package controllers

import (
	"github.com/vindafadilla/qeponcdr/qeponcdr/models"
	
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin"

	"time"
)

func GetDataperMonth(ctx *gin.Context) {
	results, day := getArray(ctx)
    resultspermonth := models.StatperMonth{}

    for i := 0; i <= (day-1); i++ {
    	// ctx.String(200, "%v",results)
    	resultspermonth.TotalUser = resultspermonth.TotalUser + results[i].StatUser.All
    	resultspermonth.TotalUserVer = resultspermonth.TotalUserVer + results[i].StatUser.AllVer
    	resultspermonth.TotalCA = resultspermonth.TotalCA + results[i].StatCDR.Init
    }

    resultspermonth.AvgCA = float64(resultspermonth.TotalCA)/ float64(day)
    resultspermonth.AvgCAUser = float64(resultspermonth.TotalCA)/ float64(resultspermonth.TotalUserVer)/ float64(day)
    resultspermonth.TotalDay = day
    ctx.JSON(200, resultspermonth)
}

func GetAllDataperMonth(ctx *gin.Context) {
	results, _ := getArray(ctx)

	ctx.JSON(200, results)
}

func getArray(ctx *gin.Context)([]models.StatperDays, int) {
	startTime,endTime, day := getTime(ctx)

	usercoll:=initializeDB(ctx,"data_intermediate")

	results := []models.StatperDays{}

	err := usercoll.Find(bson.M{"date": bson.M{"$gte": startTime, "$lte": endTime}}).All(&results)
    if err != nil {
        ctx.String(400, "%s",err)
    }

    for i := 0; i < len(results); i++ {
    	results[i].Date = results[i].Date.UTC()    	
    }

    	return results, day
    
}

func getTime(ctx *gin.Context)(time.Time, time.Time, int) {
    
    Starttime := ctx.Query("starttime")
    Endtime := ctx.Query("endtime")
    
    loc, _ := time.LoadLocation("Europe/Berlin")
    const shortForm = "20060102150405"
    // StartTime, err :=time.Parse(shortForm, Starttime)
    // if err != nil {
    //     ctx.String(400, "%s",err)
    // }
    // EndTime, err := time.Parse(shortForm, Endtime)
    // if err != nil {
    //     ctx.String(400, "%s",err)
    // }

    StartTime, err :=time.ParseInLocation(shortForm, Starttime, loc)
    if err != nil {
        ctx.String(400, "%s",err)
    }
    // StartTime = StartTime.UTC()
    EndTime, err := time.ParseInLocation(shortForm, Endtime, loc)
    if err != nil {
        ctx.String(400, "%s",err)
    }
    // EndTime = EndTime.UTC()

    day := EndTime.YearDay() - StartTime.YearDay()+1

    return StartTime, EndTime, day
}

func GetDataperDay(ctx *gin.Context) {
	startTime,endTime, _ := getTime(ctx)

	usercoll:=initializeDB(ctx,"data_intermediate")

	results := models.StatperDays{}

	err := usercoll.Find(bson.M{"date": bson.M{"$gte": startTime, "$lte": endTime}}).One(&results)
    if err != nil {
        ctx.String(400, "%s",err)
    }else{
    	results.Date = results.Date.UTC()
    	ctx.JSON(200, results)
    }
}