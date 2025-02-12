package coinbaseanalyzer

import akka.actor.{Actor, Props}
import akka.actor.typed.{ActorRef, ActorSystem, Behavior}
import akka.actor.typed.scaladsl.Behaviors

import scala.io.Source
import spray.json.*
import DefaultJsonProtocol.*
import coinbaseanalyzer.CBJsonProtocol.cbMessageFormat


class JsonInterpreter extends Actor {
  val currencyDistributor = context.actorOf(Props(new CurrencyDistributor()), "currencyDistributor")


  def receive = {
    case FileLineMessage(text) =>

      try{
      val json = text.parseJson.convertTo[CBMessage]

        currencyDistributor ! json

        } catch {
          case e: Exception => null // discard the line if it is something else
        }

    case LastMessage() =>
      println("JsonInterpreter done")
      currencyDistributor ! LastMessage()
    case DoneMessage() =>
      context.parent ! DoneMessage()

  }
}
