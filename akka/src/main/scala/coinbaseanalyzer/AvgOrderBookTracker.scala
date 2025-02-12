package coinbaseanalyzer

import akka.actor.{Actor, Props}


class AvgOrderBookTracker(productId: String) extends Actor {
  val fileWriter = context.actorOf(Props(new FileWriter(productId)), "fileWriter")


  def receive = {
    case updates: Seq[CBUpdate] =>

      val bids = updates
        .filter(_.side == "bid")
        .map(_.price_level.toFloat)
      val avgBids = bids.sum / bids.length

      val offers = updates
        .filter(_.side == "offer")
        .map(_.price_level.toFloat)
      val avgOffers = offers.sum / offers.length

      fileWriter ! WriteLineMessage(s"Bids: $avgBids")
      fileWriter ! WriteLineMessage(s"Offers: $avgOffers")

    case LastMessage() =>
      println(s"AvgOrderBookTracker $productId done")
      fileWriter ! LastMessage()
    case DoneMessage() =>
      context.parent ! DoneMessage()
  }
}
