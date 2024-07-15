package app

import (
	"fmt"
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
			a.Rover.Rotate(1)
		case "left":
			a.Rover.Rotate(-1)
		default:
			output <- "Invalid command, use: up, down, left, right."
			continue
		}

		pos := a.Rover.GetCurrentPosition()
		dir := a.Rover.GetCurrentDirection()
		output <- fmt.Sprintf("Текущие координаты: (%d, %d), направление: %s", pos.X, pos.Y, dir)
	}

	return nil
}
