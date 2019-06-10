package types

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type DateTime struct {
	time.Time
}

func (dt *DateTime) Scan(v interface{}) error {
	var nv string
	switch v.(type) {
		case []byte:
			tv, _ := v.([]byte)
			nv = string (tv)
		case string:
			nv, _ = v.(string)
		default:
			return errors.New("sql: Source type not supported.")
	}

	var errCheckTimezone error
	useTimeZone := false
	formatString := "2006-01-02 15:04:05"
	useTimeZone, errCheckTimezone = regexp.MatchString("[-+][0-9][0-9]$", nv)
	if errCheckTimezone != nil {
		return errors.New(fmt.Sprint("sql: Regular expression error:", errCheckTimezone))
	}
	if strings.Contains(nv, ".") == false {
		if useTimeZone == true {
			formatString += "-07"
		}
		newTime, errTimeParse := time.Parse(formatString, nv)
		if errTimeParse != nil {
			return errTimeParse
		}
		dt.Time = newTime
	} else {
		tempZeroPaddingCount := 0
		if useTimeZone == true {
			tempZeroPaddingCount = len(nv) - strings.Index(nv, ".") - 4
		} else {
			tempZeroPaddingCount = len(nv) - strings.Index(nv, ".") - 1
		}
		tempZeroPadding := ""
		for paddingCount := 0; paddingCount < tempZeroPaddingCount; paddingCount++ {
			tempZeroPadding = tempZeroPadding + "0"
		}
		formatString += "." + tempZeroPadding
		if useTimeZone == true {
			formatString += "-07"
		}
		newTime, errTimeParse := time.Parse(formatString, nv)
		if errTimeParse != nil {
			return errTimeParse
		}
		dt.Time = newTime
	}
	return nil
}
