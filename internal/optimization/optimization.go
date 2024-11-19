package optimization

import (
	"fmt"
	"mars-rover/internal/models"
)

type Optimizer struct{}

func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

// OptimizeRoute метод для оптимизации последовательности команд в последовательность движений
// упрощает множественные последовательности из вперёд-назад и поворотов,
// чтобы марсоход не бегал много раз назад-вперёд или не крутился на месте
func (o *Optimizer) OptimizeRoute(commands string) ([]models.Move, error) {
	if len(commands) == 0 {
		return []models.Move{}, nil
	}

	var state models.MoveType
	moves := make([]models.Move, 0, len(commands))

	turns := 0
	steps := 0

	for _, command := range commands {
		switch command {
		case 'F', 'B':
			steps = move(command, steps)
			if state == models.Rotation && turns%4 != 0 {
				moves = append(moves, models.Move{Type: models.Rotation, Value: turns % 4})
				turns = 0
			}
			state = models.Movement
		case 'R', 'L':
			turns = rotate(command, turns)
			if state == models.Movement && steps != 0 {
				moves = append(moves, models.Move{Type: models.Movement, Value: steps})
				steps = 0
			}
			state = models.Rotation
		default:
			return nil, fmt.Errorf("%w: %c", models.ErrIncorrectSymbol, command)
		}
	}

	if state == models.Rotation && turns%4 == 0 || state == models.Movement && steps == 0 {
		return moves, nil
	}

	switch state {
	case models.Rotation:
		moves = append(moves, models.Move{Type: models.Rotation, Value: turns % 4})
	case models.Movement:
		moves = append(moves, models.Move{Type: models.Movement, Value: steps})
	}

	return moves, nil
}

func move(command rune, count int) int {
	if command == 'F' {
		return count + 1
	}
	return count - 1
}

func rotate(command rune, count int) int {
	if command == 'L' {
		return count + 1
	}
	return count - 1
}
