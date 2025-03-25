#!/bin/bash

trap "exit" INT


for i in {1..100}
    do
        echo "Iteration: $i"
        /usr/bin/time -v -o results/cb_analyzer_long_time.txt mix run
    done