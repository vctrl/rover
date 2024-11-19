package rover

import (
	"github.com/stretchr/testify/assert"
	"mars-rover/internal/models"
	"testing"
)

func TestNewRover(t *testing.T) {
	r := NewRover()
	assert.Equal(t, models.North, r.Direction)
	assert.Equal(t, models.Coordinates{X: 1, Y: 1}, r.Pos)
}

func TestRover_Move(t *testing.T) {
	tests := []struct {
		name      string
		initial   models.Coordinates
		direction models.Direction
		steps     int
		expected  models.Coordinates
	}{
		{"Move North", models.Coordinates{X: 1, Y: 1}, models.North, 2, models.Coordinates{X: 1, Y: 3}},
		{"Move South", models.Coordinates{X: 1, Y: 1}, models.South, 2, models.Coordinates{X: 1, Y: -1}},
		{"Move West", models.Coordinates{X: 1, Y: 1}, models.West, 2, models.Coordinates{X: -1, Y: 1}},
		{"Move East", models.Coordinates{X: 1, Y: 1}, models.East, 2, models.Coordinates{X: 3, Y: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rover{
				Direction: tt.direction,
				Pos:       tt.initial,
			}
			r.Move(tt.steps)
			assert.Equal(t, tt.expected, r.Pos)
		})
	}
}

func TestRover_Rotate(t *testing.T) {
	tests := []struct {
		name     string
		initial  models.Direction
		steps    int
		expected models.Direction
	}{
		{"Rotate Right from North", models.North, 1, models.West},
		{"Rotate Left from North", models.North, -1, models.East},
		{"Full rotation to the right", models.North, 4, models.North},
		{"Full rotation to the left", models.North, -4, models.North},
		{"Rotate Right from East", models.East, 1, models.North},
		{"Rotate Left from East", models.East, -1, models.South},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rover{
				Direction: tt.initial,
			}
			r.Rotate(tt.steps)
			assert.Equal(t, tt.expected, r.Direction)
		})
	}
}

func TestRover_PerformRoute(t *testing.T) {
	tests := []struct {
		name        string
		initial     *Rover
		route       []models.Move
		expectedPos models.Coordinates
		expectedDir models.Direction
	}{
		{
			"Simple Move",
			NewRover(),
			[]models.Move{{Type: models.Movement, Value: 3}},
			models.Coordinates{X: 1, Y: 4},
			models.North,
		},
		{
			"Move and Rotate",
			NewRover(),
			[]models.Move{
				{Type: models.Movement, Value: 1},
				{Type: models.Rotation, Value: 1},
				{Type: models.Movement, Value: 1},
			},
			models.Coordinates{X: 0, Y: 2},
			models.West,
		},
		{
			"Complex Route",
			NewRover(),
			[]models.Move{
				{Type: models.Movement, Value: 1}, // (1, 2), N
				{Type: models.Rotation, Value: 1}, // (1, 2), W
				{Type: models.Movement, Value: 1}, // (0, 2), W
				{Type: models.Rotation, Value: 2}, // (0, 2), E
				{Type: models.Movement, Value: 1}, // (1, 2), E
			},
			models.Coordinates{X: 1, Y: 2},
			models.East,
		},
		{
			"Complex Route with Back",
			NewRover(),
			[]models.Move{
				{Type: models.Movement, Value: 1},  // (1, 2), N
				{Type: models.Rotation, Value: 1},  // (1, 2), W
				{Type: models.Movement, Value: 1},  // (0, 2), W
				{Type: models.Rotation, Value: 2},  // (0, 2), E
				{Type: models.Movement, Value: 1},  // (1, 2), E
				{Type: models.Rotation, Value: -1}, // (1, 2), S
				{Type: models.Movement, Value: 1},  // (1, 1), S
				{Type: models.Movement, Value: -1}, // (1, 2), S (Move Back)
			},
			models.Coordinates{X: 1, Y: 2},
			models.South,
		},
		{
			"Route #1 (1, 1) => (-1, 4)",
			NewRover(),
			[]models.Move{
				{Type: models.Movement, Value: 2},  // FF (1,1) => (1, 3)
				{Type: models.Rotation, Value: 1},  // L (N => W)
				{Type: models.Rotation, Value: -2}, // RR (W => E)
				{Type: models.Movement, Value: -2}, // BB (1, 3) => (-1, 3)
				{Type: models.Rotation, Value: 1},  // L (E => N)
				{Type: models.Movement, Value: 1},  // F (-1, 3) => (-1, 4)
				{Type: models.Rotation, Value: -1}, // R (N => E)
				// R
			},
			models.Coordinates{X: -1, Y: 4},
			models.East,
		},
		{
			"Rotate test",
			NewRover(),
			[]models.Move{
				{Type: models.Rotation, Value: 2},  // L (N => S)
				{Type: models.Rotation, Value: -3}, // RRR (S => W)
			},
			models.Coordinates{X: 1, Y: 1},
			models.East,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.PerformRoute(tt.route)
			assert.Equal(t, tt.expectedPos, tt.initial.GetCurrentPosition())
			assert.Equal(t, tt.expectedDir, tt.initial.GetCurrentDirection())
		})
	}
}

func TestIndexOf(t *testing.T) {
	directions := []models.Direction{models.North, models.West, models.South, models.East}
	tests := []struct {
		name     string
		dir      models.Direction
		expected int
	}{
		{"Index of North", models.North, 0},
		{"Index of West", models.West, 1},
		{"Index of South", models.South, 2},
		{"Index of East", models.East, 3},
		{"Index of Invalid", models.Direction("Invalid"), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := indexOf(tt.dir, directions)
			assert.Equal(t, tt.expected, result)
		})
	}
}
