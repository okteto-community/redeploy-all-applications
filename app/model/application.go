package model

import "time"

// Application represents an application deployed within an Okteto namespace
type Application struct {
	Branch      string    `json:"branch"`
	LastUpdated time.Time `json:"lastUpdated"`
	Name        string    `json:"name"`
	Repository  string    `json:"repository"`
	Status      string    `json:"status"`
}
