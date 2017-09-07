package models

import(
	"time"
)

//Total User and New User struct
type StatUser struct {
        All        int `json:"all" bson:"all"`
        AllVer     int `json:"allver" bson:"allver"`
        InRange    int `json:"inrange" bson:"inrange"`
        InRangeVer int `json:"inrangever" bson:"inrangever"`
}

//Total Call Attempts struct
type StatCDR struct {
    Init   int `json:"init" bson:"init"`
    Ring   int `json:"ring" bson:"ring"`
    Accept int `json:"accept" bson:"accept"`
}

//Total and Average Call Attempt struct
type StatCDRAvg struct{
    TotalInit   int `json:"totalinit" bson:"totalinit"`
    AvgInit   float64 `json:"avginit" bson:"avginit"`
}

type UserDevice struct{
    TotalLenovo int `json:"totallenovo" bson:"totallenovo"`
    TotalUnknown int `json:"totalunknown" bson:"totalunknown"`
    TotalSamsung int `json:"totalsamsung" bson:"totalsamsung"`
    TotalSony int `json:"totalsony" bson:"totalsony"`
    TotalAndromax int `json:"totalandromax" bson:"totalandromax"`
}

//User call attempts average level
type UserAvgLevel struct{
    GTFive int `json:"gtfive"`
    OneTillFive int `json:"onetillfive"`
    ZPThreeTillOne int `json:"zpthreetillone"`
    ZPOneTillZPThree int `json:"zponetillzpthree"`
    LTZPOne int `json:"ltzpone"`

}

type StatperDays struct{
    Date time.Time `bson:"date" json:"date"`
    StatUser `bson:"statuser" json:"statuser"`
    StatCDR `bson:"statcdr" json:"statcdr"`
    StatCDRAvg `bson:"statcdravg" json:"statcdravg"`
    UserDevice `bson:"userdevice" json:"userdevice"`
    UserAvgLevel `bson:"useravglevel" json:"useravglevel"`
}

type StatperMonth struct{
    AvgCA float64 `bson:"avgca" json:"avgca"`
    AvgCAUser float64 `bson:"avgcauser" json:"avgcauser"`
    TotalUser int `bson:"totaluser" json:"totaluser"`
    TotalUserVer int `bson:"totaluserver" json:"totaluserver"`
    TotalCA int `bson:"totalca" json:"totalca"`
    TotalDay int `bson:"totalday" json:"totalday"`
}