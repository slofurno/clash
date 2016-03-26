#!/bin/bash

#inotifywait 
pid=0
while :
do
files=$(find . -name "*.go" | tr '\n' ' ')
inotifywait -e modify $files
if [ $pid -ne 0 ]; then kill $pid; fi
sleep 1
(go build && ./front) & pid=$!
echo "pid ${pid}"
done
