defmodule Bench do
  def pong() do
    receive do
      {:ok} -> nil
    end

    pong()
  end

  def ping(index, pid) do
    if index == 0 do
      IO.puts("Done")
      Process.exit(pid, :kill)
    end

    send(pid, {:ok})
    ping(index - 1, pid)
  end
end

pid = spawn(Bench, :pong, [])

Bench.ping(100_000_000, pid)
