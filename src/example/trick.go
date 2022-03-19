package main

import "fmt"

/*
var (
    size := 1024
    max_size = size * 2
)
*/

// mixed named and unnamed function parameters
/*
func funcMui(x, y int) (sum int, error) {
    return x + y, nil
}
*/

// 类型不匹配
/*
func fail_with_mistype() {
	list := new([]int) // *[]int  malloc and memset to zero
	// list := make([]int) // slice
	list = append(list, 1)
	fmt.Println(list)
}
*/

type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}
func (p *People) ShowB() {
	fmt.Println("showB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

func main() {
	t := Teacher{}
	t.People.ShowB()
}
