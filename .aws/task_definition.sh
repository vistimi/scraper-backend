#!/bin/bash

for ARGUMENT in "$@"
do
   KEY=$(echo $ARGUMENT | cut -f1 -d=)

   KEY_LENGTH=${#KEY}
   VALUE="${ARGUMENT:$KEY_LENGTH+1}"

   export "$KEY"="$VALUE"
done

# "containerDefinitions": [
#     {
#         "logConfiguration": {
#             "logDriver": "awslogs",
#             "secretOptions": null,
#             "options": {
#             "awslogs-group": "/ecs/scraperDefinitionFargate",
#             "awslogs-region": "us-east-1",
#             "awslogs-stream-prefix": "ecs"
#             }
#         },
#     }
# ]

cat << EOF > ${FILE_NAME}
{
  "ipcMode": null,
  "executionRoleArn": "arn:aws:iam::${AWS_ACCOUNT_ID}:role/${AWS_EXEC_ROLE}",
  "containerDefinitions": [
    {
      "dnsSearchDomains": null,
      "environmentFiles": [
        {
          "value": "arn:aws:s3:::${AWS_BUCKET_ENV_NAME}/${AWS_FILE_ENV_NAME}",
          "type": "s3"
        }
      ],
      "logConfiguration": null,
      "entryPoint": [],
      "portMappings": [
        {
          "hostPort": 8080,
          "protocol": "tcp",
          "containerPort": 8080
        },
        {
          "hostPort": 27017,
          "protocol": "tcp",
          "containerPort": 27017
        }
      ],
      "command": [],
      "linuxParameters": null,
      "cpu": ${CPU},
      "environment": [],
      "resourceRequirements": null,
      "ulimits": null,
      "dnsServers": null,
      "mountPoints": [],
      "workingDirectory": null,
      "secrets": null,
      "dockerSecurityOptions": null,
      "memory": ${MEMORY},
      "memoryReservation": ${MEMORY_RESERVATION},
      "volumesFrom": [],
      "stopTimeout": null,
      "image": "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_NAME}:8e7e16a6e89faf1d95514a5afaa1798b7263e3a7",
      "startTimeout": null,
      "firelensConfiguration": null,
      "dependsOn": null,
      "disableNetworking": null,
      "interactive": null,
      "healthCheck": null,
      "essential": true,
      "links": null,
      "hostname": null,
      "extraHosts": null,
      "pseudoTerminal": null,
      "user": null,
      "readonlyRootFilesystem": null,
      "dockerLabels": null,
      "systemControls": null,
      "privileged": null,
      "name": ${NAME}
    }
  ],
  "placementConstraints": [],
  "memory": ${MEMORY},
  "taskRoleArn": "arn:aws:iam::${AWS_ACCOUNT_ID}:role/${AWS_TASK_ROLE}",
  "family": ${DEFINITION_FAMILY},
  "pidMode": null,
  "requiresCompatibilities": [
    "EC2"
  ],
  "networkMode": null,
  "runtimePlatform": null,
  "cpu": ${CPU},
  "inferenceAccelerators": null,
  "proxyConfiguration": null,
  "volumes": []
}
EOF

cat ${FILE_NAME}