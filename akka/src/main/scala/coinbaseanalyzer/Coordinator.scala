package coinbaseanalyzer

import akka.actor.{Actor, ActorSystem, Props}


object Coordinator extends App {
  val system = ActorSystem("CoinbaseAnalyzer")

  val startTime = System.currentTimeMillis()

  val listener = system.actorOf(Props(new Actor {
    def receive = {
      case DoneMessage() =>
        val endTime = System.currentTimeMillis()
        println(s"Time taken: ${endTime - startTime} ms")
        system.terminate()
    }
  }))

  val fileReader = system.actorOf(Props(new FileReader(listener)), "fileReader")


  fileReader ! StartMessage("../coinbase-analyzer/messages_short.log")
}