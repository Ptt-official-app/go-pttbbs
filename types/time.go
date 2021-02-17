package types

import (
	"time"
	"unsafe"
)

var (
	ZERO_TIME       = time.Unix(int64(0), 0)
	ZERO_LOCAL_TIME = ZERO_TIME.In(TIMEZONE)
	ZERO_TIME4      = Time4(0)
)

type Time4 int32

//XXX check whether INT32_SZ should be TIME4_SZ
const TIME4_SZ = unsafe.Sizeof(Time4(0))

func NowTS() Time4 {
	// We don't need to worry about time-zone when using unix-timestamp.
	return TimeToTime4(time.Now())
}

func TimeToTime4(t time.Time) Time4 {
	return Time4(t.Unix())
}

//ToLocal
//
//Instead of using Local, we specify TIME_LOCATION
//to avoid the confusion. (also good for tests)
func (t Time4) ToLocal() time.Time {
	return time.Unix(int64(t), 0).In(TIMEZONE)
}

func (t Time4) ToUtc() time.Time {
	return time.Unix(int64(t), 0).UTC()
}

//Cdate
//
//Print date-time in string.
//23+1 bytes, "01/31/2007 00:00:00 Mon\0"
func (t Time4) Cdate() string {
	return t.ToLocal().Format("01/02/2006 15:04:05 Mon")
}

//Cdatelite
//
//Light-print date-time in string.
//19+1 bytes, "01/31/2007 00:00:00\0"
func (t Time4) Cdatelite() string {
	return t.ToLocal().Format("01/02/2006 15:04:05")
}

//Cdatedate
//
//Print date in string.
//10+1 bytes, "01/31/2007\0"
func (t Time4) Cdatedate() string {
	return t.ToLocal().Format("01/02/2006")
}

//CdateMd
//
//Print month/day in string
//5+1 bytes, "01/31\0"
func (t Time4) CdateMd() string {
	return t.ToLocal().Format("01/02")
}

//Cdatemd
//
//Print month/day in string
//5+1 bytes, " 1/31\0"
func (t Time4) Cdatemd() (str string) {
	str = t.ToLocal().Format("1/02")
	if len(str) == 4 {
		str = " " + str
	}
	return str
}

//CdateMdHM
//
//Print month/day hour:minute in string
//11+1 bytes, "01/31 10:01\0"
func (t Time4) CdateMdHM() string {
	return t.ToLocal().Format("01/02 15:04")
}

//CdateMdHMS
//
//Print month/day hour:minute:second in string
//13+1 bytes, "01/31 10:01:01\0"
func (t Time4) CdateMdHMS() string {
	return t.ToLocal().Format("01/02 15:04:05")
}

//Ctime
//Sun Feb 14 13:33:27 2021
func (t Time4) Ctime() string {
	return t.ToLocal().Format("Mon Jan _2 15:04:05 2006")
}
