import matplotlib.pyplot as plt
import numpy as np

# with open("elixir/throughput_bench_results/throughput_10000000.txt") as f:    
#     elixirdata = f.read().split(",")
# elixirdata = [int(x) for x in elixirdata]


# with open("golang/throughput_bench_results/throughput_10000000.txt") as f:
#     godata = f.read().split(",")

# godata = [int(x) for x in godata]

run = "run_2"

with open(f"akka/prime_bench_results/{run}/10mil_1workers.txt") as f:
    akkadata1 = f.read().split(",")
akkadata1 = [int(x) for x in akkadata1]

with open(f"akka/prime_bench_results/{run}/10mil_2workers.txt") as f:
    akkadata2 = f.read().split(",")
akkadata2 = [int(x) for x in akkadata2]

with open(f"akka/prime_bench_results/{run}/10mil_5workers.txt") as f:
    akkadata5 = f.read().split(",")
akkadata5 = [int(x) for x in akkadata5]

with open(f"akka/prime_bench_results/{run}/10mil_10workers.txt") as f:
    akkadata10 = f.read().split(",")
akkadata10 = [int(x) for x in akkadata10]

with open(f"akka/prime_bench_results/{run}/10mil_20workers.txt") as f:
    akkadata20 = f.read().split(",")
akkadata20 = [int(x) for x in akkadata20]

with open(f"akka/prime_bench_results/{run}/10mil_50workers.txt") as f:
    akkadata50 = f.read().split(",")
akkadata50 = [int(x) for x in akkadata50]

with open(f"akka/prime_bench_results/{run}/10mil_100workers.txt") as f:
    akkadata100 = f.read().split(",")
akkadata100 = [int(x) for x in akkadata100]

with open(f"akka/prime_bench_results/{run}/10mil_1000workers.txt") as f:
    akkadata1000 = f.read().split(",")
akkadata1000 = [int(x) for x in akkadata1000]


# Convert data points from nanoseconds to seconds
akkadata1 = np.array(akkadata1) / 1e9
akkadata2 = np.array(akkadata2) / 1e9
akkadata5 = np.array(akkadata5) / 1e9
akkadata10 = np.array(akkadata10) / 1e9
akkadata20 = np.array(akkadata20) / 1e9
akkadata50 = np.array(akkadata50) / 1e9
akkadata100 = np.array(akkadata100) / 1e9
akkadata1000 = np.array(akkadata1000) / 1e9

# Create a chart
workers = [1, 2, 5, 10, 20, 50, 100, 1000]
times = [np.mean(akkadata1), np.mean(akkadata2), np.mean(akkadata5), np.mean(akkadata10), np.mean(akkadata20), np.mean(akkadata50), np.mean(akkadata100), np.mean(akkadata1000)]

plt.figure(figsize=(10, 6))
plt.plot(workers, times, marker='o')
plt.xscale('log')
plt.xlabel('Number of Workers')
plt.ylabel('Time (seconds)')
plt.title('Prime Calculation Time vs Number of Workers')
plt.grid(True)

plt.savefig('images/prime/prime_chart.svg')
plt.savefig('images/prime/prime_chart.png', dpi=300)
plt.savefig('images/prime/prime_chart.pdf')

