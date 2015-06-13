package service

import (
	"time"
)

type Token struct {
	ID          string    `json:"id" gorethink:"id"`
	Restriction []string  `json:"restriction" gorethink:"restriction"`
	LastUse     time.Time `json:"last_use" gorethink:"last_use"`
}

type Script struct {
	ID   string `json:"id" gorethink:"id"`
	Code string `json:"code" gorethink:"code"`
}
