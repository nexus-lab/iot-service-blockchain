/*
 * Smart contracts and chaincode functions.
 */
package contract

import (
	"fmt"

	"github.com/nexus-lab/iot-service-blockchain/common"
)

// StateRegistryInterface core utilities for managing a list of ledger states
type StateRegistryInterface interface {
	// PutState create or update a state in the ledger
	PutState(state common.StateInterface) error

	// GetState return a state by its key components
	GetState(keyComponents ...string) (common.StateInterface, error)

	// GetStates return a list of states by key components
	GetStates(keyComponents ...string) ([]common.StateInterface, error)

	// RemoveState remove a state from the ledger
	RemoveState(state common.StateInterface) error
}

// StateRegistry default implementations of StateRegistryInterface
type StateRegistry struct {
	ctx TransactionContextInterface

	// Name name of the state list
	Name string

	// Deserialize create a new state instance from its JSON representation
	Deserialize func([]byte) (common.StateInterface, error)
}

// PutState create or update a state in the ledger
func (r *StateRegistry) PutState(state common.StateInterface) error {
	err := state.Validate()
	if err != nil {
		return err
	}

	key, err := r.ctx.GetStub().CreateCompositeKey(r.Name, state.GetKeyComponents())
	if err != nil {
		return err
	}

	data, err := state.Serialize()
	if err != nil {
		return err
	} else if data == nil {
		return fmt.Errorf("serialized state of %T is empty", state)
	}

	return r.ctx.GetStub().PutState(key, data)
}

// GetState return a state by its key
func (r *StateRegistry) GetState(key ...string) (common.StateInterface, error) {
	key_, err := r.ctx.GetStub().CreateCompositeKey(r.Name, key)
	if err != nil {
		return nil, err
	}

	data, err := r.ctx.GetStub().GetState(key_)
	if err != nil {
		return nil, err
	} else if data == nil {
		return nil, &common.NotFoundError{What: key_}
	}

	state, err := r.Deserialize(data)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// GetStates return a list of states by their partial composite key
func (r *StateRegistry) GetStates(key ...string) ([]common.StateInterface, error) {
	iterator, err := r.ctx.GetStub().GetStateByPartialCompositeKey(r.Name, key)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	states := make([]common.StateInterface, 0)
	for iterator.HasNext() {
		result, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		state, err := r.Deserialize(result.Value)
		if err != nil {
			return nil, err
		}

		states = append(states, state)
	}

	return states, nil
}

// RemoveState remove a state from the ledger
func (r *StateRegistry) RemoveState(state common.StateInterface) error {
	key_, err := r.ctx.GetStub().CreateCompositeKey(r.Name, state.GetKeyComponents())
	if err != nil {
		return err
	}

	data, err := r.ctx.GetStub().GetState(key_)
	if err != nil {
		return err
	} else if data == nil {
		return &common.NotFoundError{What: key_}
	}

	return r.ctx.GetStub().DelState(key_)
}
