package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("begin call testDefer")
	testDefer()
	fmt.Println("end call testDefer")
}

func testDefer() {
	defer func() {
		fmt.Println("call defer")
		time.Sleep(time.Second * 3)
		//do nothing
	}()
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		fmt.Println("in testDefer")
	}
	return
}
