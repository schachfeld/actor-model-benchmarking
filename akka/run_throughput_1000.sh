#!/bin/bash

for i in {1..1000}
do
    echo "Iteration: $i"
    chrt -f 99 /usr/bin/time -v -o throughput_time.txt java -jar target/scala-3.3.4/akka-assembly-0.1.0-SNAPSHOT.jar
done