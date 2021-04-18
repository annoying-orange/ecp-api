#!/bin/bash
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 989041579659.dkr.ecr.ap-northeast-1.amazonaws.com
docker build -t annoying-orange/ecp-api .
docker tag annoying-orange/ecp-api:latest 989041579659.dkr.ecr.ap-northeast-1.amazonaws.com/wesport/ecp-api:latest
docker push 989041579659.dkr.ecr.ap-northeast-1.amazonaws.com/wesport/ecp-api:latest