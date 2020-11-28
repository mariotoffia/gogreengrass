[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/mod/github.com/mariotoffia/gogreengrass)
[![GitHub Actions](https://img.shields.io/github/workflow/status/mariotoffia/gogreengrass/Go?style=flat-square)](https://github.com/mariotoffia/gogreengrass/actions?query=workflow%3AGo)

# Overview of gogreengrass 

:bulb: **Deploy your cloud go lambda onto greengrass core devices without alteration**

This library is a enabler to deploy standard aws go lambdas onto greengrass device lambdas. It also exposes the greengrass local API functions (greengrass SDK) to e.g. communicate with MQTT, local device shadow / secrets manager etc.

It also enables a go programmer to create much more efficient greengrass specific lambdas using the simplified lambda model if that is required.

**NOTE: This is still very much in development!**

## Example

Primary mode is to use GGC C Runtime and deploy golang lambdas as **greengrass lambda executable**. In this mode the go lambda is dynamically linked to the GGC C runtime and is much more optimal.

It is possible to use e.g. _CDK_ to deploy the lambda.

For example, create this simple lambda that you want to execute in same thread as the `main` function.

```golang

//go:generate gogreengrass -sdkc

func main() {
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
}
```

Make sure to have the shared library shim installed by `gogreengrass -sdkc` - _see Command Line Tool_. Since this file is decorated with generator pattern, `go generate` will create the library shim. Just do a standard _go_ build `go build -o testlambda main.go` and include it into your deployment. 

The following _CDK_ definition can be used to deploy the above lambda (_see sample: internal/test/sdkc/lambda_).

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
_Note that the lambda runtime is in this case arn:aws:greengrass:::runtime/function/executable_.

When doing `npm run deploy` it will show up in the IoT Core Greengrass console lambda for the greengrass group.

## Command Line Tool
Install the command line tool by `go get -u github.com/mariotoffia/gogreengrass`. This tool may be used in order to install the mock shared library if using C runtime.

```gogreengrass -h``` emits the following:

```bash
gogreengrass v0.0.6
Usage: gogreengrass [--out PATH] [--sdkc]

Options:
  --out PATH, -o PATH    The out path to write the shared AWS C runtime library mock. Default is /tmp/gogreengrass
  --sdkc, -l             Installs the c runtime shared library in /tmp/gogreengrass (or if -o, some other path)
  --help, -h             display this help and exit
  --version              display version and exit
  ```

## C Runtime

This is the preferred method to create your go lambdas. Use the _sdkc_ package to interact with the lambda runtime and greengrass specific APIs such as local device shadow / secrets manager or publish data on _MQTT_ etc.

### Lambda

The go version of the lambda runtime is layered. The "slim" or simple layer and the standard AWS lambda layer. Depending on the
use-case and performance on the device you may choose one over the other. 


#### Standard Lambda

The standard AWS lambda layer is behaving exactly the same as a standard cloud lambda, hence portable. 

When using the convenience function  `Start` it starts the lambda dispatcher on the main thread and the lambda gets invoked on the main thread. This is more or less the standard cloud version of it.

```go
	sdkc.Start(func(c context.Context, data MyEvent) (MyResponse, error) {
		// process the data here
	})
```

If you want to continue on main thread and fire up a dispatcher on a separate thread, you may use `StartWithOpts` to control this behavior. 

```go
	sdkc.StartWithOpts(
		func(c context.Context, data MyEvent) (MyResponse, error) { // <1>
		// process the data here
		}, 
		RuntimeOptionSeparateThread, // <2>
		true  // <3>
	)
```
<1> This function is executed on a single background thread. Hence, the invocations is serialized on that thread. The main thread continues.
<2> The specifies the separate thread behavior.
<3> If set to `true`, it will always fetch the payload. If `false` it is up to the lambda to fetch the data (_see below_).

Since the lambda function do take `MyEvent` the payload must be set to `true` in order to `Unmarshal` into that object. If lambda wants no payload or want to handle this itself, specify `false`.

```go
// registered lambda
func(c context.Context) error {
	r := NewRequestReader()
	b := make([]byte, 256)
	for {
		n, err := 	r.Read(b)
		if n > 0 {
			// process the b[:n]
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}
}
```

If you only want to read it all directly, e.g. want to do custom `Unmarshal` or other.

```go
if buf, err := ioutil.ReadAll(NewRequestReader()); err == nil {
   // the complete request payload is in buf
}
```

In short, you may use standard portable go lambdas or very specialized on the "standard" track. You may even do more lightweight lambdas using the more low level lambda.

#### Slim Lambda

The slim lambda is a non reflective and no `Unmarshal` path and hence is more optimized. You have two options to register the lambda (as with regular lambda). The `GGStart` and `GGStartWithOpts`, it works exactly the same on registration part, except that the lambda function is fixed. You have to do all reading, writing and others yourself.

```go
	sdkc.GGStart(func(lc *sdkc.LambdaContextSlim) { // <1>

		sdk.Log(sdkc.LogLevelInfo, "%s, %s\n", lc.ClientContext, lc.FunctionARN) // <2>
		sdk.Log(sdkc.LogLevelInfo, "Payload: %s\n", string(lc.Payload)) // <3>
       // <4>
	})
```
<1> The one and only function type to register.
<2> ClientContext is a string that you may unmarshal yourself.
<3> The payload is a `[]byte` (in this case it will be populated since `GGStart` do set _payload_ to `true`)
<4> If you want to return data, you have to write either error output or return payload yourself.

As with standard lambda, one may register using `GGStartWithOpts` to change if background thread or read / write data yourself.

To write a response, do create a `NewResponseWriter` and do a `Write(buff)`. If any errors is returned, use the `NewErrorResponseWriter` and do a `Write(buff)`. Both of them implements the `io.Writer` interface.

### Install C Runtime SDK Mock Library 

You need to have the mock version of the shared library. Either follow the instructions in the [greengrass core C SDK](https://github.com/aws/aws-greengrass-core-sdk-c) or use the `gogreengrass` ability to store _libaws-greengrass-core-sdk-c.so_ in your _/tmp/gogreengrass_ folder. 

```bash
gogreengrass -sdkc
```

This writes the shared library (shim) that makes your go lambdas build and run. When deployed onto the greengrass core device, the real shared library is already present (**this library shall never be part of the package**) - see the [greengrass core C SDK](https://github.com/aws/aws-greengrass-core-sdk-c) for more information.