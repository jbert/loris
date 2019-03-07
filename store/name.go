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

func newFromNameBits(nameBits []string) (Store, error) {

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
		return NewMap(), nil
	case "mutexmap":
		return NewMutexMap(), nil
	case "mutex":
		s, err := newFromNameBits(nameBits)
		if err != nil {
			return nil, fmt.Errorf("Mutex wrapper: %s", err)
		}
		return NewMutex(s), nil
	}
	return nil, fmt.Errorf("Unrecognised name: %s", storeName)
}
