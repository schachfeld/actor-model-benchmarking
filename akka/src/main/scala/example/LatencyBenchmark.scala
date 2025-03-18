package example

import akka.actor.Actor
import akka.actor.ActorSystem
import akka.actor.Props

object LatencyBenchmark extends App {
  class BenchActor extends Actor {
    def receive: Receive = { case "ok" => // do nothing
    }
  }

  def benchLatency(n: Int, system: ActorSystem): Unit = {
    for (_ <- 1 to n) {
      val startTime = System.nanoTime()
      val actor = system.actorOf(Props[BenchActor]())
      // actor ! "ok"
      val endTime = System.nanoTime()
      val elapsed = endTime - startTime
      print(s"$elapsed,")
    }
    system.terminate()
  }

  val n = 1_000_000
  val system = ActorSystem("LatencyBenchmarkSystem")

  benchLatency(n, system)
}
