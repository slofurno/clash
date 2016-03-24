#!/bin/bash
t1=$1
t2=$2

diff -w <(node $t1) <(cat $t2)
