import matplotlib.pyplot as plt
import numpy as np

# with open("elixir/throughput_bench_results/throughput_10000000.txt") as f:    
#     elixirdata = f.read().split(",")
# elixirdata = [int(x) for x in elixirdata]


# with open("golang/throughput_bench_results/throughput_10000000.txt") as f:
#     godata = f.read().split(",")

# godata = [int(x) for x in godata]

run = "run_2"

workers_list = [1, 2, 5, 10, 20, 50, 100, 1000]
akkadata = []

for workers in workers_list:
    with open(f"akka/prime_bench_results/{run}/10mil_{workers}workers.txt") as f:
        data = f.read().split(",")
    data = [int(x) for x in data]
    akkadata.append(np.array(data) / 1e9)

akkadata1, akkadata2, akkadata5, akkadata10, akkadata20, akkadata50, akkadata100, akkadata1000 = akkadata


godata = []

for workers in workers_list:
    with open(f"golang/prime_bench_results/10mil_{workers}workers.txt") as f:
        data = f.read().split(",")
    data = [int(x) for x in data]
    godata.append(np.array(data) / 1e9)

godata1, godata2, godata5, godata10, godata20, godata50, godata100, godata1000 = godata

elixirdata = []

for workers in workers_list:
    with open(f"elixir/prime_bench_results/10mil_{workers}workers.txt") as f:
        data = f.read().split(",")
    data = [int(x) for x in data]
    elixirdata.append(np.array(data) / 1e9)

elixirdata1, elixirdata2, elixirdata5, elixirdata10, elixirdata20, elixirdata50, elixirdata100, elixirdata1000 = elixirdata

# Create a chart
workers = [1, 2, 5, 10, 20, 50, 100, 1000]

akka_times = [np.mean(akkadata1), np.mean(akkadata2), np.mean(akkadata5), np.mean(akkadata10), np.mean(akkadata20), np.mean(akkadata50), np.mean(akkadata100), np.mean(akkadata1000)]
go_times = [np.mean(godata1), np.mean(godata2), np.mean(godata5), np.mean(godata10), np.mean(godata20), np.mean(godata50), np.mean(godata100), np.mean(godata1000)]
elixir_times = [np.mean(elixirdata1), np.mean(elixirdata2), np.mean(elixirdata5), np.mean(elixirdata10), np.mean(elixirdata20), np.mean(elixirdata50), np.mean(elixirdata100), np.mean(elixirdata1000)]

plt.figure(figsize=(10, 6))
plt.plot(workers, akka_times, marker='o', label='Akka')
plt.plot(workers, go_times, marker='o', label='Golang')
plt.plot(workers, elixir_times, marker='o', label='Elixir')
plt.xscale('log')
plt.xlabel('Number of Workers')
plt.ylabel('Time (seconds)')
plt.title('Prime Calculation Time vs Number of Workers')
plt.legend()
plt.grid(True)

plt.tight_layout(pad=0)
plt.savefig('images/prime/prime_chart.svg')
plt.savefig('images/prime/prime_chart.png', dpi=300)
plt.savefig('images/prime/prime_chart.pdf')

