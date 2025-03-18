defmodule Bench do
  def actor() do
    receive do
      {:ok} -> nil
    end

    actor()
  end

  def benchLatency(index) do
    if index == 0 do
      nil
    else
      benchLatency(index - 1)

      start_time = :os.system_time(:nanosecond)
      spawn(Bench, :actor, [])
      end_time = :os.system_time(:nanosecond)
      elapsed = end_time - start_time
      IO.write("#{elapsed},")
    end
  end
end

n = 1_000_000

Bench.benchLatency(n)

Process.exit(pid, :kill)
