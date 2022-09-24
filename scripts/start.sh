#!/bin/bash
region=us-east-1
applicationLoadBalancer=scraper-alb
targetGroup=scraper-fargate-tg
securityGroup=scraper-alb-sg

# get security group
GET_SGs=$(
    aws ec2 describe-security-groups \
        --filters Name=group-name,Values=${securityGroup} \
        --region ${region} \
        --query 'SecurityGroups[0].[GroupId]' \
        --output text
    )
securityGroupID=${GET_SGs}
echo "securityGroupID = "${securityGroupID}

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
        --query 'LoadBalancers[0].[LoadBalancerArn]' \
        --output text
    )
loadBalancerArn=${CREATE_ELB_GET_ARNs}
echo "loadBalancerArn = "${loadBalancerArn}

# get the DNS for browsing
GET_ELB_DNSs=$(
    aws elbv2 describe-load-balancers \
        --names ${applicationLoadBalancer} \
        --region ${region} \
        --query 'LoadBalancers[0].[DNSName]' \
        --output text
)
elbDNS=${GET_ELB_DNSs}
echo "elbDNS = ${elbDNS}"

# get target group of elb
GET_TGs=$(
    aws elbv2 describe-target-groups \
        --names ${targetGroup} \
        --region ${region} \
        --query 'TargetGroups[0].[TargetGroupArn]' \
        --output text
)
targetGroupArn=${GET_TGs}
echo "targetGroupArn = "${targetGroupArn}

# add listener to elb
CREATE_LISTENER=$(
    aws elbv2 create-listener \
        --load-balancer-arn ${loadBalancerArn} \
        --protocol HTTP \
        --port 80 \
        --region ${region} \
        --default-actions Type=forward,TargetGroupArn=${targetGroupArn} \
        --query 'Listeners[0].[ListenerArn]' \
        --output text
)
echo "listenerArn = "${CREATE_LISTENER}

# update the service
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
${SCRIPT_DIR}/update_service.sh