import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cognito from 'aws-cdk-lib/aws-cognito';

export class CognitoStack extends cdk.Stack {
  public readonly userPool: cognito.UserPool;
  public readonly userPoolClient: cognito.UserPoolClient;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create a Cognito User Pool
    this.userPool = new cognito.UserPool(this, 'TodoUserPool', {
      userPoolName: 'TodoUserPool',
      signInAliases: {
        email: true,
        phone: true,
      },
      autoVerify: { email: true },
      passwordPolicy: {
        minLength: 8,
        requireLowercase: true,
        requireUppercase: true,
        requireDigits: true,
        requireSymbols: false,
      },
      accountRecovery: cognito.AccountRecovery.EMAIL_ONLY,
      userVerification: {
        emailSubject: 'Verify your email',
        emailBody: 'Your verification code is {####}',
        emailStyle: cognito.VerificationEmailStyle.CODE,
      },
    });

    // Create a User Pool Client
    this.userPoolClient = new cognito.UserPoolClient(
      this,
      'TodoUserPoolClient',
      {
        userPool: this.userPool,
        userPoolClientName: 'TodoUserPoolClient',
        generateSecret: false,
        authFlows: {
          adminUserPassword: true,
          custom: true,
          userSrp: true,
          userPassword: true,
        },
        preventUserExistenceErrors: true,
        oAuth: {
          flows: {
            implicitCodeGrant: true,
          },
          scopes: [
            cognito.OAuthScope.OPENID,
            cognito.OAuthScope.EMAIL,
            cognito.OAuthScope.PROFILE,
          ],
          callbackUrls: ['https://jwt.io'],
        },
        supportedIdentityProviders: [
          cognito.UserPoolClientIdentityProvider.COGNITO,
        ],
      }
    );

    new cognito.UserPoolDomain(this, 'TodoCognitoDomain', {
      userPool: this.userPool,
      cognitoDomain: {
        domainPrefix: `todo-user-pool-domain-${cdk.Aws.ACCOUNT_ID}`,
      },
    });

    new cognito.CfnUserPoolUser(this, 'TodoCognitoAuthorizerUserPoolUser', {
      userPoolId: this.userPool.userPoolId,
      username: 'sachindra.maharjan4@gmail.com',
      userAttributes: [
        {
          name: 'email',
          value: 'sachindra.maharjan4@gmail.com',
        },
      ],
    });

    // Output relevant details
    new cdk.CfnOutput(this, 'UserPoolId', {
      value: this.userPool.userPoolId,
    });

    new cdk.CfnOutput(this, 'UserPoolClientId', {
      value: this.userPoolClient.userPoolClientId,
    });
  }
}
