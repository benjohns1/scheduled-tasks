package core

// Equality defines the Equals function across the core domain
type Equality interface {
	Equals(other interface{}) bool
}
