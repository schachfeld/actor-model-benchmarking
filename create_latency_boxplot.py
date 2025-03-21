import matplotlib.pyplot as plt
import numpy as np

with open("elixir/latency.txt") as f:    
    for line in f:
        elixirdata = line.split(",")


elixirdata = [int(x) for x in elixirdata]


with open("golang/latency_bench_results/latency.txt") as f:
    godata = f.read().split(",")

godata = [int(x) for x in godata]

with open("akka/latency.txt") as f:
    akkadata = f.read().split(",")
akkadata = [int(x) for x in akkadata]


fig, ax = plt.subplots()

box = ax.boxplot([elixirdata, godata, akkadata], 
                 tick_labels=['Elixir', 'Ergo', 'Akka'], 
                 showfliers=False, 
                 showmeans=True, 
                 meanline=True, 
                 meanprops=dict(linewidth=2))

# Calculate and print the Interquartile Range (IQR)
iqr_elixir = np.percentile(elixirdata, 75) - np.percentile(elixirdata, 25)
iqr_golang = np.percentile(godata, 75) - np.percentile(godata, 25)
iqr_akka = np.percentile(akkadata, 75) - np.percentile(akkadata, 25)

print(f'IQR for Elixir: {iqr_elixir}')
print(f'IQR for Golang: {iqr_golang}')
print(f'IQR for Akka: {iqr_akka}')

# Calculate and print the farthest upper and lower outliers
def find_farthest_outliers(data):
    q1 = np.percentile(data, 25)
    q3 = np.percentile(data, 75)
    iqr = q3 - q1
    lower_bound = q1 - 1.5 * iqr
    upper_bound = q3 + 1.5 * iqr
    lower_outliers = [x for x in data if x < lower_bound]
    upper_outliers = [x for x in data if x > upper_bound]
    farthest_lower_outlier = min(lower_outliers) if lower_outliers else None
    farthest_upper_outlier = max(upper_outliers) if upper_outliers else None
    return farthest_lower_outlier, farthest_upper_outlier

farthest_lower_elixir, farthest_upper_elixir = find_farthest_outliers(elixirdata)
farthest_lower_golang, farthest_upper_golang = find_farthest_outliers(godata)
farthest_lower_akka, farthest_upper_akka = find_farthest_outliers(akkadata)

print(f'Farthest lower outlier for Elixir: {farthest_lower_elixir}')
print(f'Farthest upper outlier for Elixir: {farthest_upper_elixir}')
print(f'Farthest lower outlier for Golang: {farthest_lower_golang}')
print(f'Farthest upper outlier for Golang: {farthest_upper_golang}')
print(f'Farthest lower outlier for Akka: {farthest_lower_akka}')
print(f'Farthest upper outlier for Akka: {farthest_upper_akka}')

means = [np.mean(elixirdata), np.mean(godata), np.mean(akkadata)]
for i, mean in enumerate(means, start=1):
    ax.text(i, mean, f'{mean:.2f}', ha='center', va='bottom', color='green')
ax.set_ylabel('Latency (nanoseconds)')
ax.set_title('Actor Spawn Latency') 

# lower is better text
# ax.text(1.5, ax.get_ylim()[0] - (ax.get_ylim()[1] - ax.get_ylim()[0]) * 0.1, 'Lower is better', ha='center', va='center')


mean_patch = plt.Line2D([], [], color='green', label='Mean', linestyle='dashed', linewidth=1)
median_patch = plt.Line2D([], [], color='orange', label='Median', linewidth=1)

ax.legend(handles=[mean_patch, median_patch])

plt.savefig('images/latency_boxplot.svg')
plt.savefig('images/latency_boxplot.png', dpi=300)
plt.savefig('images/latency_boxplot.pdf')