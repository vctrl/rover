package app

import (
	"errors"
	"fmt"
	"github.com/eiannone/keyboard"
	"mars-rover/internal/models"
)

type Rover interface {
	PerformRoute(route []models.Move)
	GetCurrentPosition() models.Coordinates
	GetCurrentDirection() models.Direction
	Move(steps int)
	Rotate(steps int)
}

type Optimizer interface {
	OptimizeRoute(commands string) ([]models.Move, error)
}

type App struct {
	Rover     Rover
	Optimizer Optimizer
}

func NewApp(rover Rover, optimizer Optimizer) *App {
	return &App{
		Rover:     rover,
		Optimizer: optimizer,
	}
}

func (a *App) CalculateRoute(commands string) (models.Coordinates, models.Direction, error) {
	route, err := a.Optimizer.OptimizeRoute(commands)
	if err != nil {
		return models.Coordinates{}, "", err
	}

	a.Rover.PerformRoute(route)
	return a.Rover.GetCurrentPosition(), a.Rover.GetCurrentDirection(), nil
}

func (a *App) InteractiveControl(input <-chan string, output chan<- string) error {
	for command := range input {
		switch command {
		case "up":
			a.Rover.Move(1)
		case "down":
			a.Rover.Move(-1)
		case "right":
			a.Rover.Rotate(-1)
		case "left":
			a.Rover.Rotate(1)
		case "exit":
			close(output)
			return nil
		default:
			output <- fmt.Sprintf("Некорректная команда %v, используйте стрелки вверх, вниз, влево, вправо.", command)
			continue
		}

		pos := a.Rover.GetCurrentPosition()
		dir := a.Rover.GetCurrentDirection()
		output <- fmt.Sprintf("Текущие координаты: (%d, %d), направление: %s", pos.X, pos.Y, dir)
	}

	return nil
}

func (a *App) CaptureInput(input chan<- string) error {
	err := keyboard.Open()
	if err != nil {
		return fmt.Errorf("failed to open keyboard: %w", err)
	}
	defer keyboard.Close()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			return fmt.Errorf("failed to get key: %w", err)
		}

		if key == keyboard.KeyCtrlC {
			input <- "exit"
			close(input)
			break
		}

		switch key {
		case keyboard.KeyArrowUp:
			input <- "up"
		case keyboard.KeyArrowDown:
			input <- "down"
		case keyboard.KeyArrowRight:
			input <- "right"
		case keyboard.KeyArrowLeft:
			input <- "left"
		default:
			input <- "invalid"
		}
	}

	return nil
}

func (a *App) HandleCommands(commands string) (models.Coordinates, models.Direction, error) {
	position, direction, err := a.CalculateRoute(commands)
	if err != nil {
		return models.Coordinates{}, "", err
	}
	return position, direction, nil
}

func HandleError(err error) string {
	if errors.Is(err, models.ErrIncorrectSymbol) {
		return fmt.Sprintf("Некорректный путь: %v, путь должен состоять только из символов F, B, R, L", err)
	}
	return fmt.Sprintf("Ошибка: %v", err)
}
