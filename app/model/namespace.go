package model

type Status string

const (
	Active           = "Active"
	DestroyAllFailed = "DestroyAllFailed"
	DestroyingAll    = "DestroyingAll"
	Deleting         = "Deleting"
	Inactive         = "Inactive"
	Sleeping         = "Sleeping"
	DeleteFailed     = "DeleteFailed"
)

// Namespace represents an Okteto namespace
type Namespace struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}
