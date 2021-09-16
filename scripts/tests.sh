#!/bin/bash
set -e
set -o pipefail

_TFCODING_DOCKER_TAG="${TFCODING_DOCKER_TAG:-"unfor19/tfcoding:0.15.0-beta2-latest"}"


error_msg(){
    local msg=$1
    echo -e "\e[31m[ERROR]\e[0m $msg"
    export DEBUG=1
    exit 1
}


log_msg(){
  local msg="$1"
  echo -e "[LOG] $(date) :: $msg"
}

should(){
    local expected=$1
    local test_name=$2
    local expr=$3
    local output_code
    echo "-------------------------------------------------------"
    echo "[LOG] $test_name - Should $expected"
    echo "[LOG] Executing: $expr"
    output_msg="$(trap '$expr' EXIT)"
    output_code=$?

    echo -e "[LOG] Output Code: ${output_code}"
    echo -e "[LOG] Output Msg:\n\n${output_msg}\n"

    if [[ "$expected" == "pass" && "$output_code" -eq 0 && ! "$output_msg" =~ .*(ERROR|Error|error).* ]]; then
        echo -e "\e[92m[SUCCESS]\e[0m Test passed as expected"
    elif [[ "$expected" == "fail" && "$output_code" -ne 0 ]] || [[ "$expected" == "fail" && "$output_msg" =~ .*(ERROR|Error|error).* ]] ; then
        echo -e "\e[92m[SUCCESS]\e[0m Test failed as expected"
    else
        error_msg "Test output is not expected, terminating"
    fi
}

ssm_put_parameter(){
    local p_name="$1"
    local p_value="${2:-"empty"}"
    local p_type="${3:-"String"}"
    local p_key_id="${4:-"alias/aws/ssm"}"
    declare -a arg_key_id=(--key-id "$p_key_id")

    if [[ -z "$p_name" ]]; then
        error_msg "Must provide parameter name in ssm_put_parameter"
    fi
    if [[ "$p_type" = "SecureString" ]]; then
        aws ssm --endpoint-url="$_AWS_SSM_ENDPOINT_URL" put-parameter --overwrite --name "$p_name" --value "$p_value" --type "$p_type" ${arg_key_id[*]} 1>/dev/null
    else
        aws ssm --endpoint-url="$_AWS_SSM_ENDPOINT_URL" put-parameter --overwrite --name "$p_name" --value "$p_value" --type "$p_type" 1>/dev/null
    fi   
}

run_app(){
    local output
    set +e
    output="$(./parzival "$@" 2>&1)"
    sleep 1
    set -e
    echo "$output"
}

# Initialize Variables
export AWS_ACCESS_KEY_ID="mock_aws"
export AWS_SECRET_ACCESS_KEY="mock_aws"
export AWS_SESSION_KEY="mock_aws"
export AWS_REGION="us-east-1"
_AWS_SSM_ENDPOINT_URL="${AWS_SSM_ENDPOINT_URL:-"http://localhost:4566"}"
_SKIP_PARAM_CREATION="${SKIP_PARAM_CREATION:-"true"}"

# Tests
make up-localstack
source scripts/wait_for_endpoints.sh "http://localhost:4566/health"
if [[ "$_SKIP_PARAM_CREATION" != "true" ]]; then
    log_msg "Creating parameters ..."
    ssm_put_parameter "/myapp/dev/LOG_LEVEL" "INFO" "String"
    ssm_put_parameter "/myapp/dev/GOOGLE_CLIENT_ID" "1a2s3d4f" "SecureString"
    ssm_put_parameter "/myapp/dev/GOOGLE_CLIENT_SECRET" "W1llyW0naO0mpaL00mp4" "SecureString"
    log_msg "Completed creating parameters"
fi

log_msg "Build application ..."
go build
log_msg "Completed building application"

log_msg "Running Test Suite"
should pass "Help Menu" "run_app get --help"
should pass "Get No Arguments" "run_app get --localstack"
should pass "Get Minimum Arguments" "run_app get --localstack -p /myapp/"
should pass "Get Complete Arguments" "run_app get --localstack -p /myapp/ -o .test.json"
should fail "Get Unknown Argument" "run_app get --local"
log_msg "Completed Test Suite"

log_msg "Finished"
