package response

import (
	"encoding/json"
	"net/http"
)

type Builder struct {
	Writer http.ResponseWriter
}

func (builder *Builder) Status(code int) *Builder {
	builder.Writer.WriteHeader(code)
	return builder
}

func (builder *Builder) Header(headers map[string]string) *Builder {
	for h, v := range headers {
		builder.Writer.Header().Set(h, v)
	}
	return builder
}

func (builder *Builder) BuildWithBytes(body []byte) {
	_, err := builder.Writer.Write(body)
	if err != nil {
		panic(Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
	}
}

func (builder *Builder) BuildWithJson(body map[string]any) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
	}

	// without overriding other headers
	builder.Writer.Header().Set("Content-Type", "application/json")
	builder.BuildWithBytes(jsonData)
}

func (builder *Builder) BuildWithJsonList(body []map[string]any) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
	}

	// without overriding other headers
	builder.Writer.Header().Set("Content-Type", "application/json")
	builder.BuildWithBytes(jsonData)
}
