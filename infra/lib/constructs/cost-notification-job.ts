import * as cdk from "aws-cdk-lib";
import * as events from "aws-cdk-lib/aws-events";
import * as events_targets from "aws-cdk-lib/aws-events-targets";
import * as iam from "aws-cdk-lib/aws-iam";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as logs from "aws-cdk-lib/aws-logs";
import * as sns from "aws-cdk-lib/aws-sns";
import { Construct } from "constructs";
import path from "path";

export type CostNotificationJobProps = {
  lambdaAssetRoot: string;
  schedule: events.Schedule;
  topic: sns.ITopic;
};

export class CostNotificationJob extends Construct {
  readonly function: lambda.Function;

  constructor(scope: Construct, id: string, props: CostNotificationJobProps) {
    super(scope, id);

    this.function = new lambda.Function(this, "CostNotificationFunc", {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      architecture: lambda.Architecture.ARM_64,
      handler: "bootstrap",
      code: lambda.Code.fromAsset(
        path.join(props.lambdaAssetRoot, "cost-notification"),
      ),
      memorySize: 128,
      timeout: cdk.Duration.seconds(30),
      environment: {
        SNS_TOPIC_ARN: props.topic.topicArn,
      },
      logGroup: new logs.LogGroup(this, "CostNotificationFuncLogGroup", {
        logGroupName: `/aws/lambda/DiscordMessenger-CostNotification`,
        retention: logs.RetentionDays.ONE_YEAR,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
      }),
    });

    props.topic.grantPublish(this.function);

    this.function.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ["ce:GetCostAndUsage", "ce:GetCostForecast"],
        resources: ["*"],
      }),
    );

    new events.Rule(this, "ScheduleRule", {
      schedule: props.schedule,
      targets: [
        new events_targets.LambdaFunction(this.function, {
          retryAttempts: 3,
          maxEventAge: cdk.Duration.hours(1),
        }),
      ],
    });
  }
}
