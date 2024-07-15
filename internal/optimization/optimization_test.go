package optimization

import (
	"mars-rover/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptimizeRoute(t *testing.T) {
	tests := []struct {
		name          string
		commands      string
		expectedMoves []models.Move
		expectedErr   error
	}{
		{
			name:          "Empty commands",
			commands:      "",
			expectedMoves: []models.Move{},
			expectedErr:   nil,
		},
		{
			name:     "Simple forward move",
			commands: "F",
			expectedMoves: []models.Move{
				{Type: models.Movement, Value: 1},
			},
			expectedErr: nil,
		},
		{
			name:     "Simple backward move",
			commands: "B",
			expectedMoves: []models.Move{
				{Type: models.Movement, Value: -1},
			},
			expectedErr: nil,
		},
		{
			name:     "Simple right turn",
			commands: "R",
			expectedMoves: []models.Move{
				{Type: models.Rotation, Value: -1},
			},
			expectedErr: nil,
		},
		{
			name:     "Simple left turn",
			commands: "L",
			expectedMoves: []models.Move{
				{Type: models.Rotation, Value: 1},
			},
			expectedErr: nil,
		},
		{
			name:     "Mixed commands",
			commands: "FFLRB",
			expectedMoves: []models.Move{
				{Type: models.Movement, Value: 2},
				{Type: models.Movement, Value: -1},
			},
			expectedErr: nil,
		},
		{
			name:        "Invalid command",
			commands:    "X",
			expectedErr: models.ErrIncorrectSymbol,
		},
		{
			name:        "Multiple invalid commands",
			commands:    "FFLXBR",
			expectedErr: models.ErrIncorrectSymbol,
		},
		{
			name:     "Multiple valid and invalid commands",
			commands: "FFLXRBR",
			expectedMoves: []models.Move{
				{Type: models.Movement, Value: 2},
				{Type: models.Rotation, Value: 1},
			},
			expectedErr: models.ErrIncorrectSymbol,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewOptimizer()
			moves, err := optimizer.OptimizeRoute(tt.commands)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedMoves, moves)
			}
		})
	}
}

func TestMove(t *testing.T) {
	tests := []struct {
		name     string
		command  rune
		count    int
		expected int
	}{
		{
			name:     "Move forward",
			command:  'F',
			count:    0,
			expected: 1,
		},
		{
			name:     "Move backward",
			command:  'B',
			count:    0,
			expected: -1,
		},
		{
			name:     "Move forward increment",
			command:  'F',
			count:    2,
			expected: 3,
		},
		{
			name:     "Move backward decrement",
			command:  'B',
			count:    -2,
			expected: -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := move(tt.command, tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		name     string
		command  rune
		count    int
		expected int
	}{
		{
			name:     "Rotate left",
			command:  'L',
			count:    0,
			expected: 1,
		},
		{
			name:     "Rotate right",
			command:  'R',
			count:    0,
			expected: -1,
		},
		{
			name:     "Rotate left increment",
			command:  'L',
			count:    2,
			expected: 3,
		},
		{
			name:     "Rotate right decrement",
			command:  'R',
			count:    -2,
			expected: -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rotate(tt.command, tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}
