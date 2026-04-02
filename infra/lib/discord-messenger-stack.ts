import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { AppConfig } from "./app-config";
import { DiscordNotifier } from "./constructs/discord-notifier";
import { CostNotificationJob } from "./constructs/cost-notification-job";

export type DiscordMessengerStackProps = cdk.StackProps & {
  config: AppConfig;
};

export class DiscordMessengerStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: DiscordMessengerStackProps) {
    super(scope, id, props);

    cdk.Tags.of(this).add("app", "discord-messenger");

    const discordNotifier = new DiscordNotifier(this, "DiscordNotifier", {
      lambdaAssetRoot: props.config.lambdaAssetRoot,
      discordWebhookParameterName: props.config.discordWebhookParameterName,
    });

    new CostNotificationJob(this, "CostNotificationJob", {
      lambdaAssetRoot: props.config.lambdaAssetRoot,
      schedule: props.config.costNotificationSchedule,
      topic: discordNotifier.topic,
    });
  }
}
