#!/bin/bash

trap "exit" INT


for arg in 1 2 5 10 20 50 100 1000
do
    for i in {1..100}
    do
        echo "Workers: $arg, Iteration: $i"
        /usr/bin/time -v -o prime_time_${arg}.txt ./main $arg
    done
done

# for i in {1..100}
#     do
#         echo "Argument: $arg, Iteration: $i"
#         /usr/bin/time -v -o prime_time_10.txt java -jar target/scala-3.3.4/akka-assembly-0.1.0-SNAPSHOT.jar 10
#     done