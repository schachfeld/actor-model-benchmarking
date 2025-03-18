import sbt.Keys.libraryDependencies
import sbt.Keys.resolvers

import scala.collection.Seq

ThisBuild / version := "0.1.0-SNAPSHOT"

ThisBuild / scalaVersion := "3.3.4"

val AkkaVersion = "2.10.0"

libraryDependencies += "io.spray" %% "spray-json" % "1.3.6"

Compile / mainClass := Some("example.Pinger")

lazy val root = (project in file("."))
  .settings(
    name := "akka",
    resolvers += "Akka library repository".at("https://repo.akka.io/maven"),
    libraryDependencies ++= Seq(
      "ch.qos.logback" % "logback-classic" % "1.5.8",
      "com.typesafe.akka" %% "akka-actor-typed" % AkkaVersion,
      "com.typesafe.akka" %% "akka-actor-testkit-typed" % AkkaVersion % Test
    )
  )

enablePlugins(AssemblyPlugin)

mainClass in assembly := Some("example.LatencyBenchmark")

import sbtassembly.AssemblyPlugin.autoImport._

assemblyMergeStrategy in assembly := {
  case PathList("reference.conf")    => MergeStrategy.concat
  case PathList("application.conf")  => MergeStrategy.concat
  case PathList("META-INF", xs @ _*) => MergeStrategy.discard
  case _                             => MergeStrategy.first
}
