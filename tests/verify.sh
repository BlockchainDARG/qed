#!/bin/bash


# Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#   http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e 

QED="go run ../cmd/cli/qed.go -k path -e http://localhost:8080"

add_event(){
	local event="$1"; shift
	local value="$1"; shift
	$QED  add --key "${event}" --value "${value}"
}


#Adding key [ test event ] with value [ 2 ]
#test event
#Received snapshot with values: 
#	Event: test event
#	HyperDigest: a45fe00356dfccb20b8bc9a7c8331d5c0f89c4e70e43ea0dc0cb646a4b29e59b
#	HistoryDigest: 444f6e7eee66986752983c1d8952e2f0998488a5b038bed013c55528551eaafa
#	Version: 0

verify_event() {
	local commitment="$1"; shift
	echo "${commitment}"
	local event=$(echo "${commitment}" | grep "Event: " | awk -F': ' '{print $2;}')
	local history=$(echo "${commitment}" | grep "HistoryDigest" | awk -F': ' '{print $2;}')
	local hyper=$(echo "${commitment}" | grep "HyperDigest: " | awk -F': ' '{print $2;}')
	local version=$(echo "${commitment}" | grep "Version: " | awk -F': ' '{print $2;}')
	
	
	$QED membership --historyDigest ${history}   --hyperDigest ${hyper}  --version ${version} --key ${event} --verify
}


commitment=$(add_event "$@")
verify_event "${commitment}"
