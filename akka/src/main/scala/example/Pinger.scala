package example

import akka.actor.typed.{ActorRef, ActorSystem, Behavior}
import akka.actor.typed.scaladsl.Behaviors


//#hello-world-actor
object Pong {
  final case class Greet(whom: String)

  def apply(): Behavior[Greet] = Behaviors.receive { (context, message) =>
    Behaviors.same
  }
}
//#hello-world-actor

//#hello-world-main
object Pinger {

  final case class Ping(message: String)

  def apply(): Behavior[Ping] =
    Behaviors.setup { context =>
      val pong = context.spawn(Pong(), "pong")

      Behaviors.receiveMessage { message =>
        pong ! Pong.Greet(message.message)
        Behaviors.same
      }
    }

  //#hello-world-main
  def main(args: Array[String]): Unit = {

    val system: ActorSystem[Pinger.Ping] =
      ActorSystem(Pinger(), "hello")

    val messageNum = 100_000_000

    for (i <- 0 until messageNum) {
      system ! Pinger.Ping("World")
    }

    println(s"Sent $messageNum messages")


    Thread.sleep(3000)
    system.terminate()
  }
  //#hello-world-main
}
//#hello-world-main