package request

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
	return b.Code >= 0 && b.Size >= 0 && len(b.Content) > 0 && len(b.Content) == b.Size
}
