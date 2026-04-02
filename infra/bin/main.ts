import * as cdk from "aws-cdk-lib";
import path from "path";
import { createAppConfig } from "../lib/app-config";
import { DiscordMessengerStack } from "../lib/discord-messenger-stack";

const app = new cdk.App();

const config = createAppConfig({
  lambdaAssetRoot: path.join(__dirname, "..", "..", "lambda", "dist"),
});

new DiscordMessengerStack(app, "DiscordMessengerStack", {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
  config,
});
