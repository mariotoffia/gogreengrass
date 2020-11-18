package main

import "github.com/mariotoffia/gogreengrass/sdkc"

func main() {
	sdkc.Start(sdkc.RuntimeOptionSingleThread)
	sdkc.Log(sdkc.GreenGrassLogLevelDebug, "hello %s", "nils")
}
