package response

// UploadChunkWebsocketBody represents the response to client about the status of uploading a chunk
type UploadChunkWebsocketBody struct {
	Status       string `json:"status"`
	Size         int    `json:"size"`
	Code         int    `json:"code"`
	ErrorMessage string `json:"errorMessage"`
}

const (
	UploadChunkStatusCompleted = "completed"
	UploadChunkStatusSuccess   = "success"
	UploadChunkStatusFailed    = "failed"
)
