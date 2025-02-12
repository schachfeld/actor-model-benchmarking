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
  def writeLine(parent, filename, file \\ nil) do
    file =
      if file == nil do
        File.open!(filename, [:append, :utf8])
      else
        file
      end

    receive do
      {:writeLine, content} ->
        IO.write(file, content)
        IO.write(file, "\n")

        writeLine(parent, filename, file)

      {:lastmessage} ->
        File.close(file)
        send(parent, {:donemessage})
    end
  end
end

defmodule AvgOrderBookCalculator do
  defp to_float(value) do
    cond do
      String.contains?(value, ".") -> String.to_float(value)
      true -> String.to_integer(value) * 1.0
    end
  end

  def calculateOrderBookAvg(parent, productId, filewriterPID \\ nil) do
    filewriterPID =
      if filewriterPID == nil do
        # IO.puts("Spawning filewriter")
        spawn(FileWriter, :writeLine, [self(), "orderbooks/" <> productId <> ".txt"])
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

      {:lastmessage} ->
        send(filewriterPID, {:lastmessage})

      {:donemessage} ->
        send(parent, {:donemessage})
    end

    calculateOrderBookAvg(parent, productId, filewriterPID)
  end
end

defmodule ProductDistributor do
  def start(parent, productId, avgOrderBookPID \\ nil) do
    avgOrderBookPID =
      if avgOrderBookPID == nil do
        spawn(AvgOrderBookCalculator, :calculateOrderBookAvg, [self(), productId])
      else
        avgOrderBookPID
      end

    receive do
      {:cbevent, %CBEvent{} = event} ->
        send(avgOrderBookPID, {:cbupdate, event.updates})

      {:lastmessage} ->
        send(avgOrderBookPID, {:lastmessage})

      {:donemessage} ->
        send(parent, {:donemessage, self()})
    end

    start(parent, productId, avgOrderBookPID)
  end
end

defmodule CurrencyDistributor do
  def start(parent, productDistributorPIDs \\ %{}, doneDistributors \\ []) do
    receive do
      {:cbmessage, %CBMessage{} = content} ->
        # content.events
        # |> Enum.map(fn %CBEvent{} = event ->
        #   productId = event.product_id
        #   # if productId == nil do
        #   #   IO.inspect(event)
        #   # end
        #   productDistributorPIDs =
        #     Map.put_new(
        #       productDistributorPIDs,
        #       productId,
        #       spawn(ProductDistributor, :start, [self(), productId])
        #     )
        #   send(productDistributorPIDs[productId], {:cbevent, event})
        #   productDistributorPIDs
        # end)
        # |> Enum.reduce(productDistributorPIDs, fn x, acc -> Map.merge(x, acc) end)

        productDistributorPIDs =
          content.events
          |> Enum.reduce(
            productDistributorPIDs,
            fn event, acc ->
              productDistributorPIDs =
                if Map.has_key?(acc, event.product_id) do
                  acc
                else
                  Map.put(
                    acc,
                    event.product_id,
                    spawn(ProductDistributor, :start, [self(), event.product_id])
                  )
                end

              send(productDistributorPIDs[event.product_id], {:cbevent, event})

              productDistributorPIDs
            end
          )

        start(parent, productDistributorPIDs, doneDistributors)

      {:lastmessage} ->
        Enum.map(
          productDistributorPIDs,
          fn {productId, pid} -> send(pid, {:lastmessage}) end
        )

        start(parent, productDistributorPIDs, doneDistributors)

      {:donemessage, pid} ->
        doneDistributors = [pid | doneDistributors]

        asdf = Map.keys(productDistributorPIDs)
        IO.puts("Done: #{length(doneDistributors)} #{length(asdf)}")

        if length(doneDistributors) == length(Map.keys(productDistributorPIDs)) do
          send(parent, {:donemessage})
        end

        start(parent, productDistributorPIDs, doneDistributors)
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
      {:ok, parsed} when parsed.channel == "l2_data" ->
        parsed.events
        |> Enum.map(fn %CBEvent{} = event ->
          productId = event.product_id

          if productId == nil do
            IO.inspect(json)
            IO.inspect(parsed)
          end

          send(coinbaseDistributorPID, {:cbevent, event})
        end)

        send(coinbaseDistributorPID, {:cbmessage, parsed})

      # skip this one
      {:ok, _parsed} ->
        nil

      {:error, error} ->
        nil
    end
  end

  def start(parent, pid \\ nil) do
    pid =
      if pid == nil do
        spawn(CurrencyDistributor, :start, [self()])
      else
        pid
      end

    receive do
      {:json, json, index} -> parse_json(pid, json)
      {:lastmessage} -> send(pid, {:lastmessage})
      {:donemessage} -> send(parent, {:donemessage})
    end

    start(parent, pid)
  end
end

defmodule FileReader do
  use Application

  def start(_type, _args) do
    # IO.puts("Starting JsonInterpreter")
    pid = spawn(JsonInterpreter, :start, [self()])

    starttime = System.monotonic_time(:millisecond)

    # File.stream!("../messages.log")
    File.stream!("../messages.log")
    |> Stream.with_index()
    |> Stream.map(fn {line, index} ->
      send(pid, {:json, line, index})
    end)
    # |> Stream.map(fn {line, index} -> IO.puts("index: #{index}") end)
    |> Stream.run()

    send(pid, {:lastmessage})

    receive do
      {:donemessage} ->
        endtime = System.monotonic_time(:millisecond)
        IO.puts("Time taken: #{endtime - starttime}ms")
    end

    :ok
  end
end
