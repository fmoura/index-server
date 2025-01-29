package handler

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

const (
	valuePathParamName = "value"
	IndexValuePath     = "/index/{value}"
	maxValue           = 1000000
)

type IndexService interface {
	SearchIndex(value uint64) (index int, actualValue uint64, found bool)
}

type IndexResponse struct {
	Index int    `json:"index"`
	Value uint64 `json:"value"`
}

type IndexNotFoundResponse struct {
	Index        int    `json:"index"`
	Value        uint64 `json:"value"`
	ErrorMessage string `json:"error message"`
}

type IndexHandler struct {
	dataService IndexService
}

func NewIndexHandler(dataService IndexService) *IndexHandler {
	return &IndexHandler{dataService: dataService}
}

func (h *IndexHandler) HandleGet(ctx *gofr.Context) (interface{}, error) {

	value, err := strconv.ParseUint(ctx.PathParam(valuePathParamName), 10, 64)

	if err != nil || value > maxValue {
		return nil, http.ErrorInvalidParam{Params: []string{valuePathParamName}}
	}

	index, actualValue, found := h.dataService.SearchIndex(value)

	if !found {
		return IndexNotFoundResponse{
			Index:        -1,
			Value:        value,
			ErrorMessage: "Value not found",
		}, nil
	}

	return IndexResponse{
		Index: index,
		Value: actualValue,
	}, nil
}
