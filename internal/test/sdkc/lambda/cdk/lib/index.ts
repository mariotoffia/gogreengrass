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

