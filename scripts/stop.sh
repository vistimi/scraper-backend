#!/bin/bash
ecsClusterName=scraperClusterECS
ecsServiceName=scraperServiceFargate
region=us-east-1
applicationLoadBalancer=scraper-alb

# get elb arn
GET_ELBs=$(
    aws elbv2 describe-load-balancers \
        --names ${applicationLoadBalancer} \
        --region ${region} \
        --query 'LoadBalancers[0].[LoadBalancerArn]' \
        --output text
)
loadBalancerArn=${GET_ELBs}
echo "loadBalancerArn = "${loadBalancerArn}

# delete load balancer
aws elbv2 delete-load-balancer \
        --load-balancer-arn ${loadBalancerArn} \
        --region ${region}

# desired count to zero
UPDATE_FARGATE=$(
    aws ecs update-service \
        --cluster ${ecsClusterName} \
        --service ${ecsServiceName} \
        --desired-count 0 \
        --force-new-deployment \
        --region ${region} \
        --query 'service.[desiredCount]' \
        --output text
)
echo "desiredCount = "${UPDATE_FARGATE}

# remove the active tasks
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
${SCRIPT_DIR}/remove_tasks.sh