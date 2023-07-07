package responsebody

import (
	"encoding/json"
	"net/http"
	"playhouse-server/util"
)

type Wrapper struct {
	Writer http.ResponseWriter
}

func (builder *Wrapper) RawBody(body []byte) {
	_, err := builder.Writer.Write(body)
	if err != nil {
		panic(util.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}
}

func (builder *Wrapper) JsonBodyFromMap(body map[string]any) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(util.ResponseErr{
			Code:    http.StatusInternalServerError,
			ErrBody: err,
		})
	}

	builder.Writer.Header().Set("Content-Type", "application/json")
	builder.RawBody(jsonData)
}

func (builder *Wrapper) Status(code int) *Wrapper {
	builder.Writer.WriteHeader(code)
	return builder
}

func (builder *Wrapper) Header(headers map[string]string) *Wrapper {
	for h, v := range headers {
		builder.Writer.Header().Set(h, v)
	}
	return builder
}
