package main

import (
	"fmt"
	"time"
)

func main() {
	duration := 5
	fmt.Printf("program sleep %vs...\n", duration)
	sleep(time.Duration(duration) * time.Second)
	fmt.Printf("program complete after %vs\n", duration)
}

func sleep(t time.Duration) {
	<-time.After(t)
}
