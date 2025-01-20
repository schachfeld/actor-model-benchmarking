defmodule JsonInterpreter do
  defp parse_json(json) do
    parsed = JSON.decode!(json)
    IO.puts("parsed: #{parsed}")
  end

  def start() do
    receive do
      {:json, json, index} -> parse_json(json)
    end

    start()
  end
end

defmodule FileReader do
  def start() do
    pid = spawn(JsonInterpreter, :start, [])

    File.stream!("../messages_short.log")
    |> Stream.with_index()
    |> Stream.map(fn {line, index} ->
      send(pid, {:json, line, index})
    end)
    # |> Stream.map(fn {line, index} -> IO.puts("index: #{index}") end)
    |> Stream.run()
  end
end

FileReader.start()
