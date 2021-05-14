#!/bin/bash

set -e -o pipefail

# TODO: trigger this by some tag within the file we could grep for
# get the list of "models.go" files current staged
model_files_updated=$(git diff --name-only --cached | grep "models.go" | wc -l)

# if there are no "models.go", then abort
if [ $model_files_updated == 0 ]; then
  echo "No models.go files were updated, no action needed"
  exit 0
fi

# run model generation for browser and web portal
make export-ext-model
make export-portal-model

# do we have changes in models.ts in portal and browser-ext?
model_ts_files_changed=($(git diff --name-only | grep -E '(browser-extension/src/models\.ts|web-portal/src/models\.ts)'))


# if yes, add them and run build to verify changed didn't break
if [ ${#model_ts_files_changed[@]} == 0 ]; then
  echo "Model changes didn't result in TS file changes, no action needed"
  exit 0
fi

# run builds
root_dir=$(git rev-parse --show-toplevel)
## now loop through the above array
for i in "${model_ts_files_changed[@]}"
do
  project_path=$(echo $i | cut -d'/' -f1)
  cd "${root_dir}/${project_path}" && npm run build
  relative_path=$(realpath --relative-to="$PWD" "${root_dir}/${i}")
  git add $relative_path
done
