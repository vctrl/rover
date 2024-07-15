package models

import "errors"

var (
	ErrIncorrectSymbol = errors.New("validation error: unexpected input")
)

type Direction string

const (
	North Direction = "N"
	South Direction = "S"
	East  Direction = "E"
	West  Direction = "W"
)

type Coordinates struct {
	X int
	Y int
}

type MoveType string

const (
	Movement MoveType = "Movement"
	Rotation MoveType = "Rotation"
)

type Move struct {
	Type  MoveType
	Value int
}
