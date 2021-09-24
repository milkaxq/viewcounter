package entity

//UserAgent and Ip Struct
type ResInfo struct {
	Ip          string `json:"ip"`
	Agent       string `json:"agent"`
	ProductGuid string `json:"guid"`
	ResNo       string `json:"resno"`
}
