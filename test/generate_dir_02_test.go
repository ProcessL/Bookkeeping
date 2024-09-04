package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Node struct {
	Text     string `json:"text"`
	Children []Node `json:"children"`
}

var stRootDir02 string
var stSeparator02 string
var iRootNode Node

const jsonFileName = "dir.json"

func loadJson02() {
	stSeparator02 = string(filepath.Separator)
	stWorkDir02, _ := os.Getwd()
	stRootDir02 = stWorkDir02[:strings.LastIndex(stWorkDir02, stSeparator02)]
	gnJsonfileBytes, _ := os.ReadFile(stWorkDir02 + stSeparator02 + jsonFileName)
	err := json.Unmarshal(gnJsonfileBytes, &iRootNode)
	if err != nil {
		fmt.Println("load json err:", err)
	}
}

func pareNode(node *Node, stParentDir string) {
	if len(node.Text) != 0 {
		createDir02(node.Text, stParentDir)
	}

	if len(node.Children) != 0 {
		stParentDir = stParentDir + stSeparator02 + node.Text
	}

	for _, child := range node.Children {
		pareNode(&child, stParentDir)
	}
}

func createDir02(path string, stParentDir string) {

	if len(stParentDir) != 0 {
		path = stParentDir + stSeparator02 + path
	}

	err := os.MkdirAll(stRootDir02+stSeparator02+path, os.ModePerm)
	if err != nil {
		fmt.Println("create dir err:", err)
	}
}

func TestGenerate02(t *testing.T) {
	loadJson02()
	pareNode(&iRootNode, "")
}
