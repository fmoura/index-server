package handlers

import (
	"strconv"
	"strings"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

type IndexResponse struct {
	Index        int64  `json:"index"`
	Value        int    `json:"value"`
	ErrorMessage string `json:"error message"`
}

func HandleIndex(ctx *gofr.Context) (interface{}, error) {

	indexParam := ctx.Request.PathParam("index")

	if strings.Trim(indexParam, " ") == "" {
		return nil, http.ErrorMissingParam{Params: []string{"index"}}
	}

	index, err := strconv.ParseInt(indexParam, 10, 64)

	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"index"}}
	}

	return IndexResponse{
		Index:        index,
		Value:        -1,
		ErrorMessage: "Not Implemented yet",
	}, nil
}
