package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Direction string

const (
	North Direction = "N"
	South Direction = "S"
	East  Direction = "E"
	West  Direction = "W"
)

type Coordinates struct {
	x int
	y int
}

type Rover struct {
	direction Direction
	pos       Coordinates
}

func NewRover() Rover {
	return Rover{
		direction: North,
		pos:       Coordinates{1, 1},
	}
}

func (r *Rover) PerformRoute(route []Move) {
	for _, action := range route {
		switch action.Type {
		case Movement:
			r.move(action.Value)
		case Rotation:
			r.rotate(action.Value)
		}
	}
}

func (r *Rover) GetCurrentPosition() Coordinates {
	return r.pos
}

func (r *Rover) move(steps int) {
	switch r.direction {
	case North:
		r.pos.y += steps
	case South:
		r.pos.y -= steps
	case West:
		r.pos.x -= steps
	case East:
		r.pos.x += steps
	}
}

func (r *Rover) rotate(steps int) {
	directions := []Direction{North, East, South, West}
	currentIndex := indexOf(r.direction, directions)
	newIndex := (currentIndex + steps) % len(directions)
	if newIndex < 0 {
		newIndex += len(directions)
	}

	r.direction = directions[newIndex]
}

func indexOf(dir Direction, directions []Direction) int {
	for i, d := range directions {
		if d == dir {
			return i
		}
	}
	return -1
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

var (
	ErrIncorrectSymbol = errors.New("validation error: unexpected input")
)

func OptimizeRoute(commands string) ([]Move, error) {
	if len(commands) == 0 {
		return []Move{}, nil
	}

	var state MoveType
	moves := make([]Move, 0, len(commands))

	turns := 0
	steps := 0

	for _, command := range commands {
		switch command {
		case 'F', 'B':
			steps = move(command, steps)
			if state == Rotation {
				moves = append(moves, Move{Rotation, turns % 4})
				turns = 0
			}
			state = Movement
		case 'R', 'L':
			turns = rotate(command, turns)
			if state == Movement {
				moves = append(moves, Move{Movement, steps})
				steps = 0
			}
			state = Rotation
		default:
			return nil, fmt.Errorf("%w: %c", ErrIncorrectSymbol, command)
		}
	}

	switch state {
	case Rotation:
		moves = append(moves, Move{Rotation, turns % 4})
	case Movement:
		moves = append(moves, Move{Movement, steps})
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

func interactiveControl(rover *Rover) error {
	fmt.Println("Используйте стрелки для управления марсоходом. Нажмите ESC для выхода.")
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

		if key == keyboard.KeyEsc {
			break
		}

		switch key {
		case keyboard.KeyArrowUp:
			rover.move(1)
		case keyboard.KeyArrowDown:
			rover.move(-1)
		case keyboard.KeyArrowRight:
			rover.rotate(1)
		case keyboard.KeyArrowLeft:
			rover.rotate(-1)
		default:
			fmt.Println("Неверная команда, используйте стрелки для управления.")
		}

		pos := rover.GetCurrentPosition()
		fmt.Printf("Текущие координаты: (%d, %d), направление: %s\n", pos.x, pos.y, rover.direction)
	}

	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rover",
		Short: "Марсоход",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Добро пожаловать в центр управления марсоходом 'Curiosity'!")

			prompt := promptui.Select{
				Label: "Выберите режим",
				Items: []string{"Ввести маршрут с консоли", "Загрузить маршрут из файла", "Интерактивное управление стрелками"},
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Ошибка выбора: %v\n", err)
				return
			}

			var commands string
			rover := NewRover()

			switch result {
			case "Ввести маршрут с консоли":
				fmt.Print("Введите маршрут: ")
				reader := bufio.NewReader(os.Stdin)
				commands, _ = reader.ReadString('\n')
				commands = strings.TrimSpace(commands)
			case "Загрузить маршрут из файла":
				fmt.Print("Введите путь к файлу: ")
				var filePath string
				fmt.Scan(&filePath)
				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Ошибка чтения файла: %v\n", err)
					return
				}
				commands = string(content)
				commands = strings.TrimSpace(commands)
			case "Интерактивное управление стрелками":
				if err := interactiveControl(&rover); err != nil {
					fmt.Printf("Ошибка в интерактивном режиме: %v\n", err)
				}
				return
			}

			optimizedRoute, err := OptimizeRoute(commands)
			if err != nil {
				if errors.Is(err, ErrIncorrectSymbol) {
					fmt.Printf("Некорректный путь: %v, путь должен состоять только из символов F, B, R, L\n", err)
				} else {
					fmt.Printf("Ошибка оптимизации маршрута: %v\n", err)
				}
				return
			}

			rover.PerformRoute(optimizedRoute)

			finalPos := rover.GetCurrentPosition()
			fmt.Printf("Расчёт выполнен успешно. Конечное положение Марсохода: (%d, %d), направление: %s\n",
				finalPos.x, finalPos.y, rover.direction)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Ошибка выполнения команды: %v\n", err)
	}
}
