package main

import "C"
import (
	"fmt"
)

// Setup is the testrunner for this
func Setup() {
	//GGDataplane.Publish("hello", sdk.QueueFullPolicyAllOrException, `{"user":"mario"}`)

	s, _ := GGDataplane.GetThingShadow("myThingName")
	fmt.Printf("--> '%s'\n", s)

	s, _ = GGDataplane.UpdateThingShadow("myThingName", `{"state": { "reported": { "state": 17}}}`)
	fmt.Printf("--> '%s'\n", s)

	s, _ = GGDataplane.DeleteThingShadow("myThingName")
	fmt.Printf("--> '%s'\n", s)
}

func main() {}
