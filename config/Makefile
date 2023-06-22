SHELL:= /bin/bash
.SILENT:
MAKEFLAGS += --warn-undefined-variables

export GH_ORG GH_REPO GH_PATH GH_BRANCH OVERRIDE_EXTENSION OUTPUT_FOLDER

gh-load-folder:
	echo "curl -L \
		-H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer ${GITHUB_TOKEN}" \
		-H "X-GitHub-Api-Version: 2022-11-28" \
		https://api.github.com/repos/${GH_ORG}/${GH_REPO}/contents/${GH_PATH}?ref=${GH_BRANCH}"
	$(eval res=$(shell curl -L \
		-H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer ${GITHUB_TOKEN}" \
		-H "X-GitHub-Api-Version: 2022-11-28" \
		https://api.github.com/repos/${GH_ORG}/${GH_REPO}/contents/${GH_PATH}?ref=${GH_BRANCH} | jq -c '.[] | .path'))
	for file in ${res}; do \
		make -f ${MAKEFILE} gh-load-file OUTPUT_FOLDER=${OUTPUT_FOLDER} GH_PATH="$$file"; \
    done
gh-load-file:
	echo ${GH_PATH}
	curl -L -o ${OUTPUT_FOLDER}/$(shell basename ${GH_PATH} | cut -d. -f1)_${OVERRIDE_EXTENSION}$(shell [[ "${GH_PATH}" = *.* ]] && echo .$(shell basename ${GH_PATH} | cut -d. -f2) || echo '') \
			-H "Accept: application/vnd.github.v3.raw" \
			-H "Authorization: Bearer ${GITHUB_TOKEN}" \
			-H "X-GitHub-Api-Version: 2022-11-28" \
			https://api.github.com/repos/${GH_ORG}/${GH_REPO}/contents/${GH_PATH}?ref=${GH_BRANCH}

.ONESHELL: prepare-scraper-backend
prepare-scraper-backend:
	$(eval GH_ORG=KookaS)
	$(eval GH_REPO=scraper-backend)
	cat <<-EOF > ${OUTPUT_FOLDER}/${OVERRIDE_EXTENSION}.env
		COMMON_NAME=${COMMON_NAME}
		CLOUD_HOST=${CLOUD_HOST}
		FLICKR_PRIVATE_KEY=${FLICKR_PRIVATE_KEY}
		FLICKR_PUBLIC_KEY=${FLICKR_PUBLIC_KEY}
		UNSPLASH_PRIVATE_KEY=${UNSPLASH_PRIVATE_KEY}
		UNSPLASH_PUBLIC_KEY=${UNSPLASH_PUBLIC_KEY}
		PEXELS_PUBLIC_KEY=${PEXELS_PUBLIC_KEY}
		AWS_REGION=${AWS_REGION}
		AWS_ACCESS_KEY=${AWS_ACCESS_KEY}
		AWS_SECRET_KEY=${AWS_SECRET_KEY}
	EOF
	make -f ${MAKEFILE} gh-load-folder GH_PATH=config