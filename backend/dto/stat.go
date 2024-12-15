package dto

type StatResponse struct {
	Progress []Progress `json:"progress"`
	Goal     int64      `json:"goal"`
}
