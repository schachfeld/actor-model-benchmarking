defmodule Bench do
  def try_hashes() do
    receive do
      {:ok} -> nil
    end

    try_hashes()
  end

  def master(word) do
    pid = spawn(Bench, :try_hashes, [])
  end
end

pid = spawn(Bench, :master, [])

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
