#!/bin/bash

trap "exit" INT

    for i in {74..100}
    do
        echo "Workers: $arg, Iteration: $i"
        export TOTAL_WORKERS=1
        /usr/bin/time -v -o prime_time_${arg}.txt elixir prime_number_check.ex
    done

for arg in 2 5 10 20 50 100 1000
do
    for i in {1..100}
    do
        echo "Workers: $arg, Iteration: $i"
        export TOTAL_WORKERS=$arg
        /usr/bin/time -v -o prime_time_${arg}.txt elixir prime_number_check.ex
    done
done

# for i in {1..100}
#     do
#         echo "Argument: $arg, Iteration: $i"
#         /usr/bin/time -v -o prime_time_10.txt java -jar target/scala-3.3.4/akka-assembly-0.1.0-SNAPSHOT.jar 10
#     done