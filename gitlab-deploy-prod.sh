# !/bin/bash

# Get servers list:
set — f
# Variables from GitLab server:
string=$PROD_DEPLOY_SERVER
# Note: They can’t have spaces!!
array=(${string//,/ })
# Iterate servers for deploy and pull last commit
# Careful with the ; https://stackoverflow.com/a/20666248/1057052
for i in '${!array[@]}'; do
  echo 'Deploy project on server ${array[i]}'
ssh ubuntu@${array[i]} 'cd account-api && git pull origin master'

done
