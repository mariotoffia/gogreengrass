import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import path = require("path");

export class DeviceService extends cdk.Stack {

  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const testlambda = new lambda.Function(this, 'testlambda', {
      runtime: lambda.Runtime.GO_1_X,
      functionName: 'testlambda',
      handler: 'testlambda',
      code: lambda.Code.fromAsset(path.join(__dirname, '/../_out/testlambda')),
      timeout: cdk.Duration.seconds(30),
      environment: {
        IS_ENV_ON_GGC_SET: 'yes they are!',
      }
    });

  }
}
