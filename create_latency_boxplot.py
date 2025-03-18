import matplotlib.pyplot as plt
import numpy as np

with open("elixir/latency.txt") as f:    
    for line in f:
        elixirdata = line.split(",")


elixirdata = [int(x) for x in elixirdata]


with open("golang/latency.txt") as f:
    godata = f.read().split(",")

godata = [int(x) for x in godata]


fig, ax = plt.subplots()

ax.boxplot([elixirdata, godata], tick_labels=['Elixir', 'Go'], showfliers=False, showmeans=True, meanline=True)
ax.set_ylabel('Latency (nanoseconds)')
ax.set_title('Actor Creation Latency') 

# lower is better text
# ax.text(1.5, ax.get_ylim()[0] - (ax.get_ylim()[1] - ax.get_ylim()[0]) * 0.1, 'Lower is better', ha='center', va='center')


mean_patch = plt.Line2D([], [], color='green', label='Mean', linestyle='dashed', linewidth=1)

median_patch = plt.Line2D([], [], color='orange', label='Median', linewidth=1)

ax.legend(handles=[mean_patch, median_patch])

plt.savefig('images/latency_boxplot.png')