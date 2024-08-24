import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as lambda from 'aws-cdk-lib/aws-lambda';

interface ApiGatewayStackProps extends cdk.StackProps {
  userPool: cognito.UserPool;
  lambdaFunction: lambda.Function;
}

export class ApiGatewayStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: ApiGatewayStackProps) {
    super(scope, id, props);

    // Create a Cognito Authorizer for API Gateway
    const authorizer = new apigateway.CognitoUserPoolsAuthorizer(
      this,
      'TodoAuthorizer',
      {
        cognitoUserPools: [props.userPool],
      }
    );

    const api = new apigateway.LambdaRestApi(this, 'TodoApi', {
      handler: props.lambdaFunction,
      proxy: false,
    });

    const todos = api.root.addResource('todo');
    todos.addMethod(
      'POST',
      new apigateway.LambdaIntegration(props.lambdaFunction),
      {
        authorizer: authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );
    todos.addMethod(
      'GET',
      new apigateway.LambdaIntegration(props.lambdaFunction),
      {
        authorizer: authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const todosWithId = todos.addResource('{id}');
    todosWithId.addMethod(
      'GET',
      new apigateway.LambdaIntegration(props.lambdaFunction),
      {
        authorizer: authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );
    todosWithId.addMethod(
      'PUT',
      new apigateway.LambdaIntegration(props.lambdaFunction),
      {
        authorizer: authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );
    todosWithId.addMethod(
      'DELETE',
      new apigateway.LambdaIntegration(props.lambdaFunction),
      {
        authorizer: authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );
  }
}
