package models

type GetBind struct {
	Version int `json:"_version"`
	Data Bind `json:"data"`
}

type Bind struct {
	Address string `json:"address"`
	Port int `json:"port"`
	Name string `json:"name"`
}
