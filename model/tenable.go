package model

import "github.com/dotdancer/gogofly/global"

type Tenable struct {
	global.GlyModel
	Name         string `json:"name"`
	Description  string `json:"description"`
	Details      string `json:"details"`
	Status       int    `json:"status"`
	ImportStatus string `json:"import_status"`
}

func (t *Tenable) TableName() string {
	return "tenable"
}
