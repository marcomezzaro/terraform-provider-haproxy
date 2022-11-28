package models

type GetBackend struct {
	Version int `json:"_version"`
	Data Backend `json:"data"`
}

type Backend struct {
	Balance struct {
		Algorithm string `json:"algorithm"`
	} `json:"balance"`
	Mode string `json:"mode"`
	Name string `json:"name"`
}