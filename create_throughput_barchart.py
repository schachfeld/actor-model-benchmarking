import matplotlib.pyplot as plt
import numpy as np

with open("elixir/throughput_bench_results/throughput_10000000.txt") as f:    
    elixirdata = f.read().split(",")
elixirdata = [int(x) for x in elixirdata]


with open("golang/throughput_bench_results/throughput_10000000.txt") as f:
    godata = f.read().split(",")

godata = [int(x) for x in godata]

with open("akka//throughput_bench_results/throughput_10000000.txt") as f:
    akkadata = f.read().split(",")
akkadata = [int(x) for x in akkadata]


fig, ax = plt.subplots()

elixir_msg_per_sec = [1e7 / (x / 1e9) for x in elixirdata]
go_msg_per_sec = [1e7 / (x / 1e9)  for x in godata]
akka_msg_per_sec = [1e7 / (x / 1e9) for x in akkadata]

means = [np.mean(elixir_msg_per_sec), np.mean(go_msg_per_sec), np.mean(akka_msg_per_sec)]
tick_labels = ['Elixir', 'Go', 'Akka']

ax.bar(tick_labels, means, color=[0.53, 0.25, 0.31])
ax.set_ylabel('Messages per Second')
ax.set_title('Mean Throughput of 10 Million Messages')
plt.savefig('images/throughput/throughput_10mil_barchart.svg')
plt.savefig('images/throughput/throughput_10mil_barchart.png', dpi=300)
plt.savefig('images/throughput/throughput_10mil_barchart.pdf')


fig, ax2 = plt.subplots()
ax2.boxplot([elixir_msg_per_sec, go_msg_per_sec, akka_msg_per_sec], tick_labels=tick_labels, meanline=True, showmeans=True)
ax2.set_ylabel('Messages per Second')
ax2.set_title('Throughput Distribution of 10 Million Messages')

plt.savefig('images/throughput/throughput_10mil_boxplot.svg')
plt.savefig('images/throughput/throughput_10mil_boxplot.png', dpi=300)
plt.savefig('images/throughput/throughput_10mil_boxplot.pdf')

