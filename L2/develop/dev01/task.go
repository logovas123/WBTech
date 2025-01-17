package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	ntpServer := "0.ru.pool.ntp.org" // адрес ntp server

	ntpTime, err := getCurrentTime(ntpServer) // получаем время с сервера
	// если возникает ошибка, то выводим её в Stderr и выходим из программы с ненулевым кодом ошибки
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting current time from ntp server: ", err)
		os.Exit(1)
	}

	// сравниваем время сервера с системным временем
	currentTime := time.Now()

	fmt.Println("Текущее системное время:", currentTime.Format(time.RFC1123))
	fmt.Println("Время ntp сервера:", ntpTime.Format(time.RFC1123))
}

// функция, которая возвращает время для получения
func getCurrentTime(ntpServer string) (time.Time, error) {
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		return time.Time{}, err
	}

	return ntpTime, nil
}
