package data

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

//go:embed resources/input.txt
var inputTxtBytes []byte

type TextDataProvider struct {
	inputSlice []uint64
}

func (p *TextDataProvider) Input() []uint64 {
	return p.inputSlice
}

func NewTextDataProvider() (*TextDataProvider, error) {

	slice, err := loadInput(inputTxtBytes)

	if err != nil {
		return nil, err
	}

	return &TextDataProvider{
		inputSlice: slice,
	}, nil
}

func loadInput(input []byte) ([]uint64, error) {
	scanner := bufio.NewScanner(bytes.NewReader(input))

	slice := []uint64{}
	for scanner.Scan() {
		val, err := strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("Error reading input"), err)
		}
		slice = append(slice, val)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Join(fmt.Errorf("Error reading input"), err)
	}

	sorted := sort.SliceIsSorted(slice, func(i, j int) bool {
		return slice[i] <= slice[j]
	})

	if !sorted {
		return nil, fmt.Errorf("Input is not sorted")
	}

	return slice, nil
}
