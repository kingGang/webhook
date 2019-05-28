#! /bin/bash

/root/webhook/webhook -path="/root/webhook/gitpull.sh" -port=10065 &

cd /root/caredaily.doc/ && \
hugo server --bind "0.0.0.0" --baseURL "http://192.168.100.245:1313" -t hyde