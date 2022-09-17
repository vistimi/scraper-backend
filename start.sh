ecsClusterName=scraperClusterECS
ecsServiceName=scraperServiceFargate
nameTaskDefinition=scraperDefinitionFargate
region=us-east-1
applicationLoadBalancer=scraper-alb
targetGroup=scraper-fargate-tg
securityGroup=scraper-alb-sg

# get security group
GET_SGs=$(
    aws ec2 describe-security-groups \
        --filters Name=group-name,Values=${securityGroup} \
        --region ${region} \
        --query 'SecurityGroups[*].[GroupId]' \
        --output text
    )
# for one variable only
securityGroupID=${GET_SGs}
echo "securityGroupID = ${securityGroupID}"

# create elb
CREATE_ELB_GET_ARNs=$(
    aws elbv2 create-load-balancer \
        --name ${applicationLoadBalancer} \
        --region ${region} \
        --subnets subnet-0f537719e61603532 subnet-0e7d325654b95d975 subnet-04872fb36cef128ac \
        --security-groups ${securityGroupID} \
        --type application \
        --scheme internet-facing \
        --ip-address-type ipv4 \
        --tags Key=vpc,Value=scraper-vpc \
        --query 'LoadBalancers[*].[LoadBalancerArn]' \
        --output text
    )
loadBalancerArn=${CREATE_ELB_GET_ARNs}
echo "loadBalancerArn = ${loadBalancerArn}"

# get target group of elb
GET_TGs=$(
    aws elbv2 describe-target-groups \
        --names ${targetGroup} \
        --region ${region} \
        --query 'TargetGroups[*].[TargetGroupArn]' \
        --output text
)
targetGroupArn=${GET_TGs}
echo "targetGroupArn = ${targetGroupArn}"

# add listener to elb
CREATE_LISTENER=$(
    aws elbv2 create-listener \
        --load-balancer-arn ${loadBalancerArn} \
        --protocol HTTP \
        --port 80 \
        --region ${region} \
        --default-actions Type=forward,TargetGroupArn=${targetGroupArn} \
        --query 'Listeners[*].[ListenerArn]' \
        --output text
)
echo "listenerArn = ${CREATE_LISTENER}"

# desired tasks
UPDATE_FARGATE=$(
    aws ecs update-service \
        --cluster ${ecsClusterName} \
        --service ${ecsServiceName} \
        --desired-count 1 \
        --force-new-deployment \
        --region ${region} \
        --query 'service[*].[desiredCount]' \
        --output text
)
echo "desiredCount = ${UPDATE_FARGATE}"

# # for more than one variable
# IFS=', ' read -r -a array <<< "$ELB_ARN" 
# for arn in "${array[@]}" 
# do
#     echo ${arn} 
# done