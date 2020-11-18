# gogreengrass

This library has two modes of operation to deploy and execute the _go lambdas_. 

1. One is standard lambda and thus wrapped as a python 3.7 and may be managed in AWS console and deployed onto GGC. The go lambda will be invoked through the python wrapper.  

2. Use green grass lambda executable. In this mode the go lambda is dynamically linked to the GGC C runtime and is much more optimal. The only downside is that the lambdas is only visible inside the Greengrass Core administration utility and not in the standard Lambda view. Other than that they behave exactly as a lambda, but i much more lightweight.

Both may use e.g. _CDK_ to deploy the lambda. To e.g. deploy the following greengrass executable

```golang
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/mariotoffia/gogreengrass/sdkc"
)

func main() {
	sdkc.Log(sdkc.LogLevelInfo, "Register lambda...\n")

	type MyEvent struct {
		Data  int    `json:"data"`
		Hello string `json:"hello"`
	}

	type MyResponse struct {
		Age   int    `json:"age"`
		Topic string `json:"topic"`
	}

	sdkc.Start(func(c context.Context, data MyEvent) (MyResponse, error) {

		lc, _ := lambdacontext.FromContext(c)

		fmt.Printf(
			"context: %v, topic: %s, data: '%v'\n",
			lc, lc.ClientContext.Custom["subject"], data,
		)

		resp := MyResponse{Age: 19, Topic: "feed/myfunc"}

		sdkc.NewQueue().PublishObject(
			"feed/testlambda", sdkc.QueueFullPolicyOptionAllOrError, &resp,
		)

		return resp, nil
	})

	sdkc.Log(sdkc.LogLevelInfo, "Exit main()\n")
}
```

The following _CDK_ definition was used (_see sample: internal/test/sdkc/lambda_).

```typescript
import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import path = require("path");

const GREENGRASS_EXECUTABLE = new lambda.Runtime('arn:aws:greengrass:::runtime/function/executable')
export class TestLambda extends cdk.Stack {

  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const testlambda = new lambda.Function(this, 'testlambda', {
      runtime: GREENGRASS_EXECUTABLE,
      functionName: 'testlambda',
      handler: 'testlambda',
      code: lambda.Code.fromAsset(path.join(__dirname, '../../_out/testlambda')),
      timeout: cdk.Duration.seconds(30),
      currentVersionOptions: {
        removalPolicy: cdk.RemovalPolicy.RETAIN,
      }
    });

    testlambda.currentVersion.addAlias('live')
  }
}
```
_Note that the lambda runtime is in this case arn:aws:greengrass:::runtime/function/executable_ (if python use `Runtime.PYTHON_3_7`).

When doing `npm run deploy` it will show up in the greengrass console to be added to a group.

## Python wrapper
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

If you stick with the auto generated, make sure that your handler is set to: *glue.function_handler*. If you bind the IoT Core -> Lambda on e.g. topic: _invoke/ggtest_, open the test utility to and post to _invoke/ggtest_:

```json
{
  "data": 49,
  "hello": "world!"
}
```

The log should print out the following in _greengrass/ggc/var/log/user/your region/account/ggtest.log_

```
[2020-11-16T21:55:08.122Z][INFO]-context: &{c760d735-5e69-4ffc-5a76-1460197dcd71 arn:aws:lambda:eu-central-1:033549287452:function:ggtest:11 { } {{   } map[] map[subject:invoke/ggtest]}}, (from-)topic: invoke/ggtest, data: '{49 world!}'
```

The above line is logged out from the following statement (in sample code)

```go
fmt.Printf(
			"context: %v, (from-)topic: %s, data: '%v'\n",
			lc, lc.ClientContext.Custom["subject"], data,
		)
```

`gogreengrass` has the following flags to control the glue generation:

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
