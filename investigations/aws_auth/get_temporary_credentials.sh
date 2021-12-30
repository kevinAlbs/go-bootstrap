# Run script with . ./get_temporary_credentials.sh
# Create temporary credentials with the minimum lifetime (15 minutes).
TEMPCREDS=$(aws sts get-session-token --duration-seconds 900)

export TEST_AWS_TEMP_ACCESS_KEY_ID=$(echo $TEMPCREDS | jq -r ".Credentials.AccessKeyId")
export TEST_AWS_TEMP_SECRET_ACCESS_KEY=$(echo $TEMPCREDS | jq -r ".Credentials.SecretAccessKey")
export TEST_AWS_TEMP_SESSION_TOKEN=$(echo $TEMPCREDS | jq -r ".Credentials.SessionToken")

CALLERIDENTITY=$(aws sts get-caller-identity)
export TEST_AWS_ARN=$(echo $CALLERIDENTITY | jq -r ".Arn")