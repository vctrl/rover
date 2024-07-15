package app

import (
	"errors"
	"mars-rover/internal/mocks"
	"mars-rover/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRover := mocks.NewMockRover(ctrl)
	mockOptimizer := mocks.NewMockOptimizer(ctrl)
	app := NewApp(mockRover, mockOptimizer)

	tests := []struct {
		name              string
		commands          string
		expectedRoute     []models.Move
		expectedPosition  models.Coordinates
		expectedDirection models.Direction
		optimizeError     error
		expectError       bool
	}{
		{
			name:     "Successful route",
			commands: "FFLRB",
			expectedRoute: []models.Move{
				{Type: models.Movement, Value: 2},
				{Type: models.Rotation, Value: 0},
				{Type: models.Movement, Value: -1},
			},
			expectedPosition:  models.Coordinates{X: 1, Y: 2},
			expectedDirection: models.North,
			optimizeError:     nil,
			expectError:       false,
		},
		{
			name:              "Optimize error",
			commands:          "FFLRB",
			expectedRoute:     nil,
			expectedPosition:  models.Coordinates{},
			expectedDirection: "",
			optimizeError:     errors.New("optimization error"),
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOptimizer.EXPECT().OptimizeRoute(tt.commands).Return(tt.expectedRoute, tt.optimizeError)

			if tt.optimizeError == nil {
				mockRover.EXPECT().PerformRoute(tt.expectedRoute)
				mockRover.EXPECT().GetCurrentPosition().Return(tt.expectedPosition)
				mockRover.EXPECT().GetCurrentDirection().Return(tt.expectedDirection)
			}

			position, direction, err := app.CalculateRoute(tt.commands)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedPosition, position)
				assert.Equal(t, tt.expectedDirection, direction)
			}
		})
	}
}

func TestInteractiveControl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRover := mocks.NewMockRover(ctrl)
	app := NewApp(mockRover, nil)

	tests := []struct {
		name            string
		commands        []string
		expectedOutputs []string
	}{
		{
			name:     "Valid commands",
			commands: []string{"up", "down", "right", "left"},
			expectedOutputs: []string{
				"Текущие координаты: (1, 1), направление: N",
				"Текущие координаты: (1, 1), направление: N",
				"Текущие координаты: (1, 1), направление: N",
				"Текущие координаты: (1, 1), направление: N",
			},
		},
		{
			name:     "Invalid command",
			commands: []string{"up", "invalid", "left"},
			expectedOutputs: []string{
				"Текущие координаты: (1, 1), направление: N",
				"Некорректная команда invalid, используйте стрелки вверх, вниз, влево, вправо.",
				"Текущие координаты: (1, 1), направление: N",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make(chan string)
			output := make(chan string)

			go func() {
				defer close(output)
				err := app.InteractiveControl(input, output)
				require.NoError(t, err)
			}()

			mockRover.EXPECT().Move(1).AnyTimes()
			mockRover.EXPECT().Move(-1).AnyTimes()
			mockRover.EXPECT().Rotate(1).AnyTimes()
			mockRover.EXPECT().Rotate(-1).AnyTimes()
			mockRover.EXPECT().GetCurrentPosition().AnyTimes().Return(models.Coordinates{X: 1, Y: 1})
			mockRover.EXPECT().GetCurrentDirection().AnyTimes().Return(models.North)

			go func() {
				for _, command := range tt.commands {
					input <- command
				}
				close(input)
			}()

			for _, expected := range tt.expectedOutputs {
				assert.Equal(t, expected, <-output)
			}
		})
	}
}
