package service

import (
	"slices"

	"gofr.dev/pkg/gofr/logging"
)

type DataProvider interface {
	Input() []uint64
}

type IndexService struct {
	dataProvider  DataProvider
	conformFactor float64
	logger        logging.Logger
}

func NewIndexService(logger logging.Logger, dataProvider DataProvider, conformValue uint64) *IndexService {
	return &IndexService{
		dataProvider:  dataProvider,
		conformFactor: float64(conformValue) / 100.0,
		logger:        logger,
	}
}

func (s *IndexService) SearchIndex(value uint64) (index int, actualValue uint64, found bool) {
	s.logger.Debugf("Searching for value %d", value)
	slice := s.dataProvider.Input()

	index, found = slices.BinarySearch(
		slice,
		value,
	)

	if found {
		actualValue = slice[index]
		s.logger.Debugf("Found value %d at index %d", actualValue, index)
		return index, actualValue, true
	}

	threshold := float64(value) * (s.conformFactor)
	// Get value from i
	if index < len(slice) && index > 0 {
		actualValue = slice[index]
		if float64(actualValue-value) < threshold {
			s.logger.Debugf("Found value %d at index %d", actualValue, index)
			return index, actualValue, true
		}
	}

	// get value from i-1
	index = index - 1
	if index < len(slice) && index > 0 {
		actualValue = slice[index]
		if float64(value-actualValue) < threshold {
			s.logger.Debugf("Found value %d at index %d", actualValue, index)
			return index, actualValue, true
		}
	}

	s.logger.Debugf("Value %d not found", value)
	return -1, 0, false
}
