defmodule Bench do
  def pong() do
    receive do
      {:ok, _index} -> nil
    end

    pong()
  end

  def ping(index, pid) do
    if index == 0 do
      IO.puts("Done")
      nil
    else
      send(pid, {:ok, index})
      ping(index - 1, pid)
    end
  end
end

pid = spawn(Bench, :pong, [])

n = 100_000_000

start_time = :os.system_time(:millisecond)
IO.puts("Start")
Bench.ping(n, pid)
end_time = :os.system_time(:millisecond)

elapsed = end_time - start_time
elapsedSeconds = elapsed / 1000

IO.puts("received #{n} messages. #{n / elapsedSeconds} msg/sec")
IO.puts("Total time: #{elapsedSeconds}")

Process.exit(pid, :kill)
