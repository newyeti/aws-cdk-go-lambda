import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as go from '@aws-cdk/aws-lambda-go-alpha';
import { Runtime } from 'aws-cdk-lib/aws-lambda';
import { RetentionDays } from 'aws-cdk-lib/aws-logs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as iam from 'aws-cdk-lib/aws-iam';
import { UserPool, UserPoolClient } from 'aws-cdk-lib/aws-cognito';

interface AuthLambdaStackProps extends cdk.StackProps {
  userPool: UserPool;
  userPoolClient: UserPoolClient;
}

export class AuthLambdaStack extends cdk.Stack {
  public readonly authLambda: lambda.Function;

  constructor(scope: Construct, id: string, props: AuthLambdaStackProps) {
    super(scope, id, props);

    this.authLambda = new go.GoFunction(this, 'AuthFuntion', {
      entry: '../function/auth',
      runtime: Runtime.PROVIDED_AL2023,
      logRetention: RetentionDays.ONE_WEEK,
      environment: {
        UserpoolID: props.userPool.userPoolId,
        ClientsID: props.userPoolClient.userPoolClientId,
      },
    });

    // Grant the Lambda function permission to invoke AdminInitiateAuth
    this.authLambda.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['cognito-idp:AdminInitiateAuth'],
        resources: [props.userPool.userPoolArn],
      })
    );
  }
}
