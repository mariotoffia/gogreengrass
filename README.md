# gogreengrass
Python wrapper to deploy go lambdas into green grass. It also encapsulates the python greengrass SDK so go code may invoke
the python SDK to do device shadow operations, secrets manager and publish data onto _MQTT_.

Install the command line tool by `go get -u github.com/mariotoffia/gogreengrass`. This tool may be used in order to generate
the necessary go and python shim. For example (_see internal/example/lambda/main_):

The _main.go_ contains the lambda code in this sample.

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/mariotoffia/gogreengrass/sdk"
)

// once is invoked once when ggc lambda startup
func once() {

	type MyEvent struct {
		Data  int    `json:"data"`
		Hello string `json:"hello"`
	}

	type MyResponse struct {
		Age   int    `json:"age"`
		Topic string `json:"topic"`
	}

	sdk.Register(func(c context.Context, data MyEvent) (MyResponse, error) {

		lc, _ := lambdacontext.FromContext(c)

		fmt.Printf(
			"context: %v, (from-)topic: %s, data: '%v'\n",
			lc, lc.ClientContext.Custom["subject"], data,
		)

		return MyResponse{Age: 19, Topic: "feed/myfunc"}, nil
	})
}
func main() {
	// Standard AWS cloud lambda initialization code
	// Your standard lambda SDK registration code here
}
```

When you invoke the `gogreengrass -d`, it will generate _glue.py_ and _glue.go_ where the python code has the lambda entry
```python
# The actual function that you need to bind the lambda entry point to
def function_handler(event, context):
    result = invokeJSON(context, event)
    return result
```

You may remove that and provide with your own. The _glue.go_ give the python runtime the ability to invoke your lambda code as well as
able to call python SDK functions such as publish to MQTT etc. You may modify the glue code if you wish (you always have the ability to re-generate using `gogreengrass -d`).

This is only needed to be done once (if not upgraded the gogreengrass sdk and breaking changes has been introduced) and thus, you may
check-in the code. Publish the lambda as python 3.x runtime and then it will be deployable in a greengrass core group.

```bash
gogreengrass v0.0.4
Usage: gogreengrass [--out PATH] [--package PACKAGE] [--binary BINARY] [--downloadsdk] [--force]

Options:
  --out PATH, -o PATH    The out path to write the generated go and python files. Default is current directory
  --package PACKAGE, -p PACKAGE
                         an optional package name instead of main
  --binary BINARY, -b BINARY
                         an optional name of the binary that the build system produces, default is foldername.o
  --downloadsdk, -d      If set to true, it will download the python sdk (it is needed to be in current folder)
  --force, -f
  --help, -h             display this help and exit
  --version              display version and exit
```

Above is the available options to generate the glue code.
