package entity

//UserAgent and Ip Struct
type MediaInfo struct {
	Ip        string `json:"ip"`
	Agent     string `json:"agent"`
	MediaGuid string `json:"guid"`
}
