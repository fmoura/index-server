package data

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const suiteTimeout = 300 * time.Second

type TextDataProviderSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
}

func TestTextDataProviderSuite(t *testing.T) {
	suite.Run(t, new(TextDataProviderSuite))
}

func (s *TextDataProviderSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), suiteTimeout)

}

func (s *TextDataProviderSuite) inputBytesFromSlice(slice []uint64) []byte {
	var buffer bytes.Buffer

	// Iterate over the uint64 slice and write each value as a string with a newline
	for _, num := range slice {
		_, err := buffer.WriteString(strconv.FormatUint(num, 10) + "\n")
		s.Nil(err)
	}

	return buffer.Bytes()
}

func (s *TextDataProviderSuite) TestItReturnsValidSlice() {
	s.Run("WhenInputIsWellFormed", func() {

		data := []uint64{0, 1, 2, 3, 4}

		slice, err := loadInput(s.inputBytesFromSlice(data))
		s.Nil(err)
		s.NotNil(slice)

		s.ElementsMatch(data, slice)
	})

	s.Run("WhenInputIsEmpty", func() {

		data := []uint64{}

		slice, err := loadInput(s.inputBytesFromSlice(data))
		s.Nil(err)
		s.NotNil(slice)

		s.ElementsMatch(data, slice)
	})
}

func (s *TextDataProviderSuite) TestItReturnsError() {
	s.Run("WhenInputIsMalformed", func() {
		data := []uint64{0, 1, 2, 3}

		var buffer bytes.Buffer

		for _, num := range data {
			_, err := buffer.WriteString(strconv.FormatUint(num, 10) + "\n")
			s.Nil(err)
		}

		_, err := buffer.WriteString("invalid" + "\n")

		s.Nil(err)

		slice, err := loadInput(buffer.Bytes())

		s.NotNil(err)
		s.Nil(slice)

	})

	s.Run("WhenInputIsNotSorted", func() {
		data := []uint64{5, 1, 2, 3}

		slice, err := loadInput(s.inputBytesFromSlice(data))

		s.NotNil(err)
		s.Nil(slice)

	})

	s.Run("WhenInputHasNegativeNumbers", func() {
		data := []int64{-1, 0, 1, 2, 3}

		var buffer bytes.Buffer

		for _, num := range data {
			_, err := buffer.WriteString(strconv.FormatInt(num, 10) + "\n")
			s.Nil(err)
		}

		slice, err := loadInput(buffer.Bytes())

		s.NotNil(err)
		s.Nil(slice)

	})
}
