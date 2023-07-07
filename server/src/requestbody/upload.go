package requestbody

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
	isValidSize := r.VideoSize > 0
	isValidName := len(r.VideoName) > 0
	return isValidType && isValidName && isValidSize
}
