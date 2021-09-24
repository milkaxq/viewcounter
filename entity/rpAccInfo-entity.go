package entity

//UserAgent and Ip Struct
type RpAccInfo struct {
	Ip        string `json:"ip"`
	Agent     string `json:"agent"`
	RpAccGuid string `json:"guid"`
	RpAccNo   string `json:"resno"`
}
