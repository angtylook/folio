package main

import (
	"fmt"

	"base"
)

func main() {
	e1 := base.NewExecutor(10)
	e2 := base.NewExecutor(10)
	go e1.Start()
	go e2.Start()

	e1.Dispatch(func() {
		fmt.Println("Dispatch on e1")
	})
	fmt.Println("after dispatch e1")

	e2.Dispatch(func() {
		fmt.Println("Dispatch on e2")
	})
	fmt.Println("after dispatch e2")

	for i := 0; i < 10; i++ {
		v := i
		e1.Post(func() {
			fmt.Printf("post on e1 ==> %d \n", v)
		})

		e2.Post(func() {
			fmt.Printf("post on e2 ==> %d \n", v)
		})
	}
	fmt.Println("after post")

	e1.Stop(base.StopOptDiscard)
	e2.Stop(base.StopOptDiscard)

	e1.WaitForStop()
	e2.WaitForStop()
}
