package main

import (
	"fmt"
	"time"
)

func main() {
	after := time.Date(2020, 07, 13, 0, 0, 0, 0, time.Local).AddDate(0, 0, 42)
	fmt.Println(after)
}
