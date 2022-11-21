#!/bin/bash
gh workflow run cicd.yml --repo KookaS/scraper-backend --ref production\
    -f aws-account-name="KookaS" \
    -f aws-account-id="401582117818" \
    -f aws-region="us-east-1" \
    -f environment-name="test" \
    -f container-cpu="256" \
    -f container-memory="512" \
    -f container-memory-reservation="500" \
    -f aws-exec-role="Role_ECS_S3" \
    -f aws-task-role="Role_ECS_S3" \
    -f keep-images-amount="1" 
    # -f keep-images-days="1"

echo "Sleep 10 seconds for spawning action"
sleep 10s
echo "Continue to check the status"

# while workflow status == in_progress, wait
workflowStatus=$(gh run list --workflow CI/CD --limit 1 | awk '{print $1}')
while [ "${workflowStatus}" != "completed" ]
do
    echo "Waiting for status workflow to complete: "${workflowStatus}
    sleep 5s
    workflowStatus=$(gh run list --workflow CI/CD --limit 1 | awk '{print $1}')
done

echo "Workflow finished: "${workflowStatus}