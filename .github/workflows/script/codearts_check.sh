
IAM_DATA=$(cat <<EOF
{
    "auth": {
        "identity": {
            "methods": [
                "password"
            ],
            "password": {
                "user": {
                    "domain": {
                        "name": "freesky-edward"
                    },
                    "name": "codearts_test",
                    "password": "$CODEARTS_PASSWORD"
                }
            }
        },
        "scope": {
            "project": {
                "name": "cn-north-4"
            }
        }
    }
}
EOF
)

response=$(curl -s -i --location 'https://iam.myhuaweicloud.com/v3/auth/tokens?nocatalog=true' \
  --header 'Content-Type: application/json' \
  --data "$IAM_DATA")

# Extract the X-Subject-Token from the response
token=$(echo "$response" | grep "X-Subject-Token" | awk '{print $2}' | tr -d '\r')

echo "X-Subject-Token: $token"


DATA=$(cat <<EOF
{
  "sources" : [ {
    "type" : "code",
    "params" : {
      "git_type" : "github",
      "default_branch" : "feture_experence",
      "git_url" : "https://github.com/opensourceways/message-push.git",
      "build_params" : {
        "build_type" : "branch",
        "event_type" : "Manual",
        "target_branch" : "$BRANCH_NAME"
      }
    }
  } ],
  "description" : "运行描述",
  "variables" : [ {
    "name" : "repo",
    "value" : "message-push"
  } ,
  {
    "name" : "owner",
    "value" : "opensourceways"
  }
  ,
  {
    "name" : "pr_id",
    "value" : "$PR_ID"
  }
  ]
}
EOF
)

CODEARTS_PIPELINE="https://cloudpipeline-ext.cn-north-4.myhuaweicloud.com/v5/3a76c1785dda4b13a399937d1978f240/api/pipelines/ea5489fc52984d36a1b33fc6378ddb85/run"

curl --location "$CODEARTS_PIPELINE" \
--header "X-Auth-Token:$token" \
--header "Content-Type: application/json" \
--data "$DATA"