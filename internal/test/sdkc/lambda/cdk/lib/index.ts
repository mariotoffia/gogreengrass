import * as cdk from '@aws-cdk/core';
import * as common from 'dist-common';
import {ConfigChangeService} from './services/configchange'
import { IQueueSubscription, IIoTCore } from 'cdk-custom-resources'

export interface DeviceServiceProperties extends common.Stack.ServiceProperties {
    /**
   * The endpoint URL to IoT core. Example https://a3gec3zlwkdb2-ats.iot.eu-central-1.amazonaws.com
   */
  iotCore: IIoTCore;
}

export class DeviceService extends common.Stack.ServiceConstruct {
  /**
   * Subscriptions that should be fed into buildingService for config-change
   */
  public readonly buildingSubscriptions: Array<IQueueSubscription>

  constructor(scope: cdk.Construct, id: string, props: DeviceServiceProperties) {
    super(scope, id, props);

    const configChange = new ConfigChangeService(this, 'devicesvc/configchange', {
      deployEnvironment: this.deployEnvironment,
      logLevel: this.logLevel,
      trace: this.trace,
      iotCore: props.iotCore,
    });

    this.buildingSubscriptions = [ configChange.buildingServiceSubscription ]
  }
}
