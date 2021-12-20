export ACCESS_KEY_ID=$(cat ~/.csfle/aws_creds_assumerole.json | jq -r ".accessKeyId")
export SECRET_ACCESS_KEY=$(cat ~/.csfle/aws_creds_assumerole.json | jq -r ".secretAccessKey")

TEMPCREDS=$(
    AWS_ACCESS_KEY_ID=$ACCESS_KEY_ID \
    AWS_SECRET_ACCESS_KEY=$SECRET_ACCESS_KEY \
    AWS_DEFAULT_REGION=us-east-1 \
    aws sts assume-role --role-arn "arn:aws:iam::524754917239:role/assumeme2" --role-session-name "foo"
)

export TEMP_ACCESS_KEY_ID=$(echo $TEMPCREDS | jq -r ".Credentials.AccessKeyId")
export TEMP_SECRET_ACCESS_KEY=$(echo $TEMPCREDS | jq -r ".Credentials.SecretAccessKey")
export TEMP_SESSION_TOKEN=$(echo $TEMPCREDS | jq -r ".Credentials.SessionToken")
export CMK_ARN="arn:aws:kms:us-east-1:524754917239:key/c6e5a131-b3db-4886-85fe-1083cd752881"
export CMK_REGION="us-east-1"