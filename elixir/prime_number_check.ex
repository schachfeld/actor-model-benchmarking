defmodule PrimeCalculator do
  def is_prime?(2), do: true
  def is_prime?(n) when n < 2, do: false

  def is_prime?(n) do
    max_divisor = :math.sqrt(n) |> ceil
    range = 2..max_divisor
    Enum.all?(range, fn x -> rem(n, x) != 0 end)
  end
end

defmodule PrimeWorker do
  def start(worker_id, range, parent_pid) do
    prime_count = 0

    for num <- range do
      if PrimeCalculator.is_prime?(num) do
        prime_count = prime_count + 1
      end
    end

    send(parent_pid, {:result, worker_id, prime_count})
  end
end

defmodule PrimeCoordinator do
  def start(total_workers, range) do
    chunk_size = div(Enum.count(range), total_workers)
    parent_pid = self()

    1..total_workers
    |> Enum.map(fn worker_id ->
      start_value = Enum.at(range, (worker_id - 1) * chunk_size)
      end_value = Enum.at(range, min(worker_id * chunk_size - 1, Enum.count(range) - 1))
      worker_range = start_value..end_value

      spawn(PrimeWorker, :start, [worker_id, worker_range, parent_pid])
    end)

    collect_results(total_workers, 0)
  end

  defp collect_results(0, acc), do: acc

  defp collect_results(remaining, acc) do
    receive do
      {:result, _worker_id, prime_count} ->
        collect_results(remaining - 1, acc + prime_count)
    end
  end
end

defmodule PrimeApp do
  def run do
    range = 1..10_000_000
    total_workers = 10

    IO.puts("Starting prime number calculation with #{total_workers} workers...")

    start_time = :os.system_time(:millisecond)
    prime_count = PrimeCoordinator.start(total_workers, range)
    end_time = :os.system_time(:millisecond)

    IO.puts("Found #{prime_count} prime numbers.")
    IO.puts("Calculation took #{end_time - start_time} milliseconds.")
  end
end

PrimeApp.run()

# start_time = :os.system_time(:millisecond)
# is_prime = PrimeCalculator.is_prime?(9_999_991)
# end_time = :os.system_time(:millisecond)

# is_prime |> IO.inspect()

# IO.puts("Prime check took #{end_time - start_time} milliseconds.")
