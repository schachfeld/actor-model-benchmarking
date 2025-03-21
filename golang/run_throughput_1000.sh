#!/bin/bash

for i in {1..1000}
do
    echo "Iteration: $i"
    chrt -f 99 /usr/bin/time -v -o bench.txt /home/valentin/projects/bachelorarbeit/actor-model-benchmarking/golang/main
done