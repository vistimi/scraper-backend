#!/bin/bash
ecsClusterName=scraperClusterECS
ecsServiceName=scraperServiceFargate
nameTaskDefinition=scraperDefinitionFargate
region=us-east-1

# get latest task
GET_LATEST_TASK_ARN=$(aws ecs list-task-definitions \
    --region ${region} \
    --family-prefix ${nameTaskDefinition} \
    --sort DESC \
    --query 'taskDefinitionArns[0]' \
    --output text
)
latestTaskArn=${GET_LATEST_TASK_ARN}
echo "latestTaskArn task = "${latestTaskArn}

# desired tasks
UPDATE_FARGATE=$(
    aws ecs update-service \
        --cluster ${ecsClusterName} \
        --service ${ecsServiceName} \
        --desired-count 1 \
        --force-new-deployment \
        --region ${region} \
        --task-definition ${latestTaskArn} \
        --query 'service.[desiredCount]' \
        --output text
)
echo "desiredCount = "${UPDATE_FARGATE}

# remove the active tasks
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
${SCRIPT_DIR}/remove_tasks.sh