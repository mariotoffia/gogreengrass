#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import {TestLambda} from '../lib/index'

const app = new cdk.App();
new TestLambda(app, 'TestLambdaStack');
app.synth();