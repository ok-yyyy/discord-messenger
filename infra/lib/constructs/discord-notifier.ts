import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as logs from "aws-cdk-lib/aws-logs";
import * as sns from "aws-cdk-lib/aws-sns";
import * as sns_subscriptions from "aws-cdk-lib/aws-sns-subscriptions";
import * as ssm from "aws-cdk-lib/aws-ssm";
import { Construct } from "constructs";
import path from "path";

export type DiscordNotifierProps = {
  lambdaAssetRoot: string;
  discordWebhookParameterName: string;
};

export class DiscordNotifier extends Construct {
  readonly topic: sns.Topic;
  readonly function: lambda.Function;

  constructor(scope: Construct, id: string, props: DiscordNotifierProps) {
    super(scope, id);

    this.topic = new sns.Topic(this, "Topic", {
      topicName: "DiscordMessengerTopic",
    });

    const discordWebhookParameter =
      ssm.StringParameter.fromSecureStringParameterAttributes(
        this,
        "DiscordWebhookParameter",
        {
          parameterName: props.discordWebhookParameterName,
        },
      );

    this.function = new lambda.Function(this, "PostDiscordFunc", {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      architecture: lambda.Architecture.ARM_64,
      handler: "bootstrap",
      code: lambda.Code.fromAsset(
        path.join(props.lambdaAssetRoot, "post-discord"),
      ),
      memorySize: 128,
      timeout: cdk.Duration.seconds(30),
      environment: {
        DISCORD_WEBHOOK_PARAMETER_NAME: discordWebhookParameter.parameterName,
      },
      logGroup: new logs.LogGroup(this, "PostDiscordFuncLogGroup", {
        logGroupName: `/aws/lambda/DiscordMessenger-PostDiscord`,
        retention: logs.RetentionDays.ONE_YEAR,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
      }),
    });

    discordWebhookParameter.grantRead(this.function);

    this.topic.addSubscription(
      new sns_subscriptions.LambdaSubscription(this.function),
    );
  }
}
