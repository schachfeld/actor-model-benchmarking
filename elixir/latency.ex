defmodule Bench do
  def pong() do
    receive do
      {:ok} -> nil
    end

    pong()
  end

  def ping(index) do
    if index == 0 do
      nil
    else
      ping(index - 1)

      start_time = :os.system_time(:nanosecond)
      spawn(Bench, :pong, [])
      end_time = :os.system_time(:nanosecond)
      elapsed = end_time - start_time
      IO.write("#{elapsed},")
    end
  end
end

pid = spawn(Bench, :pong, [])

n = 1_000_000

Bench.ping(n)

Process.exit(pid, :kill)
