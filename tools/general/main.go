package main

import (
	"fmt"
	"os"
)

// T: Type parameter 类型形参
// int|float32|float64 类型约束 Type constraint
// T int|float32|float64 类型形参列表 type parameter list
// Slice[T] 范型类型 general type
type Slice[T int | float32 | float64] []T

// 范型接收参数
func (s Slice[T]) Sum() T {
	var sum T
	for _, value := range s {
		sum += value
	}
	return sum
}

// Slice[int] int是类型实参 Type Arguments
// 传入类型实参实现的实例化 Instantiations
var a Slice[int] = []int{1, 2, 3}

type MyMap[KEY int | string, VALUE float32 | float64] map[KEY]VALUE

type MyStruct[T int | string, S []T] struct {
	Name  string
	Data  T
	Array S
}

type IPrintData[T int | float32 | string] interface {
	Print(data T)
}

type MyChan[T int | string] chan T

// 范型函数
func Add[T int | float32 | float64](a T, b T) T {
	return a + b
}

// 类型约束定义为接口
type Int interface {
	int | int8 | int16 | int32 | int64
}

type Uint interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

// 类型约束定义为接口，允许低层类型
type Int1 interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Interface 变为类型集 type set

// Int 和 Float的交集，下面这个接口实际上是个空集，没有类型能满足这个接口
type A interface {
	Int
	Float
}

// interface {} 变成了一个全集，使用新的any表示
var anyType any

// 没有定义方法的接口interface 称为Basic interface，基本接口，即范型之前的接口
type MyError interface {
	Error() string
}

// 有方法和类型的接口称为一般接口General Interface
// 一般接口类型不能用来定义变量，只能用于泛型的类型约束中。
type ReadWriter interface { // ReadWriter 接口既有方法也有类型
	~string | ~[]rune
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

// 范型接口
type DataProcessor[T any] interface {
	Process(oriData T) (newData T)
	int | ~struct{ Data any }
}

func main() {
	fmt.Println(os.Args)
}
