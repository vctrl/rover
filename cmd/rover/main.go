package main

import (
	"bufio"
	"errors"
	"fmt"
	"mars-rover/internal/app"
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
				fmt.Println("Используйте стрелки для управления марсоходом. Нажмите Ctrl+C для выхода.")
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
				a.HandleCommands(commands)
			case ModeFile:
				commands, err := GetCommandsFromFile(filePath)
				if err != nil {
					fmt.Printf("Ошибка получения команд: %v\n", err)
					return
				}
				a.HandleCommands(commands)
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
		err := a.CaptureInput(input)
		if err != nil {
			fmt.Printf("Ошибка ввода: %v\n", err)
		}
	}()

	for msg := range output {
		fmt.Println(msg)
	}

	return nil
}
