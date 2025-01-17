package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах.
Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// функция для смены директории
func changeDirectory(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Println("Error change dir: ", err)
	}
}

// функция для вывода текущего рабочего каталога
func printWorkingDirectory() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error print working dir: ", err)
	} else {
		fmt.Println(dir)
	}
}

// функция для вывода текста
func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// функция для завершения процесса
func killProcess(pid string) {
	pidNum, err := strconv.Atoi(pid)
	if err != nil {
		fmt.Println("Error parse pid:", err)
		return
	}
	err = syscall.Kill(pidNum, syscall.SIGTERM)
	if err != nil {
		fmt.Println("Error kill process: ", err)
	}
}

// функция для вывода процессов
func listProcesses() {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error exec command ps: ", err)
	}
}

// функция для выполнения внешних команд с помощью fork/exec
func executeCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error exec command:", err)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("shell> ")

		if !scanner.Scan() {
			fmt.Println("\nProgram exit!")
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Fields(input)

		command := args[0]
		args = args[1:]

		// Обработка команд
		switch command {
		case "cd":
			if len(args) > 0 {
				changeDirectory(args[0])
			} else {
				fmt.Println("Error: dir not specified")
			}
		case "pwd":
			printWorkingDirectory()
		case "echo":
			echo(args)
		case "kill":
			if len(args) > 0 {
				killProcess(args[0])
			} else {
				fmt.Println("Error: no process in args")
			}
		case "ps":
			listProcesses()
		default:
			// если команда неизвестна, пытаемся выполнить её
			executeCommand(command, args)
		}

	}
}
