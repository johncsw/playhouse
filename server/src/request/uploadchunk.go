package request

import "strings"

type UploadChunkBody struct {
	VideoID string `json:"videoID"`
}

// "request" body from WebSocket for client messages about metadata of a chunk
type UploadChunkWebsocketBody struct {
	Size    int `json:"size"`
	Code    int `json:"code"`
	Content []byte
}

func (b *UploadChunkWebsocketBody) isValid() bool {
	isFirstChunk := b.Code == 0
	if isFirstChunk {
		headerBytes := b.Content[:10]
		headerStr := string(headerBytes)
		// header string for mp4 file
		isValidHeader := strings.Contains(headerStr, "ftyp")
		if !isValidHeader {
			return false
		}
	}

	isValidCode := b.Code >= 0
	isValidSize := b.Size >= 0
	isValidContent := len(b.Content) > 0 && len(b.Content) == b.Size

	return isValidCode && isValidSize && isValidContent
}
