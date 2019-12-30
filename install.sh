#!/bin/bash
set -e
cd ./main
gox -osarch="linux/amd64"
cp main_linux_amd64 ../
cd ../
