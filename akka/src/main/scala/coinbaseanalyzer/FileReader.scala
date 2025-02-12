package coinbaseanalyzer

import akka.actor.{Actor, ActorRef, Props}

import scala.io.Source

class FileReader(listener: ActorRef) extends Actor {
  val jsonInterpreter = context.actorOf(Props(new JsonInterpreter()), "jsonInterpreter")

  def receive = {
    case StartMessage(filename) =>
      val source = Source.fromFile(
        filename
      )
      for (line <- source.getLines())
        jsonInterpreter ! FileLineMessage(line)

      jsonInterpreter ! LastMessage()
      println("FileReader done")
      source.close()

    case DoneMessage() =>
      listener ! DoneMessage()
  }
}
