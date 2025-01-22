defmodule CBMessage do
  @derive [Poison.Encoder]
  defstruct [
    :channel,
    :client_id,
    :timestamp,
    :sequence_num,
    :events
  ]
end

defmodule CBEvent do
  @derive [Poison.Encoder]
  defstruct [
    :type,
    :product_id,
    :updates
  ]
end

defmodule CBUpdate do
  @derive [Poison.Encoder]
  defstruct [
    :side,
    :event_time,
    :price_level,
    :new_quantity
  ]
end

defmodule FileWriter do
  def writeLine(filename, file \\ nil) do
    file =
      if file == nil do
        File.open!(filename, [:write, :utf8])
      else
        file
      end

    receive do
      {:writeLine, content} ->
        IO.write(file, content)
        IO.write(file, "\n")
    end

    writeLine(filename, file)
  end
end

defmodule AvgOrderBookCalculator do
  defp to_float(value) do
    cond do
      String.contains?(value, ".") -> String.to_float(value)
      true -> String.to_integer(value) * 1.0
    end
  end

  def calculateOrderBookAvg(productId, filewriterPID \\ nil) do
    filewriterPID =
      if filewriterPID == nil do
        spawn(FileWriter, :writeLine, ["orderbooks/" <> productId <> ".txt"])
      else
        filewriterPID
      end

    receive do
      {:cbupdate, updates} when is_list(updates) ->
        if length(updates) > 0 do
          allBids =
            updates
            |> Enum.filter(&(&1.side == "bid"))

          # IO.inspect(updates)

          if(length(allBids) > 0) do
            bidsPrice =
              updates
              |> Enum.reduce(0, fn %CBUpdate{price_level: price_level}, acc ->
                acc + to_float(price_level)
              end)

            avgBids =
              (bidsPrice / length(allBids))
              |> Float.to_string()

            send(filewriterPID, {:writeLine, "Average bids: #{avgBids}"})
          end

          allOffers =
            updates
            |> Enum.filter(&(&1.side == "offer"))

          if length(allOffers) > 0 do
            offersPrice =
              updates
              |> Enum.reduce(0, fn %CBUpdate{price_level: price_level}, acc ->
                acc + to_float(price_level)
              end)

            avgOffers =
              (offersPrice / length(allOffers))
              |> Float.to_string()

            send(filewriterPID, {:writeLine, "Average offers: #{avgOffers}"})
          end
        end
    end

    calculateOrderBookAvg(productId, filewriterPID)
  end
end

defmodule ProductDistributor do
  def start(productId, avgOrderBookPID \\ nil) do
    avgOrderBookPID =
      if avgOrderBookPID == nil do
        spawn(AvgOrderBookCalculator, :calculateOrderBookAvg, [productId])
      else
        avgOrderBookPID
      end

    receive do
      {:cbevent, %CBEvent{} = event} -> send(avgOrderBookPID, {:cbupdate, event.updates})
    end

    start(productId, avgOrderBookPID)
  end
end

defmodule CurrencyDistributor do
  def start(productDistributorPIDs \\ %{}) do
    receive do
      {:cbmessage, %CBMessage{} = content} ->
        content.events
        |> Enum.map(fn %CBEvent{} = event ->
          productId = event.product_id

          productDistributorPIDs =
            Map.put_new(
              productDistributorPIDs,
              productId,
              spawn(ProductDistributor, :start, [productId])
            )

          send(productDistributorPIDs[productId], {:cbevent, event})

          productDistributorPIDs
        end)
        |> Enum.reduce(productDistributorPIDs, fn x, acc -> Map.merge(x, acc) end)
        |> start
    end
  end
end

defmodule JsonInterpreter do
  defp parse_json(coinbaseDistributorPID, json) do
    case Poison.decode(json,
           as: %CBMessage{
             events: [
               %CBEvent{
                 updates: [%CBUpdate{}]
               }
             ]
           }
         ) do
      {:ok, parsed} ->
        send(coinbaseDistributorPID, {:cbmessage, parsed})

      {:error, error} ->
        nil
    end
  end

  def start(pid \\ nil) do
    pid =
      if pid == nil do
        spawn(CurrencyDistributor, :start, [])
      else
        pid
      end

    receive do
      {:json, json, index} -> parse_json(pid, json)
    end

    start(pid)
  end
end

defmodule FileReader do
  def start() do
    # IO.puts("Starting JsonInterpreter")
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
