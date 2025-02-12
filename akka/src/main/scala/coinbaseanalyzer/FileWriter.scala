package coinbaseanalyzer

import akka.actor.typed.{ActorRef, ActorSystem, Behavior}
import akka.actor.{Actor, Props}

import java.io.{File, PrintWriter}


class FileWriter(productId: String) extends Actor {

  val file = new File(s"avgorderbook/$productId.txt")
  val writer = new PrintWriter(file)

  def receive = {
    case WriteLineMessage(
    text
    ) =>
      writer.println(text)

    case LastMessage() =>
      println(s"FileWriter $productId done")
      writer.close()
      context.parent ! DoneMessage()

  }


}
