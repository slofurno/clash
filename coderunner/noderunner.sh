#!/bin/bash
t1=$1
t2=$2

diff <(node $t1) <(node $t2)
