ecsClusterName=scraperClusterECS
ecsServiceName=scraperServiceFargate
nameTaskDefinition=scraperDefinitionFargate
region=us-east-1
applicationLoadBalancer=scraper-alb
targetGroup=scraper-fargate-tg
securityGroup=scraper-alb-sg

# get elb arn
GET_ELBs=$(
    aws elbv2 describe-load-balancers \
        --names ${applicationLoadBalancer} \
        --region ${region} \
        --query 'LoadBalancers[*].[LoadBalancerArn]' \
        --output text
)
loadBalancerArn=${GET_ELBs}
echo "loadBalancerArn = ${loadBalancerArn}"

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
        --query 'service[*].[desiredCount]' \
        --output text
)
echo "desiredCount = ${UPDATE_FARGATE}"

# get running tasks
GET_TASKs=$(
    aws ecs list-tasks \
        --cluster ${ecsClusterName} \
        --region ${region} \
        --service-name ${ecsServiceName} \
        --query 'taskArns[*]' \
        --output text
)
taskArn=${GET_TASKs}
echo "taskArn = ${taskArn}"

# stop running tasks
STOP_TASKs=$(
    aws ecs stop-task \
        --cluster ${ecsClusterName} \
        --region ${region} \
        --task ${taskArn} \
        --query 'task[*].desiredStatus' \
        --output text
)
echo "desiredStatus = ${STOP_TASKs}"