defmodule Bench do
  def receiver() do
    receive do
      {:ok} -> nil
    end

    receiver()
  end

  def sender(index, pid) do
    if index == 0 do
      IO.puts("Done")
      nil
    else
      send(pid, {:ok})
      sender(index - 1, pid)
    end
  end
end

pid = spawn(Bench, :receiver, [])

n = 10_000_000

start_time = :os.system_time(:millisecond)
IO.puts("Start")
Bench.sender(n, pid)
end_time = :os.system_time(:millisecond)

elapsed = end_time - start_time
elapsedSeconds = elapsed / 1000

IO.puts("received #{n} messages. #{n / elapsedSeconds} msg/sec")
IO.puts("Total time: #{elapsedSeconds}")

Process.exit(pid, :kill)
