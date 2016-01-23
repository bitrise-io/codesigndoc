#!/bin/bash
set -x

temp_dir="$(mktemp -d -t codesigndoc)"

cd $temp_dir

curl -fL https://github.com/bitrise-tools/codesigndoc/releases/download/0.9.1/codesigndoc-Darwin-x86_64 > ./codesigndoc
chmod +x ./codesigndoc
./codesigndoc scan
