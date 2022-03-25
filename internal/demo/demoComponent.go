package demo

import (
	"go.uber.org/zap"
	"go_client/pkg/component"
)

type Component struct {
	*component.LocalComponent
}

func create(kind, name string) *Component {
	comp, err := component.NewLocalComponent(kind, name)
	if err != nil {
		panic(err)
	}
	comp.Connect()
	return &Component{comp}
}

func (c *Component) log(msg string) {
	zap.L().Debug(msg)
}
