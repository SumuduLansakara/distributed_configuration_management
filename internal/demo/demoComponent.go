package demo

import (
	"fmt"
	"go_client/pkg/component"
)

type DemoComponent struct {
	*component.LocalComponent
}

func create(kind, name string) *DemoComponent {
	comp, err := component.NewLocalComponent(kind, name)
	if err != nil {
		panic(err)
	}
	comp.Connect()
	return &DemoComponent{comp}
}

func (c *DemoComponent) log(msg string) {
	fmt.Printf("%s\n", msg)
}
