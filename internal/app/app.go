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
		default:
			output <- fmt.Sprintf("Invalid command %v, use: up, down, left, right.", command)
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
	defer func() {
		if err := keyboard.Close(); err != nil {
			fmt.Printf("failed to close keyboard: %v\n", err)
		}
	}()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			return fmt.Errorf("failed to get key: %w", err)
		}

		if key == keyboard.KeyCtrlC {
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
			fmt.Println("Неверная команда, используйте стрелки для управления.")
		}
	}

	return nil
}

func (a *App) HandleCommands(commands string) {
	position, direction, err := a.CalculateRoute(commands)
	if err != nil {
		fmt.Println(HandleError(err))
		return
	}

	fmt.Printf("Расчёт выполнен успешно. Конечное положение Марсохода: (%d, %d), направление: %s\n",
		position.X, position.Y, direction)
}

func HandleError(err error) string {
	if errors.Is(err, models.ErrIncorrectSymbol) {
		return fmt.Sprintf("Некорректный путь: %v, путь должен состоять только из символов F, B, R, L", err)
	}
	return fmt.Sprintf("Ошибка: %v", err)
}
