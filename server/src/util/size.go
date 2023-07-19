package util

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
)

func MaxChunkSize(videoSize int) int {
	var maxChunkSize int
	switch {
	case videoSize < 1*MB:
		maxChunkSize = 1
	case videoSize >= 1*MB:
		maxChunkSize = MB
	}

	return maxChunkSize
}
