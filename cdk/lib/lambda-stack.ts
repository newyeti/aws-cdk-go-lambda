import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as go from '@aws-cdk/aws-lambda-go-alpha';
import { Runtime } from 'aws-cdk-lib/aws-lambda';
import { RetentionDays } from 'aws-cdk-lib/aws-logs';
import * as lambda from 'aws-cdk-lib/aws-lambda';

export class LambdaStack extends cdk.Stack {
  public readonly todoLambda: lambda.Function;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const table = new dynamodb.Table(this, 'Todo', {
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      // sortKey: { name: 'createdOn', type: dynamodb.AttributeType.STRING },
      tableName: 'todo',
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    this.todoLambda = new go.GoFunction(this, 'TodoFuntion', {
      entry: '../function',
      runtime: Runtime.PROVIDED_AL2023,
      logRetention: RetentionDays.ONE_WEEK,
      environment: {
        TABLE_NAME: table.tableName,
        PrimaryKey: 'id',
      },
    });

    table.grantReadWriteData(this.todoLambda);
  }
}
