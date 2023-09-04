package request

import (
	"encoding/json"
	"io"
	"net/http"
	"playhouse-server/response"
)

type RequestBody interface {
	UploadRegistrationBody | UploadChunkBody | AuthRegistrationBody
}

type UnCheckedRequestBody interface {
	isValid() bool
}

func ToRequestBody[T RequestBody](b *T, r *http.Request) {
	bodyStreamer := r.Body
	defer func(bodyStreamer io.ReadCloser) {
		err := bodyStreamer.Close()
		if err != nil {
			panic(response.Error{
				Code:  http.StatusInternalServerError,
				Cause: err,
			})
		}
	}(bodyStreamer)

	decoder := json.NewDecoder(bodyStreamer)

	if err := decoder.Decode(b); err != nil {
		panic(response.Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
	}
}

func IsNotValid(b UnCheckedRequestBody) bool {
	return !b.isValid()
}
