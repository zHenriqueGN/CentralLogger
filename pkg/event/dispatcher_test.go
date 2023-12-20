package event

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventTest struct {
	Name    string
	Payload interface{}
}

func (e *EventTest) GetName() string {
	return e.Name
}

func (e *EventTest) GetPayload() interface{} {
	return e.Payload
}

func (e *EventTest) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *EventTest) GetDateTime() time.Time {
	return time.Now()
}

type HandlerTest struct {
	ID int
}

func (h *HandlerTest) Handle(event EventInterface, wg *sync.WaitGroup) {}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

type DispatcherTestSuite struct {
	suite.Suite
	event1     EventInterface
	event2     EventInterface
	handler1   HandlerInterface
	handler2   HandlerInterface
	dispatcher *Dispatcher
}

func TestDispatcherSuite(t *testing.T) {
	suite.Run(t, new(DispatcherTestSuite))

}

func (suite *DispatcherTestSuite) SetupTest() {
	suite.event1 = &EventTest{"event 1", "payload 1"}
	suite.event2 = &EventTest{"event 2", "payload 2"}
	suite.handler1 = &HandlerTest{ID: 1}
	suite.handler2 = &HandlerTest{ID: 2}
	suite.dispatcher = NewDispatcher()
}

func (suite *DispatcherTestSuite) TestGivenEqualHandlers_WhenRegisteringHandlers_ShouldReturnError() {
	err := suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.Nil(err)
	suite.Len(suite.dispatcher.handlers[suite.event1.GetName()], 1)
	err = suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.ErrorIs(err, ErrHandlerAlreadyRegistered)
}

func (suite *DispatcherTestSuite) TestGivenEventsAndItsHandlers_WhenRegisteringHandlers_ShouldRegisterWithoutErrors() {
	err := suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.Nil(err)
	suite.Len(suite.dispatcher.handlers[suite.event1.GetName()], 1)
	err = suite.dispatcher.Register(suite.event1, suite.handler2)
	suite.Nil(err)
	suite.Len(suite.dispatcher.handlers[suite.event1.GetName()], 2)
	suite.Equal(suite.dispatcher.handlers[suite.event1.GetName()][0], suite.handler1)
	suite.Equal(suite.dispatcher.handlers[suite.event1.GetName()][1], suite.handler2)

	err = suite.dispatcher.Register(suite.event2, suite.handler1)
	suite.Nil(err)
	suite.Len(suite.dispatcher.handlers[suite.event2.GetName()], 1)
	err = suite.dispatcher.Register(suite.event2, suite.handler2)
	suite.Nil(err)
	suite.Len(suite.dispatcher.handlers[suite.event2.GetName()], 2)
	suite.Equal(suite.dispatcher.handlers[suite.event2.GetName()][0], suite.handler1)
	suite.Equal(suite.dispatcher.handlers[suite.event2.GetName()][1], suite.handler2)
}

func (suite *DispatcherTestSuite) TestGivenAnEventNotRegistered_WhenRemovingAHandler_ShouldReturnError() {
	err := suite.dispatcher.Remove(suite.event1, suite.handler1)
	suite.ErrorIs(err, ErrEventNotRegistered)
}

func (suite *DispatcherTestSuite) TestGivenAnHandlerNotRegistered_WhenRemovingAHandler_ShouldReturnError() {
	err := suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.Nil(err)
	err = suite.dispatcher.Remove(suite.event1, suite.handler2)
	suite.ErrorIs(err, ErrHandlerNotFound)
}

func (suite *DispatcherTestSuite) TestGivenExistentsEventAndHandler_WhenRemovingAHandler_ShouldReturnNil() {
	err := suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.Nil(err)
	err = suite.dispatcher.Register(suite.event1, suite.handler2)
	suite.Nil(err)
	err = suite.dispatcher.Remove(suite.event1, suite.handler1)
	suite.Nil(err)
	suite.Len(suite.dispatcher.handlers[suite.event1.GetName()], 1)
}

func (suite *DispatcherTestSuite) TestGivenAnEventNotRegistered_WhenUsingHas_ShouldReturnError() {
	has, err := suite.dispatcher.Has(suite.event1, suite.handler1)
	suite.ErrorIs(err, ErrEventNotRegistered)
	suite.False(has)
}

func (suite *DispatcherTestSuite) TestGivenAHandlerNotRegistered_WhenUsingHas_ShouldReturnFalse() {
	err := suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.Nil(err)
	has, err := suite.dispatcher.Has(suite.event1, suite.handler2)
	suite.Nil(err)
	suite.False(has)
}

func (suite *DispatcherTestSuite) TestGivenAHandlerNotRegistered_WhenUsingHas_ShouldReturnTrue() {
	err := suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.Nil(err)
	has, err := suite.dispatcher.Has(suite.event1, suite.handler1)
	suite.Nil(err)
	suite.True(has)
}

func (suite *DispatcherTestSuite) TestGivenHandlersRegistered_WhenUsingClear_ShouldClearAllHandlers() {
	err := suite.dispatcher.Register(suite.event1, suite.handler1)
	suite.Nil(err)
	err = suite.dispatcher.Register(suite.event2, suite.handler2)
	suite.Nil(err)
	suite.Len(suite.dispatcher.handlers, 2)
	suite.dispatcher.Clear()
	suite.Len(suite.dispatcher.handlers, 0)
}

func (suite *DispatcherTestSuite) TestGivenANotRegisteredEvent_WhenUsingDispatch_ShoulReturnError() {
	err := suite.dispatcher.Dispatch(suite.event1)
	suite.ErrorIs(err, ErrEventNotRegistered)
}

func (suite *DispatcherTestSuite) TestGivenARegisteredEvent_WhenUsingDispatch_ShoulDispatchTheEvent() {
	handler := &MockHandler{}
	handler.On("Handle", suite.event1)
	err := suite.dispatcher.Register(suite.event1, handler)
	suite.Nil(err)
	err = suite.dispatcher.Dispatch(suite.event1)
	suite.Nil(err)
	handler.AssertExpectations(suite.T())
	handler.AssertNumberOfCalls(suite.T(), "Handle", 1)
}
