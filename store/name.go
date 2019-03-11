package store

import (
	"fmt"
	"strings"
)

func ValidateName(storeName string) error {
	// Worth optimising?
	_, err := NewFromName(storeName)
	return err
}

func NewFromName(storeName string) (Store, error) {
	bits := strings.Split(storeName, ":")
	return newFromNameBits(bits)
}

func ctorFromNameBits(nameBits []string) (func() Store, error) {
	if len(nameBits) < 1 {
		return nil, fmt.Errorf("Empty store name bits")
	}

	nameBit := nameBits[0]
	nameBits = nameBits[1:]

	nameArgs := strings.Split(nameBit, ",")
	storeName := nameArgs[0]
	nameArgs = nameArgs[1:]

	switch storeName {
	case "map":
		return func() Store { return NewMap() }, nil
	case "mutexmap":
		return func() Store { return NewMutexMap() }, nil
	case "mutex":
		ctor, err := ctorFromNameBits(nameBits)
		if err != nil {
			return nil, fmt.Errorf("Mutex wrapper: %s", err)
		}
		return func() Store { return NewMutex(ctor) }, nil
	case "shard":
		ctor, err := ctorFromNameBits(nameBits)
		if err != nil {
			return nil, fmt.Errorf("Sharded store: %s", err)
		}
		return func() Store { return NewSharded(ctor) }, nil
	}
	return nil, fmt.Errorf("Unrecognised name: %s", storeName)
}

func newFromNameBits(nameBits []string) (Store, error) {
	ctor, err := ctorFromNameBits(nameBits)
	if err != nil {
		return nil, err
	}
	return ctor(), nil
}
