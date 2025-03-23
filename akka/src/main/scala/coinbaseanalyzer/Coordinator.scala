package coinbaseanalyzer

import akka.actor.Actor
import akka.actor.ActorSystem
import akka.actor.Props

import java.nio.file.Files
import java.nio.file.Paths
import java.nio.file.StandardOpenOption

object Coordinator {
  def main(args: Array[String]): Unit = {
    val system = ActorSystem("CoinbaseAnalyzer")

    val startTime = System.nanoTime()

    val listener = system.actorOf(Props(new Actor {
      def receive = { case DoneMessage() =>
        val endTime = System.nanoTime()
        println(s"Time taken: ${endTime - startTime} ns")

        val elapsed = endTime - startTime
        val result = s"$elapsed,"

        Files.write(
          Paths.get("cb_analyzer_results/cb_analyzer.txt"),
          result.getBytes,
          StandardOpenOption.CREATE,
          StandardOpenOption.APPEND
        )

        system.terminate()
      }
    }))

    val fileReader =
      system.actorOf(Props(new FileReader(listener)), "fileReader")

    fileReader ! StartMessage("../coinbase-analyzer/messages.log")
  }
}
