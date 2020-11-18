#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { DeviceService } from '../lib/index';

const app = new cdk.App();
var stack = new cdk.Stack(app, 'deviceservicestack');
new DeviceService(stack, 'deviceservice');