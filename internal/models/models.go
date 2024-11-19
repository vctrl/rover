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

// Move структура для описания движения марсохода
type Move struct {
	// Type тип движения
	Type MoveType
	// Value при Type = Movement Value означает количество шагов, при Type = Rotation Value означает количество поворотов на 90 градусов против часовой стрелки
	Value int
}
