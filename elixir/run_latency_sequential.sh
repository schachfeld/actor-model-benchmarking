#!/bin/bash

# filepath: /home/valentin/projects/bachelorarbeit/actor-model-benchmarking/elixir/run_benchmark.sh

for i in {1..100}
do
  elixir latency.ex > "latency_results/output_$i.txt"
done