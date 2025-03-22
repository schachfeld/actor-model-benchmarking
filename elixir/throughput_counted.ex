defmodule Bench do
  def receiver() do
    receive do
      {:ok, _index} -> nil
    end

    receiver()
  end

  def sender(index, pid) do
    if index == 0 do
      IO.puts("Done")
      nil
    else
      send(pid, {:ok, index})
      sender(index - 1, pid)
    end
  end
end

pid = spawn(Bench, :receiver, [])

n = 10_000_000

start_time = :os.system_time(:nanosecond)
IO.puts("Start")
Bench.sender(n, pid)
end_time = :os.system_time(:nanosecond)

elapsed = end_time - start_time
elapsedSeconds = elapsed / 1_000_000

IO.puts("received #{n} messages. #{n / elapsedSeconds} msg/sec")
IO.puts("Total time: #{elapsedSeconds}")

File.write!("throughput_bench_results/throughput_#{n}.txt", "#{elapsed},", [:append])

Process.exit(pid, :kill)
