package request

import "playhouse-server/util"

// don't need to make it configurable now
var validVideoTypes = map[string]bool{
	"video/mp4": true,
}

type UploadRegistrationBody struct {
	VideoName string `json:"videoName"`
	VideoSize int    `json:"videoSize"`
	VideoType string `json:"videoType"`
}

func (r UploadRegistrationBody) isValid() bool {
	_, isValidType := validVideoTypes[r.VideoType]
	// don't need to make it configurable now
	isValidSize := r.VideoSize > 0 && r.VideoSize < 3*util.GB
	isValidName := len(r.VideoName) > 0
	return isValidType && isValidName && isValidSize
}
