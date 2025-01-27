package data

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

const suiteTimeout = 300 * time.Second

type TextDataProviderSuite struct {
	suite.Suite
	ctx         context.Context
	cancel      context.CancelFunc
	gofrContext *gofr.Context
}

func TestTextDataProviderSuite(t *testing.T) {
	suite.Run(t, new(TextDataProviderSuite))
}

func (s *TextDataProviderSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), suiteTimeout)

	mockContainer, _ := container.NewMockContainer(s.T())
	s.gofrContext = &gofr.Context{
		Context:   s.ctx,
		Request:   nil,
		Container: mockContainer,
	}
}

func (s *TextDataProviderSuite) TearDownSuite() {
	s.cancel()
}

func (s *TextDataProviderSuite) inputBytesFromSlice(slice []uint64) []byte {
	var buffer bytes.Buffer

	// Iterate over the uint64 slice and write each value as a string with a newline
	for _, num := range slice {
		_, err := buffer.WriteString(strconv.FormatUint(num, 10) + "\n")
		s.Require().Nil(err)
	}

	return buffer.Bytes()
}

func (s *TextDataProviderSuite) TestItReturnsValidSlice() {
	s.Run("WhenInputIsWellFormed", func() {

		data := []uint64{0, 1, 2, 3, 4}

		slice, err := loadInput(s.gofrContext, s.inputBytesFromSlice(data))
		s.Require().Nil(err)
		s.Require().NotNil(slice)

		s.Require().ElementsMatch(data, slice)
	})

	s.Run("WhenInputIsEmpty", func() {

		data := []uint64{}

		slice, err := loadInput(s.gofrContext, s.inputBytesFromSlice(data))
		s.Require().Nil(err)
		s.Require().NotNil(slice)

		s.ElementsMatch(data, slice)
	})
}

func (s *TextDataProviderSuite) TestItReturnsError() {
	s.Run("WhenInputIsMalformed", func() {
		data := []uint64{0, 1, 2, 3}

		var buffer bytes.Buffer

		for _, num := range data {
			_, err := buffer.WriteString(strconv.FormatUint(num, 10) + "\n")
			s.Require().Nil(err)
		}

		_, err := buffer.WriteString("invalid" + "\n")

		s.Require().Nil(err)

		slice, err := loadInput(s.gofrContext, buffer.Bytes())

		s.Require().NotNil(err)
		s.Require().Nil(slice)

	})

	s.Run("WhenInputIsNotSorted", func() {
		data := []uint64{5, 1, 2, 3}

		slice, err := loadInput(s.gofrContext, s.inputBytesFromSlice(data))

		s.Require().NotNil(err)
		s.Require().Nil(slice)

	})

	s.Run("WhenInputHasNegativeNumbers", func() {
		data := []int64{-1, 0, 1, 2, 3}

		var buffer bytes.Buffer

		for _, num := range data {
			_, err := buffer.WriteString(strconv.FormatInt(num, 10) + "\n")
			s.Nil(err)
		}

		slice, err := loadInput(s.gofrContext, buffer.Bytes())

		s.Require().NotNil(err)
		s.Require().Nil(slice)

	})
}
