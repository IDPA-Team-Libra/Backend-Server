#!/bin/bash
set -e
bash install.sh
sudo docker-compose build server
sudo docker-compose build db
sudo docker-compose build frontend
rm main_linux_amd64
