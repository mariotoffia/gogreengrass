import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as common from 'dist-common';
import * as iam from '@aws-cdk/aws-iam';
import * as sqs from '@aws-cdk/aws-sqs';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources'
import { IQueueSubscription, QueueSubscriptionBuilder, IIoTCore } from 'cdk-custom-resources'
import * as statement from 'cdk-iam-floyd'

import path = require("path");

export interface ConfigChangeProperties extends common.Stack.ServiceProperties {
    /**
   * The endpoint URL to IoT core. Example https://a3gec3zlwkdb2-ats.iot.eu-central-1.amazonaws.com
   */
  iotCore: IIoTCore;
}

export class ConfigChangeService extends common.Stack.ServiceConstruct {
  /** 
 * Microservice within the Device Service provider that receives via SQS a ConfigChange event
 * and handles things, device shadows and rules accordingly.
 * 
 * @returns the function that manages the things, device shadow et.al using ConfigChange as it's basis
 **/
  public readonly configchange: lambda.Function;

  /**
   * Exposes the config-change queue that needs a subscription of Building.
   */
  public readonly configchangeQueue: sqs.IQueue;
  /**
   * The subscription to the buildingsvc config-change event. It subscribes to 
   * when building instances are created, updated or deleted.
   */
  public readonly buildingServiceSubscription: IQueueSubscription;


  constructor(scope: cdk.Construct, id: string, props: ConfigChangeProperties) {
    super(scope, id, props);

    const lambdaTimeout = 60
    this.configchange = new lambda.Function(this, 'devicesvc/devconfigchange', {
      runtime: lambda.Runtime.GO_1_X,
      functionName: 'devicesvc-devconfigchange',
      handler: 'devconfigchange',
      code: lambda.Code.fromAsset(path.join(__dirname, '../../../_output/devconfigchange')),
      tracing: this.trace,
      timeout: cdk.Duration.seconds(lambdaTimeout),
      environment: {
        CB_LAMBDA_TYPE: 'native',
        CB_LOG_LEVEL: this.logLevel,
        CB_DEPLOY_ENV: this.deployEnvironment,
        CB_CONSOLE_LOG: 'false',
        CB_DS_EP: props.iotCore.customEndpointURL,
      }
    })

    this.configchangeQueue = new sqs.Queue(this, 'devicesvc/devconfigchange-queue', {
      queueName: 'devconfigchange-queue',
      encryption: sqs.QueueEncryption.KMS_MANAGED,
      visibilityTimeout: cdk.Duration.seconds(lambdaTimeout * 6),
      retentionPeriod: cdk.Duration.days(14),
      receiveMessageWaitTime: cdk.Duration.seconds(20),
    });

    this.configchange.addEventSource(new SqsEventSource(this.configchangeQueue, {
      batchSize: 10 // default
    }));

    this.buildingServiceSubscription = new QueueSubscriptionBuilder("devconfigchange").
      type("BLD#BLD").
      queue(this.configchangeQueue).
      build();

    const ps = new statement.Iot()
      .toGetThingShadow()
      .toDescribeThing()
      .toCreateThing()
      .toDeleteThing()
      .toCreateThingType()
      .toUpdateThingShadow();      
    // TODO: arn to IoT core?
    ps.addAllResources()
    this.configchange.addToRolePolicy(ps);
  }
}
