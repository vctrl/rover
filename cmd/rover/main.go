package main

import (
	"bufio"
	"errors"
	"fmt"
	"mars-rover/internal/app"
	"mars-rover/internal/models"
	"mars-rover/internal/optimization"
	"mars-rover/internal/rover"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	ModeInteractive = "interactive"
	ModeConsole     = "console"
	ModeFile        = "file"
)

func main() {
	var (
		mode     string
		filePath string
	)

	var rootCmd = &cobra.Command{
		Use:   "rover",
		Short: "Марсоход",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Добро пожаловать в центр управления марсоходом 'Curiosity'!")

			if mode == "" {
				var err error
				mode, err = SelectMode()
				if err != nil {
					fmt.Printf("Ошибка выбора: %v\n", err)
					return
				}
			}

			r := rover.NewRover()
			optimizer := optimization.NewOptimizer()
			a := app.NewApp(r, optimizer)

			switch mode {
			case ModeInteractive:
				err := HandleInteractiveMode(a)
				if err != nil {
					fmt.Printf("Ошибка в интерактивном режиме: %v\n", err)
				}
			case ModeConsole:
				commands, err := GetCommandsFromConsole()
				if err != nil {
					fmt.Printf("Ошибка получения команд: %v\n", err)
					return
				}
				HandleCommands(a, commands)
			case ModeFile:
				commands, err := GetCommandsFromFile(filePath)
				if err != nil {
					fmt.Printf("Ошибка получения команд: %v\n", err)
					return
				}
				HandleCommands(a, commands)
			default:
				fmt.Println("Неизвестный режим")
			}
		},
	}

	rootCmd.Flags().StringVarP(&mode, "mode", "m", "", "Режим работы (console, file, interactive)")
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "Путь к файлу с командами")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Ошибка выполнения команды: %v\n", err)
	}
}

func HandleError(err error) string {
	if errors.Is(err, models.ErrIncorrectSymbol) {
		return fmt.Sprintf("Некорректный путь: %v, путь должен состоять только из символов F, B, R, L", err)
	}
	return fmt.Sprintf("Ошибка: %v", err)
}

func SelectMode() (string, error) {
	prompt := promptui.Select{
		Label: "Выберите режим",
		Items: []string{"Ввести маршрут с консоли", "Загрузить маршрут из файла", "Интерактивное управление стрелками"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	switch result {
	case "Ввести маршрут с консоли":
		return ModeConsole, nil
	case "Загрузить маршрут из файла":
		return ModeFile, nil
	case "Интерактивное управление стрелками":
		return ModeInteractive, nil
	default:
		return "", errors.New("неверный выбор режима")
	}
}

func GetCommandsFromConsole() (string, error) {
	fmt.Print("Введите маршрут: ")
	reader := bufio.NewReader(os.Stdin)
	commands, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка чтения команды: %v", err)
	}
	return strings.TrimSpace(commands), nil
}

func GetCommandsFromFile(filePath string) (string, error) {
	if filePath == "" {
		fmt.Print("Введите путь к файлу: ")
		fmt.Scan(&filePath)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %v", err)
	}
	return strings.TrimSpace(string(content)), nil
}

func HandleInteractiveMode(a *app.App) error {
	input := make(chan string)
	output := make(chan string)

	go func() {
		err := a.InteractiveControl(input, output)
		if err != nil {
			fmt.Printf("Ошибка в интерактивном режиме: %v\n", err)
		}
	}()

	go func() {
		for {
			var command string
			fmt.Scan(&command)
			if command == "exit" {
				close(input)
				break
			}
			input <- command
		}
	}()

	for msg := range output {
		fmt.Println(msg)
	}
	return nil
}

func HandleCommands(a *app.App, commands string) {
	position, direction, err := a.CalculateRoute(commands)
	if err != nil {
		fmt.Println(HandleError(err))
		return
	}

	fmt.Printf("Расчёт выполнен успешно. Конечное положение Марсохода: (%d, %d), направление: %s\n",
		position.X, position.Y, direction)
}
