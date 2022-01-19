![Chatterbug and Slack](/docs/chatterbug-and-slack.png)

Integrate you Chatterbug leaderboard directly in to Slack. In this repo you'll
find all the tools you need to query Chatterbug's GraphQL API, an integration
with AWS Lambda and API Gateway, and a Lambda function that returns messages in
a format consumable by Slack.

After you have everything deployed, you can [create a custom slash command in
Slack](https://www.hongkiat.com/blog/custom-slash-command-slack/) to get the
relevant messages in Slack itself.

## Development

There's a small binary that can be built for testing purposes in `cmd/chatterbug-leaderboard`.

```bash
$ go build ./cmd/chatterbug-leaderboard
```

## Deployment

### Packaging and deployment

There's a make target that builds the Lambda function for deployment. It'll
output a zip file that can be uploaded to S3.

```bash
$ make package
rm -rf bin
mkdir bin
GOOS=linux go build -o ./bin/get-chatterbug-leaderboard ./pkg/lambda/get-chatterbug-leaderboard
cd ./bin && zip get-chatterbug-leaderboard.zip get-chatterbug-leaderboard
  adding: get-chatterbug-leaderboard (deflated 48%)
$ ls -lah ./bin
total 31408
drwxr-xr-x   4 brad  staff   128B Jan 19 09:58 .
drwxr-xr-x  16 brad  staff   512B Jan 19 09:58 ..
-rwxr-xr-x   1 brad  staff    10M Jan 19 09:58 get-chatterbug-leaderboard
-rw-r--r--   1 brad  staff   5.2M Jan 19 09:58 get-chatterbug-leaderboard.zip
```

Whenever you update the code you need to tell Lambda that you've done so.

```bash
$ aws lambda update-function-code \
		--function-name get-chatterbug-leaderboard 
		--s3-bucket chatterbug-leaderboard-artifacts-1-bucket-abc12cedf3456 \
		--s3-key get-chatterbug-leaderboard.zip \
		--profile my-profile \
		--region eu-central-1
```

### CloudFormation

Everything for running this in Lambda is captured in CloudFormatiion. There are
two stacks to manage. The first is for storing the code that will be executed
in Lambda and the second is for lambda and API gateway itself. They're split
between two stacks because the relevant code needs to exist before the Lambda
function can be created.

The process looks like this.

* Create the `artifacts` stack.
* Build, package, and upload the code to S3.
* Create the `lambda` stack.

So first, let's create the artifacts stack.

```bash
$ aws cloudformation create-stack \
		--stack-name chatterbug-leaderboard-artifacts-1 \
		--template-body file://./scripts/cloudformation-artifacts.yaml \
		--profile my-profile \
		--region eu-central-1 \
		--capabilities CAPABILITY_IAM
```

Then you need to upload the code to the bucket that was create. The bucket name
is randomly generated by CloudFormation so let's get the bucket name using the
following.

```bash
$ aws cloudformation describe-stacks \
		--stack-name chatterbug-leaderboard-artifacts-1 \
		--profile my-profile \
		--region eu-central-1
{
    "Stacks": [
        {
            "StackId": "arn:aws:cloudformation:eu-central-1:0123456789:stack/chatterbug-leaderboard-artifacts-1/b62f58e3-d4ea-4ab4-ad70-1d30b37fc19f",
            "StackName": "chatterbug-leaderboard-artifacts-1",
            "Description": "Chatterbug Leaderboard API artifacts deployment",
            "CreationTime": "2022-01-19T08:38:50.115Z",
            "RollbackConfiguration": {},
            "StackStatus": "CREATE_COMPLETE",
            "DisableRollback": false,
            "NotificationARNs": [],
            "Capabilities": [
                "CAPABILITY_IAM"
            ],
            "Outputs": [
                {
                    "OutputKey": "BucketARN",
                    "OutputValue": "arn:aws:s3:::chatterbug-leaderboard-artifacts-1-bucket-abc12cedf3456"
                },
                {
                    "OutputKey": "BucketName",
                    "OutputValue": "chatterbug-leaderboard-artifacts-1-bucket-abc12cedf3456"
                }
            ],
            "Tags": [],
            "EnableTerminationProtection": false,
            "DriftInformation": {
                "StackDriftStatus": "NOT_CHECKED"
            }
        }
    ]
}
```

Now upload the relevant code. Notice the reference in the bucket name to the
previous CloudFormation output.

```bash
$ make package
$ aws s3 cp ./bin/get-chatterbug-leaderboard.zip s3://chatterbug-leaderboard-artifacts-1-bucket-abc12cedf3456/get-chatterbug-leaderboard.zip \
		--profile my-profile \
		--region eu-central-1
```

Now create the lambda stack. Note the parameters need to be updated to reflect
the actual stack name you used as well as include your Chatterbug API token.

```bash
$ aws cloudformation create-stack \
		--stack-name chatterbug-leaderboard-lambda-1 \
		--template-body file://./scripts/cloudformation-lambda.yaml \
		--parameters ParameterKey=ArtifactsStackName,ParameterValue=chatterbug-leaderboard-artifacts-1 ParameterKey=ChatterbugAPIToken,ParameterValue=DEADBEEF\
		--profile my-profile \
		--region eu-central-1 \
		--capabilities CAPABILITY_IAM
```
