package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func main() {
	// парсим аргументы командной строки

	go ServerStart()

	timeoutFlag := flag.String("timeout", "10s", "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Error args < 2")
		os.Exit(1)
	}
	host := args[0]
	port := args[1]

	// парсим таймаут
	timeout, err := time.ParseDuration(*timeoutFlag)
	if err != nil {
		fmt.Printf("Error timeout: %v\n", err)
		os.Exit(1)
	}

	// устанавливаем соединение
	url := host + ":" + port
	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		fmt.Printf("Error of connect to %s: %v\n", url, err)
		os.Exit(1)
	}
	fmt.Printf("connected to %s\n", url)

	defer conn.Close()

	// канал ожидания завершения выполнения работы горутин
	done := make(chan struct{})

	// чтение данных из сокета и запись в STDOUT
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Printf("error of connection: %v\n", err)
		}
		close(done)
	}()

	// чтение данных из STDIN и запись в сокет
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			fmt.Printf("error writing to connection: %v\n", err)
		}
		conn.Close()
		close(done)
	}()

	<-done
	fmt.Println("Disconnected. Program complete.")
}
