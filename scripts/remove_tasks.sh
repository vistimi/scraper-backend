#!/bin/bash
ecsClusterName=scraperClusterECS
ecsServiceName=scraperServiceFargate
nameTaskDefinition=scraperDefinitionFargate
region=us-east-1
applicationLoadBalancer=scraper-alb
targetGroup=scraper-fargate-tg
securityGroup=scraper-alb-sg

# get running tasks
GET_TASKs=$(
    aws ecs list-tasks \
        --cluster ${ecsClusterName} \
        --region ${region} \
        --service-name ${ecsServiceName} \
        --query 'taskArns[0]' \
        --output text
)
taskArn=${GET_TASKs}
echo "taskArn = ${taskArn}"

if [ "${taskArn}" = "None" ]; then
  exit 1
fi

# stop running tasks
STOP_TASKs=$(
    aws ecs stop-task \
        --cluster ${ecsClusterName} \
        --region ${region} \
        --task ${taskArn} \
        --query 'task.[desiredStatus]' \
        --output text
)
echo "desiredStatus = ${STOP_TASKs}"