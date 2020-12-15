#!/bin/sh

####################################################
# should run this file like ... 
# ./load_data.sh | jq . > camas_config.json
####################################################


# print pre-part of template
cat pre-camas-config.json.template

# token hunter
curl -s https://gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/token-hunter/-/raw/master/regexes.json | jq 'to_entries[] | {"name": .key, "content": .value, "type": "regex", "noiseLevel": 1, "analysisLayer": "contents", "source": "token-hunter"}' | sed 's/^}/},/'

# rusty hog
curl -s https://raw.githubusercontent.com/newrelic/rusty-hog/master/src/default_rules.json | jq 'to_entries[] | {"name": .key, "content": .value, "type": "regex", "noiseLevel": 1, "analysisLayer": "contents", "source": "rusty-hog"} | select(.name | startswith("Generic") | not)' | sed 's/^}/},/'

# shhgit
curl -s https://raw.githubusercontent.com/eth0izzle/shhgit/master/config.yaml > tmp_shhgit_config.yaml
# 2 separate runs through the doc - 1 to catch "match" and 1 to catch "regex"
docker run --rm -v "${PWD}":/workdir mikefarah/yq yq r -j tmp_shhgit_config.yaml | jq '.signatures[] | select(.match != null) | {"name": .name, "content": .match, "type": "simple", "noiseLevel": 1, "analysisLayer": .part, "source": "shhgit"}' | sed 's/^}/},/'
docker run --rm -v "${PWD}":/workdir mikefarah/yq yq r -j tmp_shhgit_config.yaml | jq '.signatures[] | select(.regex != null) | {"name": .name, "content": .regex, "type": "regex", "noiseLevel": 1, "analysisLayer": .part, "source": "shhgit"}' | sed 's/^}/},/'
rm tmp_shhgit_config.yaml

# repo security scanner
# The ___ sed '$ s/.$//' ___ bit at the end trims the very last char (the last , that was added) -> this needs to be done on the LAST set of rules pulled in but no others
curl -s https://raw.githubusercontent.com/UKHomeOffice/repo-security-scanner/master/rules/gitrob.json | jq '.[] | {"name": .caption, "content": .pattern, "type":.type, "noiseLevel": 1, "analysisLayer": .part, "source": "repo-security-scanner-gitrob"} | . | if .type =="match" then .type = "simple" else . end' | sed 's/^}/},/' | sed '$ s/.$//' 

# print post-part of template
cat post-camas-config.json.template
