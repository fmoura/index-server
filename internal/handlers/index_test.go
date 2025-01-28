package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/http"
)

const suiteTimeout = 300 * time.Second

type IndexHandlerSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc

	service     *MockIndexService
	handler     *IndexHandler
	gofrContext *gofr.Context
	request     *MockRequest
}

func TestIndexHandlerSuite(t *testing.T) {
	suite.Run(t, new(IndexHandlerSuite))
}

func (s *IndexHandlerSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), suiteTimeout)
}

func (s *IndexHandlerSuite) TearDownSuite() {
	s.cancel()
}

func (s *IndexHandlerSuite) SetupTest() {

	s.service = newMockIndexService()
	s.handler = NewIndexHandler(s.service)

	s.request = newMockRequest()
	container, _ := container.NewMockContainer(s.T())
	s.gofrContext = &gofr.Context{
		Context:   s.ctx,
		Container: container,
		Request:   s.request,
	}

}

func (s *IndexHandlerSuite) TestItReturnsValidResponse() {
	s.Run("WhenItReceivesWellFormedRequestAndFindsValue", func() {
		response, err := s.handler.HandleGet(s.gofrContext)

		s.Require().Nil(err)
		s.Require().NotNil(response)
		s.Require().IsType(IndexResponse{}, response)
		s.Require().Equal(1, response.(IndexResponse).Index)
		s.Require().Equal(uint64(200), response.(IndexResponse).Value)
	})

	s.Run("WhenItReceivesWellFormedRequestAndDoNotFindValue", func() {

		s.service.Unset("SearchIndex")
		s.service.On("SearchIndex",
			mock.Anything,
		).Return(-1, uint64(0), false)

		response, err := s.handler.HandleGet(s.gofrContext)

		s.Require().Nil(err)
		s.Require().NotNil(response)
		s.Require().IsType(IndexNotFoundResponse{}, response)
		s.Require().Equal(-1, response.(IndexNotFoundResponse).Index)
		s.Require().Equal(uint64(200), response.(IndexNotFoundResponse).Value)
		s.Require().NotEmpty(response.(IndexNotFoundResponse).ErrorMessage)
	})
}

func (s *IndexHandlerSuite) TestItReturnsInvalidParamError() {
	s.Run("WhenItReceivesNanParam", func() {

		s.request.Unset("PathParam")
		s.request.On("PathParam",
			"value",
		).Return("invalid")

		_, err := s.handler.HandleGet(s.gofrContext)

		s.Require().NotNil(err)

		s.Require().IsType(http.ErrorInvalidParam{}, err)
	})

	s.Run("WhenItReceivesNegativeNumberParam", func() {

		s.request.Unset("PathParam")
		s.request.On("PathParam",
			"value",
		).Return("-200")

		_, err := s.handler.HandleGet(s.gofrContext)

		s.Require().NotNil(err)

		s.Require().IsType(http.ErrorInvalidParam{}, err)
	})

	s.Run("WhenITReceivesNegativeNumberParam", func() {

		s.request.Unset("PathParam")
		s.request.On("PathParam",
			"value",
		).Return("-200")

		_, err := s.handler.HandleGet(s.gofrContext)

		s.Require().NotNil(err)

		s.Require().IsType(http.ErrorInvalidParam{}, err)
	})

	s.Run("WhenITReceivesFloatNumberParam", func() {

		s.request.Unset("PathParam")
		s.request.On("PathParam",
			"value",
		).Return("200.0")

		_, err := s.handler.HandleGet(s.gofrContext)

		s.Require().NotNil(err)

		s.Require().IsType(http.ErrorInvalidParam{}, err)
	})
}

type MockIndexService struct {
	mock.Mock
}

func newMockIndexService() *MockIndexService {
	service := &MockIndexService{}

	service.On("SearchIndex",
		mock.Anything,
	).Return(1, uint64(200), true)

	return service
}

// SearchIndex(value uint64) (index int, actualValue uint64, found bool)
func (m *MockIndexService) SearchIndex(value uint64) (index int, actualValue uint64, found bool) {
	args := m.Called(value)
	return args.Get(0).(int), args.Get(1).(uint64), args.Get(2).(bool)
}

func (m *MockIndexService) Unset(methodName string) {
	for _, call := range m.ExpectedCalls {
		if call.Method == methodName {
			call.Unset()
		}
	}
}

type MockRequest struct {
	mock.Mock
}

func newMockRequest() *MockRequest {
	service := &MockRequest{}

	service.On("PathParam",
		mock.Anything,
	).Return("200")

	return service
}

// Context() context.Context
func (m *MockRequest) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// Param(string) string
func (m *MockRequest) Param(name string) string {
	args := m.Called(name)
	return args.Get(0).(string)
}

// PathParam(string) string
func (m *MockRequest) PathParam(name string) string {
	args := m.Called(name)
	return args.Get(0).(string)
}

// Bind(interface{}) error
func (m *MockRequest) Bind(i interface{}) error {
	args := m.Called(i)
	return args.Get(0).(error)
}

// HostName() string
func (m *MockRequest) HostName() string {
	args := m.Called()
	return args.Get(0).(string)
}

// Params(string) []string
func (m *MockRequest) Params(name string) []string {
	args := m.Called(name)
	return args.Get(0).([]string)
}

func (m *MockRequest) Unset(methodName string) {
	for _, call := range m.ExpectedCalls {
		if call.Method == methodName {
			call.Unset()
		}
	}
}
