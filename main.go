package main

import (
	"go.uber.org/zap"
	"go_client/internal/communicator"
	"go_client/internal/component"
	"go_client/internal/utils"
)

func checkErrors(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	utils.InitLogging()
	communicator.InitClient()
	defer communicator.DestroyClient()

	// connect two components
	comp1, err := component.NewLocalComponent("thekind", "thename1")
	checkErrors(err)
	comp1.Connect()
	zap.L().Info("components list", zap.Any("components", comp1.ListComponents("")))

	// comp 1 watch comp 2 parameter events
	var remoteComp2 *component.RemoteComponent
	comp1.WatchComponents(
		comp1.Kind,
		func(stub *component.RemoteComponentStub) {
			remoteComp2 = stub.Spawn()
			comp1.WatchParameters(
				stub,
				func(key, val string) {
					zap.L().Info("param set", zap.String("key", key), zap.String("val", val))
				},
				func(key string) {
					zap.L().Info("param delete", zap.String("key", key))
				},
			)
		},
		func(stub *component.RemoteComponentStub) {
			zap.L().Info("component deleted", zap.Any("comp", stub))
		},
	)

	comp2, err := component.NewLocalComponent("thekind", "thename1")
	checkErrors(err)
	comp2.Connect()

	// add new comp2 parameter
	comp2.SetParam("foo", "bar")
	// re-read before remote change
	println(">>>")
	println(comp2.GetParam("foo"))
	for remoteComp2 == nil {
	}
	println(remoteComp2.GetParam("foo"))
	// change via remote
	println("===")
	remoteComp2.SetParam("foo", "baz")
	// re-read after remote change
	println(comp2.GetParam("foo"))
	println(remoteComp2.GetParam("foo"))
	println("<<<")

	comp2.DeleteParam("foo")

	// disconnect
	comp2.Disconnect()
	comp1.Disconnect()
}
