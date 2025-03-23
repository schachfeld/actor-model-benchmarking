#!/bin/bash

trap "exit" INT

for i in {1..100}
do
    echo "Iteration: $i"
    /usr/bin/time -v -o bench_data/time.txt ./main
done