#!/bin/bash

for arg in 1 2 5 10 20 50 100 1000
do
    for i in {1..100}
    do
        echo "Argument: $arg, Iteration: $i"
        chrt -f 99 /usr/bin/time -v -o prime_time_${arg}.txt java -jar target/scala-3.3.4/akka-assembly-0.1.0-SNAPSHOT.jar $arg
    done
done
