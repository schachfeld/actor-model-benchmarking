import matplotlib.pyplot as plt
import numpy as np

data = []

with open("elixir/latency.txt") as f:    
    for line in f:
        data = line.split(",")


data = [int(x) for x in data]

fig, ax = plt.subplots()

ax.plot(data, label='Latency')

median = np.median(data)
ax.axhline(median, color='r', linestyle='--', label=f'Median: {median:.2f}')

average = np.mean(data)
ax.axhline(average, color='b', linestyle='-', label=f'Average: {average:.2f}')

ax.set_title('Latency over Time')
ax.set_xlabel('Sample Number')
ax.set_ylabel('Latency (ns)')

ax.grid(True)

ax.legend()

plt.savefig('latency.png')