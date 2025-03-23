#!/bin/bash

trap "exit" INT


for i in {1..100}
    do
        echo "Iteration: $i"
        /usr/bin/time -v -o cb_analyzer_results/cb_analyzer_results.txt java -jar target/scala-3.3.4/akka-assembly-0.1.0-SNAPSHOT.jar
    done