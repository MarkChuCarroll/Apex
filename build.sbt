name := "Apex Editor"

version := "0.1"

scalaVersion := "2.9.1"

libraryDependencies ++= Seq(
    "org.specs2" %% "specs2" % "1.6.1",
    "org.specs2" %% "specs2-scalaz-core" % "6.0.1" % "test",
    "org.specs2" %% "specs2" % "1.5",
    "org.specs2" %% "specs2-scalaz-core" % "5.1-SNAPSHOT" % "test" )

  resolvers ++= Seq("snapshots" at "http://scala-tools.org/repo-snapshots",
                    "releases"  at "http://scala-tools.org/repo-releases")
