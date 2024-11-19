package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var testFilePath string

func TestMain(m *testing.M) {
	// Setup phase
	testFilePath = "testfile.txt"
	content := []byte("FFLRB\n")
	err := os.WriteFile(testFilePath, content, 0644)
	if err != nil {
		fmt.Printf("Ошибка при создании тестового файла: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	exitVal := m.Run()

	// Teardown phase
	err = os.Remove(testFilePath)
	if err != nil {
		fmt.Printf("Ошибка при удалении тестового файла: %v\n", err)
		os.Exit(1)
	}

	os.Exit(exitVal)
}

func TestMainE2E(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		input          string
		expectedOutput []string
	}{
		{
			name:           "Console mode with valid commands",
			args:           []string{"--mode=console"},
			input:          "FFLRB\n",
			expectedOutput: []string{"Добро пожаловать в центр управления марсоходом 'Curiosity'!", "Введите маршрут:", "Расчёт выполнен успешно. Конечное положение Марсохода: (1, 2), направление: N\n"},
		},
		{
			name:           "File mode with valid commands",
			args:           []string{"--mode=file", "--file=testfile.txt"},
			expectedOutput: []string{"Добро пожаловать в центр управления марсоходом 'Curiosity'!", "Расчёт выполнен успешно. Конечное положение Марсохода: (1, 2), направление: N\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture the output
			var stdout, stderr bytes.Buffer

			// Prepare the command
			cmd := exec.Command("go", "run", ".")
			cmd.Args = append(cmd.Args, tt.args...)
			cmd.Stdin = bytes.NewBufferString(tt.input)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Run the command
			err := cmd.Run()
			if err != nil && tt.name != "Interactive mode" {
				t.Fatalf("Ошибка выполнения команды: %v\nstderr: %v", err, stderr.String())
			}

			output := stdout.String()
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Ожидаемый вывод должен содержать %q, но получили %q", expected, output)
				}
			}
		})
	}
}
