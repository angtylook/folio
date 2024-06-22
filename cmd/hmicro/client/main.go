package main

import (
	"context"
	"fmt"

	"github.com/angtylook/folio/api/greeter"

	"github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()
	greeterClient := greeter.NewGreeterService("greeter", service.Client())
	rsp, err := greeterClient.Hello(context.TODO(), &greeter.Request{Name: "zzw"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.Greeting)
}
