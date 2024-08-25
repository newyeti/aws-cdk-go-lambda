import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CognitoStack } from './congnito-stack';
import { LambdaStack } from './lambda-stack';
import { AuthLambdaStack } from './lambda-auth-stack';
import { ApiGatewayStack } from './apigateway-stack';

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Instantiate the Cognito Stack
    const cognitoStack = new CognitoStack(this, 'CognitoStack');

    // Instantiate the Lambda Stack
    const lambdaStack = new LambdaStack(this, 'LambdaStack');

    const authLambdaStack = new AuthLambdaStack(this, 'AuthLambdaStack', {
      userPool: cognitoStack.userPool,
      userPoolClient: cognitoStack.userPoolClient,
    });

    // Instantiate the API Gateway Stack, passing in the User Pool and Lambda
    new ApiGatewayStack(this, 'ApiGatewayStack', {
      userPool: cognitoStack.userPool,
      authFunction: authLambdaStack.authLambda,
      todoFunction: lambdaStack.todoLambda,
    });
  }
}
