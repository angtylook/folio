package main

import (
	"context"
	"fmt"
	"hmicro/greeter"

	"github.com/micro/go-micro/v2"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	rsp.Greeting = "Hello" + req.Name
	return nil
}

func main() {
	service := micro.NewService(micro.Name("greeter"))
	service.Init()
	greeter.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
