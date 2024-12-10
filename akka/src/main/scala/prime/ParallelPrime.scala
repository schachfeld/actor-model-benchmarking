package prime

import akka.actor.{Actor, ActorRef, ActorSystem, Props}


import scala.math.sqrt
import scala.collection.mutable.ListBuffer

object PrimeCalculator {
  def isPrime(n: Int): Boolean = {
    if (n < 2) false
    else if (n == 2) true
    else !(2 to sqrt(n).toInt).exists(x => n % x == 0)

    // try implementing as a loop
  }
  
  def isPrimeRecursive(n: Int): Boolean = {
    
    if (n < 2) false
    else if (n == 2) true
    else 
  }
}

case class CalculatePrimes(range: Range)
case class PrimeResult(primes: List[Int])
case object StartCalculation

class PrimeWorker extends Actor {
  def receive = {
    case CalculatePrimes(range) =>
      val primes = range.filter(PrimeCalculator.isPrime).toList
      sender() ! PrimeResult(primes)
  }
}


class PrimeCoordinator(totalWorkers: Int, range: Range) extends Actor {
  var collectedPrimes: ListBuffer[Int] = ListBuffer()
  var resultsReceived = 0
  val startTime = System.currentTimeMillis()

  val chunkSize = range.length / totalWorkers
  val workers: List[ActorRef] = (1 to totalWorkers).map { _ =>
    context.actorOf(Props[PrimeWorker]())
  }.toList

  def receive = {
    case StartCalculation =>
      for ((worker, i) <- workers.zipWithIndex) {
        val start = range.start + i * chunkSize
        val end = if (i == totalWorkers - 1) range.end else start + chunkSize - 1
        worker ! CalculatePrimes(start to end)
      }


    case PrimeResult(primes) =>
      collectedPrimes ++= primes
      resultsReceived += 1

      if (resultsReceived == totalWorkers) {
        val endTime = System.currentTimeMillis()
        println(s"Found ${collectedPrimes.length} prime numbers.")
        println(s"Calculation took ${endTime - startTime} milliseconds.")
        context.system.terminate()
      }
  }
}
IO
// Main object to run the application
object ParallelPrimeApp extends App {
  val system = ActorSystem("PrimeSystem")
  val range = 1 to 1_000_000
  val totalWorkers = 1

  val coordinator = system.actorOf(Props(new PrimeCoordinator(totalWorkers, range)), "coordinator")
  coordinator ! StartCalculation
}

object HighestPrime {
  def main(args: Array[String]) = {
    val startTime = System.currentTimeMillis()
    val isPrime = PrimeCalculator.isPrime(9_999_991)
    val endTime = System.currentTimeMillis()
    println(s"${isPrime}")
    println(s"Calculation took ${endTime - startTime} milliseconds.")
  }
}