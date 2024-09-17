package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var stRootDir string
var stSeparator string
var iJsonData map[string]any

const strJsonFileName = "dir.json"

func loadJson() {
	stSeparator = string(filepath.Separator)
	stWorkDir, _ := os.Getwd()
	stRootDir = stWorkDir[:strings.LastIndex(stWorkDir, stSeparator)]
	gnJsonBytes, _ := os.ReadFile(stWorkDir + stSeparator + strJsonFileName)
	err := json.Unmarshal(gnJsonBytes, &iJsonData)
	if err != nil {
		panic("load json fail:" + err.Error())
	}
}

func parseMap(mapData map[string]any, strParentDir string) {
	for k, v := range mapData {
		switch v.(type) {
		case string:
			path, _ := v.(string)
			if path == "" {
				continue
			}

			if strParentDir != "" {
				path = strParentDir + stSeparator + path
				if k == "text" {
					strParentDir = path
				}
			} else {
				strParentDir = path
			}
			createDir(path)
		case []any:
			parseArray(v.([]any), strParentDir)
		}
	}
}

func parseArray(giJsonData []any, strParentDit string) {
	for _, v := range giJsonData {
		mapV, _ := v.(map[string]any)
		parseMap(mapV, strParentDit)
	}
}

func createDir(path string) {
	if len(path) == 0 {
		return
	}

	err := os.MkdirAll(stRootDir+stSeparator+path, os.ModePerm)
	if err != nil {
		panic("create dir fail:" + err.Error())
	}
}

func TestGenerateDir(t *testing.T) {
	loadJson()
	parseMap(iJsonData, "")
}

func TestWriteJsonFile(t *testing.T) {
	var filePath string
	fileDir := "../docs/markdown/"
	dirList := strings.Split(fileDir, "/")
	for _, dir := range dirList {
		filePath = filepath.Join(filePath, dir)
	}
	fmt.Println(filePath)
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		panic(err)
	}
	filePath = filepath.Join(filePath, "dir.json")
	now := time.Now()
	offset := TimeOffset{
		StartTime: now.AddDate(0, 0, -30).Unix(),
		EndTime:   now.Unix(),
	}
	marshal, err := json.MarshalIndent(offset, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filePath, marshal, os.ModePerm); err != nil {
		panic(err)
	}
}

type TimeOffset struct {
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

func TestReadJsonFile(t *testing.T) {
	var fileJson string
	filePath := "../docs/markdown/dir.json"
	pathSplit := strings.Split(filePath, "/")
	for _, file := range pathSplit {
		fileJson = filepath.Join(fileJson, file)
	}
	fmt.Println(fileJson)
	file, err := os.ReadFile(fileJson)
	if err != nil {
		panic(err)
	}
	var offset TimeOffset
	if err := json.Unmarshal(file, &offset); err != nil {
		panic(err)
	}
	fmt.Println(offset)
	offset.StartTime = offset.EndTime
	unix := time.Unix(offset.StartTime, 0)
	fmt.Println(unix)
	date := unix.AddDate(0, 0, 30).Unix()
	fmt.Println(date)
}
