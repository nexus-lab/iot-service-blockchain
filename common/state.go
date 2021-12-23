/*
 * Common ledger asset objects and utility functions.
 */
package common

// StateInterface common ledger state interface
type StateInterface interface {
	// GetKeyComponents return components that compose the state key
	GetKeyComponents() []string

	// Serialize transform current state object to JSON string
	Serialize() ([]byte, error)

	// Validate check if the state properties are valid
	Validate() error
}
