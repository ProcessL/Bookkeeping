package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
