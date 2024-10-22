package main

import (
	"encoding/json"
	"fmt"
)

// "ruleId"
var SubHeaders = []string{"areaCate", "subTotalTime", "START", "basicInfoTime", "preOpeTime", "checkDefectTime", "saveDefectTime", "sitoTime", "getGlassTime", "sheetTupleListTime", "redisRuleKeysTime", "getRuleBySheetsTime", "getSampCntTime", "checkRuleTime", "sendAlarmTime", "toFileCheckOutTime", "FINISH"}

const mainHeaderJSON = `[
	{
		"Name": "主线程开始",
		"Key": "START"
	},
	{
		"Name": "basicOpe",
		"Key": "basicOpe"
	},
	{
		"Name": "basicProd",
		"Key": "basicProd"
	},
	{
		"Name": "basicOwnTime",
		"Key": "basicOwn"
	},
	{
		"Name": "rpushTime",
		"Key": "rpushTime"
	},
	{
		"Name": "zaddTime",
		"Key": "zaddTime"
	},
	{
		"Name": "matchedRuleTime",
		"Key": "matchedRuleTime"
	},
	{
		"Name": "zaddHistory",
		"Key": "zaddHistory"
	},
	{
		"Name": "主线程调用子线程(MAIN CALL SUB)",
		"Key": "CALL_SUB"
	},
	{
		"Name": "回主线程 (MAIN RETURN)",
		"Key": "FINISH"
	},
	{
		"Name": "总时间",
		"Key": "TotalTime"
	}
]`

type MainHeader struct {
	Name string
	Key  string
}

func GetMainHeaders() []MainHeader {
	var headers []MainHeader
	err := json.Unmarshal([]byte(mainHeaderJSON), &headers)
	if err != nil {
		fmt.Println(err)
	}
	return headers
}

var MainHeaders = GetMainHeaders()
