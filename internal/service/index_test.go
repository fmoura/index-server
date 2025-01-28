package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

const suiteTimeout = 300 * time.Second

type IndexServiceSuite struct {
	suite.Suite
	ctx      context.Context
	cancel   context.CancelFunc
	provider DataProvider
	service  *IndexService
}

func TestIndexServiceSuite(t *testing.T) {
	suite.Run(t, new(IndexServiceSuite))
}

func (s *IndexServiceSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), suiteTimeout)

}

func (s *IndexServiceSuite) SetupTest() {
	s.provider = newMockDataProvider()

	mockContainer, _ := container.NewMockContainer(s.T())
	ctx := &gofr.Context{
		Context:   s.ctx,
		Container: mockContainer,
	}

	s.service = &IndexService{
		dataProvider:  s.provider,
		conformFactor: 0.10,
		logger:        ctx,
	}
}

func (s *IndexServiceSuite) TearDownSuite() {
	s.cancel()
}

func (s *IndexServiceSuite) TestItReturnsValidValue() {
	s.Run("WhenItReceivesExactValue", func() {
		i, value, found := s.service.SearchIndex(100)

		s.True(found)
		s.EqualValues(uint64(100), value)
		s.EqualValues(1, i)
	})

	s.Run("WhenItReceivesConformValue", func() {
		i, value, found := s.service.SearchIndex(110)

		s.True(found)
		s.EqualValues(uint64(100), value)
		s.EqualValues(1, i)
	})
}

func (s *IndexServiceSuite) TestItReturnsNoValue() {
	s.Run("WhenItReceivesNonConformingValue", func() {
		i, value, found := s.service.SearchIndex(150)

		s.False(found)
		s.EqualValues(uint64(0), value)
		s.EqualValues(-1, i)
	})
}

type MockDataProvider struct {
	mock.Mock
}

func newMockDataProvider() *MockDataProvider {
	provider := &MockDataProvider{}

	provider.On("Input").Return([]uint64{0, 100, 200, 300, 400, 500})

	return provider
}

func (m *MockDataProvider) Input() []uint64 {
	args := m.Called()
	return args.Get(0).([]uint64)
}
