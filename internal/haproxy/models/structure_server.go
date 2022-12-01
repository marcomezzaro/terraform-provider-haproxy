package models

type GetServer struct {
	Version int `json:"_version"`
	Data Server `json:"data"`
}
	
type Server struct {
	Check          string      `json:"check"`
	Address        string      `json:"address"`
	Name           string      `json:"name"`
	Port           int         `json:"port"`
}
