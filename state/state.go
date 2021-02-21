package state

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Store
type Store interface{}

// State
type State struct {
	Gtx    layout.Context
	states map[string]Store
	Th     *material.Theme
}

// NewState creates a new empty State
func NewState() *State {
	return &State{
		states: make(map[string]Store),
	}
}

// Get returns the value of the store by its key and the existence of the desired item
func (s State) Get(key string) (Store, bool) {
	v, ok := s.states[key]
	return v, ok
}

// Get returns the value of the store by its key
func (s State) MustGet(key string) Store {
	return s.states[key]
}

// Set registers a new Store
func (s State) Set(key string, v Store) {
	s.states[key] = v
}
