package util

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func IsInWeekdaySchedule(now time.Time, startTime string, endTime string,
	mon bool, tue bool, wed bool, thu bool, fri bool, sat bool, sun bool) bool {


	refNow := (int(now.Weekday()) * 1440) + (int(now.Hour()) * 60) + (int(now.Minute()))
	h, m := parseTimeString(startTime)
	refStart := (h * 60) + m
	h, m = parseTimeString(endTime)
	refEnd := (h * 60) + m

	if mon {
		if checkDayInSchedule(1, refStart, refEnd, refNow) {
			return true
		}
	}
	if tue {
		if checkDayInSchedule(2, refStart, refEnd, refNow) {
			return true
		}
	}
	if wed {
		if checkDayInSchedule(3, refStart, refEnd, refNow) {
			return true
		}
	}
	if thu {
		if checkDayInSchedule(4, refStart, refEnd, refNow) {
			return true
		}
	}
	if fri {
		if checkDayInSchedule(5, refStart, refEnd, refNow) {
			return true
		}
	}
	if sat {
		if checkDayInSchedule(6, refStart, refEnd, refNow) {
			return true
		}
	}
	if sun {
		if checkDayInSchedule(0, refStart, refEnd, refNow) {
			return true
		}
	}
	return false

}

func checkDayInSchedule(day int, refTimeStart int, refTimeEnd int, refNow int) bool {

	start := (day * 1440) + refTimeStart
	end := (day * 1440) + refTimeEnd

	//add day overflow
	if end <= start {
		end += 1440
	}

	//fix end of week overflow
	if (end >= 7 * 1440) && (refNow < refTimeStart){
		refNow += 7 * 1440
	}

	if refNow < start {return false}

	if refNow >= end {return false}
	return true

}

func parseTimeString(timeStr string) (int, int) {
	split := strings.Split(timeStr, ":")
	if len(split) != 2 {
		log.Fatal("Failed to parse time: " + timeStr)
	}

	hr, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		log.Fatal("Failed to parse time: " + timeStr)
	}
	min, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		log.Fatal("Failed to parse time: " + timeStr)
	}
	return int(hr), int(min)

}
