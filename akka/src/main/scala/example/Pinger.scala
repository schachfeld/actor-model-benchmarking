package example

import akka.actor.typed.ActorRef
import akka.actor.typed.ActorSystem
import akka.actor.typed.Behavior
import akka.actor.typed.scaladsl.Behaviors

import java.nio.charset.StandardCharsets
import java.nio.file.Files
import java.nio.file.Paths
import java.nio.file.StandardOpenOption

object Pinger {

  final case class Ping(message: Int)

  def apply(): Behavior[Ping] =
    Behaviors.setup { context =>
      Behaviors.receiveMessage { message =>
        Behaviors.same
      }
    }

  def main(args: Array[String]): Unit = {

    println("Start")

    val system: ActorSystem[Pinger.Ping] =
      ActorSystem(Pinger(), "hello")

    val messageNum = 10_000_000

    val startTime = System.nanoTime()
    for (i <- 0 until messageNum) {
      system ! Pinger.Ping(i)
    }
    val endTime = System.nanoTime()
    val elapsedTime = endTime - startTime
    val elapsedTimeSec = elapsedTime / 1_000_000_000

    println(s"Sent $messageNum messages")
    println(s"Time taken: $elapsedTimeSec s")
    println(s"Throughput: ${messageNum / elapsedTimeSec} msg/s")

    val n = 1 // You can change this to any number you want
    val filePath = s"throughput_bench_results/throughput_$n.txt"
    val content = s"$elapsedTime,"

    Files.write(
      Paths.get(filePath),
      content.getBytes(StandardCharsets.UTF_8),
      StandardOpenOption.CREATE,
      StandardOpenOption.APPEND
    )

    // Thread.sleep(3000)
    system.terminate()
  }
}
