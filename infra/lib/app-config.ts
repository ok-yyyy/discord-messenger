import * as events from "aws-cdk-lib/aws-events";

export type AppConfig = {
  lambdaAssetRoot: string;
  discordWebhookParameterName: string;
  costNotificationSchedule: events.Schedule;
};

export function createAppConfig(input: { lambdaAssetRoot: string }): AppConfig {
  return {
    lambdaAssetRoot: input.lambdaAssetRoot,
    discordWebhookParameterName: "/discord-messenger/discord-webhook",
    costNotificationSchedule: events.Schedule.cron({
      minute: "0",
      hour: "9",
      day: "5/5",
      month: "*",
      year: "*",
    }),
  };
}
