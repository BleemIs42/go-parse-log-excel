package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	file, err := ioutil.ReadFile(`./222.log`)
	if err == nil {
		content := string(file)
		mainThread := matchMainThread(content)
		for key, dateTotal := range mainThread {
			fmt.Println(key, dateTotal)
		}

		fmt.Println("-----")

		subThread := matchSubThread(content)
		for pid, keyDate := range subThread {
			for key, date := range keyDate {
				fmt.Println(pid, key, date)
			}
		}
	}
}

func matchMainThread(file string) map[string][]string {
	timeKeyReg := regexp.MustCompile(`MAIN_(\w+)\]\s+time:(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}.\d{3})`)
	timeKeyMatched := timeKeyReg.FindAllStringSubmatch(file, -1)

	timeKeyTotalLineReg := regexp.MustCompile(`.*MAIN_TIME.*\s(\S+)`)
	timeKeyTotalLineMatched := timeKeyTotalLineReg.FindStringSubmatch(file)

	timeKeyTotalReg := regexp.MustCompile(`(\w+):(\d+)`)
	timeKeyTotalMatched := timeKeyTotalReg.FindAllStringSubmatch(timeKeyTotalLineMatched[1], -1)

	mainTread := map[string][]string{}
	for _, keyDate := range timeKeyMatched {
		mainTread[keyDate[1]] = []string{keyDate[2]}
		fmt.Println(keyDate[1], mainTread[keyDate[1]])
	}
	fmt.Println("-----")

	for _, keyTotal := range timeKeyTotalMatched {
		key, total := keyTotal[1], keyTotal[2]
		mainTread[key] = append(mainTread[key], total)
	}

	return mainTread
}

func matchSubThread(file string) map[string]map[string][]string {
	timeKeyReg := regexp.MustCompile(`SUB_(\w+)_(\w+)\]\s+time:(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}.\d{3})`)
	timeKeyMatched := timeKeyReg.FindAllStringSubmatch(file, -1)

	subThread := map[string]map[string][]string{}
	for _, pidKeyDate := range timeKeyMatched {
		pid, key, date := pidKeyDate[1], pidKeyDate[2], pidKeyDate[3]
		// fmt.Println(pid, key, date)
		if _, pidOk := subThread[pid]; pidOk {
			if _, keyOk := subThread[pid][key]; !keyOk {
				subThread[pid][key] = []string{date}
			}
		} else {
			subThread[pid] = map[string][]string{key: []string{date}}
		}
	}

	timeTotalKeyLineReg := regexp.MustCompile(`.*SUB_TIME_(\d+)\S+\s(.*)`)
	timeTotalKeyLineMatched := timeTotalKeyLineReg.FindAllStringSubmatch(file, -1)

	timeKeyTotalReg := regexp.MustCompile(`(\w+):(\d+)`)
	for _, keyTotalLineMatched := range timeTotalKeyLineMatched {
		pid, keyTotalLine := keyTotalLineMatched[1], keyTotalLineMatched[2]
		timeKeyTotalMatched := timeKeyTotalReg.FindAllStringSubmatch(keyTotalLine, -1)
		for _, keyTotal := range timeKeyTotalMatched {
			key, total := keyTotal[1], keyTotal[2]
			if _, ok := subThread[pid][key]; ok {
				subThread[pid][key] = append(subThread[pid][key], total)
			}
		}
	}

	return subThread
}

func createExcel() {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "Hello world.")

	err := f.SaveAs("./log-parsed.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}
