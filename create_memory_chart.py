import re

import matplotlib.pyplot as plt
import numpy as np

data = []

with open("mem_beam.log") as f:    
    for line in f:
        data += re.sub(' +', ' ',line.strip()).split(" ")



data = [int(x) for x in data]


vsz = data[1::3]
mem = data[2::3]

fig, ax1 = plt.subplots()

vsz_mb = [x / 1024 for x in vsz]
mem_mb = [x / 1024 for x in mem]

time = np.arange(0, len(vsz_mb) * 0.5, 0.5)

ax1.fill_between(time, vsz_mb, color='blue', alpha=0.5, label='Virtual Memory')
ax1.set_xlabel('Time (seconds)')
ax1.set_ylabel('Memory (MB)', color='b')

ax1.fill_between(time, mem_mb, color='red', alpha=0.5, label='Real Memory')

plt.title('Memory Usage')

ax1.set_ylim(bottom=0)

ax1.legend(loc='upper left')

plt.savefig('memory_usage.png')



