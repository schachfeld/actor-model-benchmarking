import sbt.Keys.{libraryDependencies, resolvers}

import scala.collection.Seq

ThisBuild / version := "0.1.0-SNAPSHOT"

ThisBuild / scalaVersion := "3.3.4"

val AkkaVersion = "2.10.0"


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


