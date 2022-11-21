#!/bin/bash

for ARGUMENT in "$@"
do
   KEY=$(echo $ARGUMENT | cut -f1 -d=)

   KEY_LENGTH=${#KEY}
   VALUE="${ARGUMENT:$KEY_LENGTH+1}"

   export "$KEY"="$VALUE"
done

cat << EOF > ${FILE_NAME}
{
   "rules": [
        {
            "rulePriority": 1,
            "description": "Keep last stable images",
            "selection": {
                "tagStatus": "tagged",
                "countType": "imageCountMoreThan",
                "countNumber": ${KEEP_IMAGES_AMOUNT}
            },
            "action": {
                "type": "expire"
            }
        },
       {
           "rulePriority": 2,
           "description": "Delete untagged images",
            "selection": {
                "tagStatus": "untagged",
                "countType": "sinceImagePushed",
                "countUnit": "days",
                "countNumber": ${KEEP_IMAGES_DAYS}
            },
            "action": {
                "type": "expire"
            }
       }
   ]
}
EOF

cat ${FILE_NAME}