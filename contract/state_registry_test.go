package contract

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockState struct {
	Id    string
	Value int
}

func (s *mockState) GetKeyComponents() []string {
	return []string{s.Id}
}

func (s *mockState) Serialize() ([]byte, error) {
	if s.Value == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *mockState) Validate() error {
	if s.Id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	return nil
}

type MockStateRegistry struct {
	mock.Mock
}

func (r *MockStateRegistry) PutState(state common.StateInterface) error {
	args := r.Called(state)
	return args.Error(0)
}

func (r *MockStateRegistry) GetState(keyComponents ...string) (common.StateInterface, error) {
	args := r.Called(keyComponents)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(common.StateInterface), args.Error(1)
}

func (r *MockStateRegistry) GetStates(keyComponents ...string) ([]common.StateInterface, error) {
	args := r.Called(keyComponents)
	return args.Get(0).([]common.StateInterface), args.Error(1)
}

func (r *MockStateRegistry) RemoveState(state common.StateInterface) error {
	args := r.Called(state)
	return args.Error(0)
}

type StateRegistryTestSuite struct {
	suite.Suite
	stub     *shimtest.MockStub
	registry *StateRegistry
}

func (s *StateRegistryTestSuite) SetupTest() {
	s.stub = shimtest.NewMockStub("StateRegistryTest", nil)
	identity, _ := cid.New(s.stub)

	ctx := &TransactionContext{}
	ctx.SetStub(s.stub)
	ctx.SetClientIdentity(identity)

	s.registry = &StateRegistry{
		ctx:  ctx,
		Name: "states",
		Deserialize: func(data []byte) (common.StateInterface, error) {
			st := new(mockState)

			if err := json.Unmarshal(data, st); err != nil {
				return nil, err
			}

			return st, nil
		},
	}
}

func (s *StateRegistryTestSuite) TestPutState() {
	state := mockState{Id: "", Value: 0}

	s.stub.MockTransactionStart("PutState")
	err := s.registry.PutState(&state)
	s.stub.MockTransactionEnd("PutState")
	assert.Error(s.T(), err, "should validate state before put state into ledger")
	state.Id = "123456"

	s.stub.MockTransactionStart("PutState")
	err = s.registry.PutState(&state)
	s.stub.MockTransactionEnd("PutState")
	assert.Error(s.T(), err, "should serialize to non-nil value")
	state.Value = 1

	s.stub.MockTransactionStart("PutState")
	err = s.registry.PutState(&state)
	s.stub.MockTransactionEnd("PutState")
	assert.Nil(s.T(), err, "should put state into ledger without error")

	key, _ := s.stub.CreateCompositeKey(s.registry.Name, state.GetKeyComponents())
	data, _ := s.stub.GetState(key)
	assert.Equal(s.T(), "{\"Id\":\"123456\",\"Value\":1}", string(data), "should put state into ledger")
}

func (s *StateRegistryTestSuite) TestGetState() {
	key, _ := s.stub.CreateCompositeKey(s.registry.Name, []string{"123456"})
	s.stub.MockTransactionStart("GetState")
	s.stub.PutState(key, []byte("{\"Id\":\"123456\",\"Value\":1}"))
	s.stub.MockTransactionEnd("GetState")

	state, err := s.registry.GetState("1")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
	assert.Nil(s.T(), state, "should return empty state")

	state, err = s.registry.GetState("123456")
	assert.Nil(s.T(), err, "should get state from ledger without error")
	assert.Equal(s.T(), "123456", state.(*mockState).Id, "should get correct state ID from ledger")
	assert.Equal(s.T(), 1, state.(*mockState).Value, "should get correct state value from ledger")
}

func (s *StateRegistryTestSuite) TestGetStates() {
	key1, _ := s.stub.CreateCompositeKey(s.registry.Name, []string{"A", "B"})
	key2, _ := s.stub.CreateCompositeKey(s.registry.Name, []string{"A", "C"})
	key3, _ := s.stub.CreateCompositeKey(s.registry.Name, []string{"A", "D"})
	s.stub.MockTransactionStart("GetStates")
	s.stub.PutState(key1, []byte("{\"Id\":\"1\",\"Value\":1}"))
	s.stub.PutState(key2, []byte("{\"Id\":\"2\",\"Value\":2}"))
	s.stub.PutState(key3, []byte("{\"Id\":\"3\",\"Value\":3}"))
	s.stub.MockTransactionEnd("GetStates")

	states, err := s.registry.GetStates("B")
	assert.Nil(s.T(), err, "should return zero state without error")
	assert.Zero(s.T(), len(states), "should return zero state")

	states, err = s.registry.GetStates("A")
	assert.Nil(s.T(), err, "should get states from ledger without error")
	assert.Equal(s.T(), 3, len(states), "should get correct number of states")
	assert.Equal(s.T(), 1, states[0].(*mockState).Value, "should get correct states")
	assert.Equal(s.T(), 2, states[1].(*mockState).Value, "should get correct states")
	assert.Equal(s.T(), 3, states[2].(*mockState).Value, "should get correct states")
}

func (s *StateRegistryTestSuite) TestRemoveState() {
	key, _ := s.stub.CreateCompositeKey(s.registry.Name, []string{"123456"})
	s.stub.MockTransactionStart("RemoveState")
	s.stub.PutState(key, []byte("{\"Id\":\"123456\",\"Value\":1}"))
	s.stub.MockTransactionEnd("RemoveState")

	s.stub.MockTransactionStart("RemoveState")
	err := s.registry.RemoveState(&mockState{Id: "1", Value: 0})
	s.stub.MockTransactionEnd("RemoveState")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")

	s.stub.MockTransactionStart("RemoveState")
	err = s.registry.RemoveState(&mockState{Id: "123456", Value: 0})
	s.stub.MockTransactionEnd("RemoveState")
	assert.Nil(s.T(), err, "should remove state from ledger without error")
	data, _ := s.stub.GetState(key)
	assert.Nil(s.T(), data, "should remove state form ledger")
}

func TestStateRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(StateRegistryTestSuite))
}
