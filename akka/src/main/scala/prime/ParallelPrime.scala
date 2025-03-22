package prime

import akka.actor.Actor
import akka.actor.ActorRef
import akka.actor.ActorSystem
import akka.actor.Props

import java.io.*
import scala.annotation.tailrec
import scala.collection.mutable.ListBuffer
import scala.math.sqrt

object PrimeCalculator {
  def isPrimeI(n: Int): Boolean = {
    if (n < 2) false
    else if (n == 2) true
    else {
      var isPrime = true
      val limit = math.sqrt(n).toInt
      for (i <- 2 to limit) {
        if (n % i == 0) {
          isPrime = false
        }
      }
      isPrime
    }

  }

  def isPrime(n: Int): Boolean = {
    @tailrec
    def checkDivisibility(i: Int): Boolean = {
      if (i > math.sqrt(n).toInt) true
      else if (n % i == 0) false
      else checkDivisibility(i + 1)
    }

    if (n < 2) false
    else if (n == 2) true
    else checkDivisibility(2)
  }
}

case class CalculatePrimes(range: Range)

case class PrimeResult(primes: List[Int])

case object StartCalculation

class PrimeWorker extends Actor {
  def receive = { case CalculatePrimes(range) =>
    val primes = range.filter(PrimeCalculator.isPrime).toList
    sender() ! PrimeResult(primes)
  }
}

class PrimeCoordinator(totalWorkers: Int, range: Range) extends Actor {
  var collectedPrimes: ListBuffer[Int] = ListBuffer()
  var resultsReceived = 0
  val startTime = System.nanoTime()

  val chunkSize = range.length / totalWorkers
  val workers: List[ActorRef] = (1 to totalWorkers).map { _ =>
    context.actorOf(Props[PrimeWorker]())
  }.toList

  def receive = {
    case StartCalculation =>
      for ((worker, i) <- workers.zipWithIndex) {
        val start = range.start + i * chunkSize
        val end =
          if (i == totalWorkers - 1) range.end else start + chunkSize - 1
        worker ! CalculatePrimes(start to end)
      }

    case PrimeResult(primes) =>
      collectedPrimes ++= primes
      resultsReceived += 1

      if (resultsReceived == totalWorkers) {
        val endTime = System.nanoTime()
        println(s"Found ${collectedPrimes.length} prime numbers.")
        println(s"Calculation took ${endTime - startTime} nanoseconds.")

        val fileName =
          s"prime_bench_results/run_2/10mil_${totalWorkers}workers.txt"
        val writer = new PrintWriter(new FileWriter(fileName, true))
        writer.append(s"${endTime - startTime},")
        writer.close()
        context.system.terminate()
      }
  }
}

object ParallelPrimeApp {
  def main(args: Array[String]): Unit = {
    if (args.length < 1) {
      println("Usage: ParallelPrimeApp <totalWorkers>")
      System.exit(1)
    }

    val totalWorkers = args(0).toInt
    val system = ActorSystem("PrimeSystem")
    val range = 1 to 10_000_000

    val coordinator = system.actorOf(
      Props(new PrimeCoordinator(totalWorkers, range)),
      "coordinator"
    )
    coordinator ! StartCalculation
  }
}

object HighestPrime extends App {
  val startTime = System.nanoTime()
  val isPrime = PrimeCalculator.isPrime(9_999_991)
//  val isPrime = PrimeCalculator.isPrime(7)

  val endTime = System.nanoTime()
  println(s"${isPrime}")
  println(s"Calculation took ${endTime - startTime} nanoseconds.")
}

object RangeWithOne extends App {
  val system = ActorSystem("PrimeSystem")
  val worker = system.actorOf(Props(new PrimeWorker()), "worker")

  val startTime = System.currentTimeMillis()

  val listener = system.actorOf(Props(new Actor {
    def receive = { case PrimeResult(primes) =>
      val endTime = System.currentTimeMillis()
      println(s"Found ${primes.length} prime numbers.")
      println(s"Calculation took ${endTime - startTime} milliseconds.")
      context.system.terminate()
    }
  }))

  worker.tell(CalculatePrimes(1 to 1_000_000), listener)
}
