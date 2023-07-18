package responsebody

// "response" body from WebSocket for client messages about status of upload a chunk
type UploadChunkWSBody struct {
	Status       string `json:"status"`
	Size         int    `json:"size"`
	Code         int    `json:"code"`
	ErrorMessage string `json:"errorMessage"`
}

const (
	UploadChunkWSBodyStatusCompleted = "completed"
	UploadChunkWSBodyStatusSuccess   = "success"
	UploadChunkWSBodyStatusFailed    = "failed"
)
