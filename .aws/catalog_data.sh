#!/bin/bash

for ARGUMENT in "$@"
do
   KEY=$(echo $ARGUMENT | cut -f1 -d=)

   KEY_LENGTH=${#KEY}
   VALUE="${ARGUMENT:$KEY_LENGTH+1}"

   export "$KEY"="$VALUE"
done

# https://docs.aws.amazon.com/AmazonECR/latest/public/getting-started-cli.html
cat << EOF > ${FILE_NAME}
{
    "architectures": [
        "x86"
    ],
    "operatingSystems": [
        "Linux"
    ]
}
EOF

cat ${FILE_NAME}