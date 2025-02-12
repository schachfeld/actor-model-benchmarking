package coinbaseanalyzer

import akka.actor.typed.scaladsl.Behaviors
import akka.actor.{Actor, Props}


class MessagePasser(productId: String) extends Actor {
  private val avgOrderBookTracker = context.actorOf(Props(new AvgOrderBookTracker(productId)), "avgOrderBookTracker")


  def receive = {
    case CBEvent(
    event_type,
    product_id,
    updates,
    ) =>
      avgOrderBookTracker ! updates

    case LastMessage() =>
      println(s"MessagePasser $productId done")
      avgOrderBookTracker ! LastMessage()

    case DoneMessage() =>
      context.parent ! DoneMessageProduct(productId)
  }
}
