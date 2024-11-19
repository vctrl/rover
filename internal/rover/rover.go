package rover

import "mars-rover/internal/models"

type Rover struct {
	Direction models.Direction
	Pos       models.Coordinates
}

func NewRover() *Rover {
	return &Rover{
		Direction: models.North,
		Pos:       models.Coordinates{X: 1, Y: 1},
	}
}

func (r *Rover) PerformRoute(route []models.Move) {
	for _, action := range route {
		switch action.Type {
		case models.Movement:
			r.Move(action.Value)
		case models.Rotation:
			r.Rotate(action.Value)
		}
	}
}

func (r *Rover) GetCurrentPosition() models.Coordinates {
	return r.Pos
}

func (r *Rover) GetCurrentDirection() models.Direction {
	return r.Direction
}

func (r *Rover) Move(steps int) {
	switch r.Direction {
	case models.North:
		r.Pos.Y += steps
	case models.South:
		r.Pos.Y -= steps
	case models.West:
		r.Pos.X -= steps
	case models.East:
		r.Pos.X += steps
	}
}

func (r *Rover) Rotate(steps int) {
	directions := []models.Direction{models.North, models.West, models.South, models.East}
	currentIndex := indexOf(r.Direction, directions)
	newIndex := (currentIndex + steps) % len(directions)
	if newIndex < 0 {
		newIndex += len(directions)
	}

	r.Direction = directions[newIndex]
}

// todo move to some common package...
func indexOf(dir models.Direction, directions []models.Direction) int {
	for i, d := range directions {
		if d == dir {
			return i
		}
	}
	return -1
}
