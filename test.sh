#!/bin/bash -ex

# ${USER}
# ${PWRD}
# are set from env
# USER=xxx PWD=yyy sh test.sh

go build -o gh-api-cli main.go

echo "test file" > testfile

./gh-api-cli rm-auth -n test_token -u ${USER} -p ${PWRD} || echo "ok, keep going"

NEWTOKEN=`./gh-api-cli add-auth -n test_token -r repo -u ${USER} -p ${PWRD}`
NEWTOKEN=`echo "${NEWTOKEN}" | tail -n 1`

(./gh-api-cli add-auth -n test_token -r repo -u ${USER} -p ${PWRD} && echo "must have failed!") || echo "ok, did not created auth twice"

LIST_AUTHS=`./gh-api-cli list-auth -u ${USER} -p ${PWRD}`
echo "${LIST_AUTHS}" | grep ${NEWTOKEN}

AUTH=`./gh-api-cli get-auth -n test_token`
echo "${AUTH}" | grep ${NEWTOKEN}

(./gh-api-cli get-auth -n nopnop && echo "must fail!") || echo "ok, it failed."

./gh-api-cli rm-release --ver 0.0.1 -o mh-cbon -r test-repo -n test_token || echo "ok, keep going" #ensure 0.0.1 does not exist

CREATE_RELEASE_WITH_AUTH=`./gh-api-cli create-release -n test_token -o mh-cbon -r test-repo --ver 0.0.1`
echo "${CREATE_RELEASE_WITH_AUTH}" | grep '"tag_name": "0.0.1"'

UPLOAD_ASSET_WITH_AUTH=`./gh-api-cli upload-release-asset -n test_token -o mh-cbon -r test-repo --ver 0.0.1 -g testfile`
echo "${UPLOAD_ASSET_WITH_AUTH}" | grep 'Assets uploaded!'

DL_ASSET_WITH_AUTH=`./gh-api-cli dl-assets -n test_token -o mh-cbon -r test-repo --ver 0.0.1 -g testfile --out testfileout`
echo "${DL_ASSET_WITH_AUTH}" | grep 'testfileout'

RM_ASSET_WITH_AUTH=`./gh-api-cli rm-assets -n test_token -o mh-cbon -r test-repo --ver 0.0.1 -g testfile`
echo "${RM_ASSET_WITH_AUTH}" | grep "Removed 'testfile'"

RM_RELEASE_WITH_AUTH=`./gh-api-cli rm-release -n test_token -o mh-cbon -r test-repo --ver 0.0.1`
echo "${RM_RELEASE_WITH_AUTH}" | grep "Release deleted with success!"




CREATE_RELEASE_WITH_TOKEN=`./gh-api-cli create-release -t ${NEWTOKEN} -o mh-cbon -r test-repo --ver 0.0.1`
echo "${CREATE_RELEASE_WITH_TOKEN}" | grep '"tag_name": "0.0.1"'

UPLOAD_ASSET_WITH_TOKEN=`./gh-api-cli upload-release-asset -t ${NEWTOKEN} -o mh-cbon -r test-repo --ver 0.0.1 -g testfile`
echo "${UPLOAD_ASSET_WITH_TOKEN}" | grep 'Assets uploaded!'

DL_ASSET_WITH_TOKEN=`./gh-api-cli dl-assets -t ${NEWTOKEN} -o mh-cbon -r test-repo --ver 0.0.1 -g testfile --out testfileout`
echo "${DL_ASSET_WITH_TOKEN}" | grep 'testfileout'

DL_ASSET_WITH_ANON=`./gh-api-cli dl-assets -o mh-cbon -r test-repo --ver 0.0.1 -g testfile --out testfileanon`
echo "${DL_ASSET_WITH_ANON}" | grep 'testfileanon'

RM_ASSET_WITH_TOKEN=`./gh-api-cli rm-assets -t ${NEWTOKEN} -o mh-cbon -r test-repo --ver 0.0.1 -g testfile`
echo "${RM_ASSET_WITH_TOKEN}" | grep "Removed 'testfile'"

RM_RELEASE_WITH_TOKEN=`./gh-api-cli rm-release -t ${NEWTOKEN} -o mh-cbon -r test-repo --ver 0.0.1`
echo "${RM_RELEASE_WITH_TOKEN}" | grep "Release deleted with success!"


DL_ASSET_WITH_TOKEN_AND_GUESS=`./gh-api-cli dl-assets -t ${NEWTOKEN} --guess --ver 3.0.4 -g gh-api-cli-386.deb --out gh-api-cli-386.deb`
echo "${DL_ASSET_WITH_TOKEN_AND_GUESS}" | grep 'gh-api-cli-386.deb'



DEL=`./gh-api-cli rm-auth -n test_token -u ${USER} -p ${PWRD}`
echo "${DEL}" | grep "Deleted authorization: test_token"



set +ex
rm testfile testfileout testfileanon ./gh-api-cli ./gh-api-cli-386.deb
echo ""
echo "OK, ALL FINE"
