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
           "description": "keep last ${KEEP_IMAGES_AMOUNT}",
           "selection": {
               "tagStatus": "any",
               "countType": "imageCountMoreThan",
               "countNumber": ${KEEP_IMAGES_AMOUNT}
           },
           "action": {
               "type": "expire"
           }
       }
   ]
}
EOF

cat ${FILE_NAME}