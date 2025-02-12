package coinbaseanalyzer

import akka.actor.{Actor, ActorRef, Props}

import scala.collection.mutable

class CurrencyDistributor extends Actor {

  val messagePassers = mutable.Map[String, ActorRef]()

  def receive = {
    case CBMessage(
    channel,
    client_id,
    timestamp,
    sequence_num,
    events,
    ) =>
      for (event <- events) {
        val containsasdf = !messagePassers.contains(event.product_id)

        val messagePasser = messagePassers.getOrElseUpdate(event.product_id, context.actorOf(Props(new MessagePasser(event.product_id)), "messagePasser-" + event.product_id))


        messagePasser ! event
      }

    case LastMessage() =>
      println("CurrencyDistributor done")
      for ((_, messagePasser) <- messagePassers) {
        messagePasser ! LastMessage()
      }

    case DoneMessageProduct(productId) =>
      messagePassers.remove(productId)
      if (messagePassers.isEmpty) {
        context.parent ! DoneMessage()
      }
  }
}
