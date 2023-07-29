package response

// UploadChunkWSBody represents the response to client about the status of uploading a chunk
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
