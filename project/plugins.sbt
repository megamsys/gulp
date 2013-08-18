// Comment to get more information during initialization
logLevel := Level.Info

// Typesafe snapshots

resolvers += "Sonatype Snapshots"  at "https://oss.sonatype.org/content/repositories/snapshots"

resolvers += "Typesafe Snapshots" at "http://repo.typesafe.com/typesafe/snapshots/"

resolvers += "Scala-Tools Maven2 Snapshots Repository" at "http://scala-tools.org/repo-snapshots"

// The Typesafe repository 
resolvers += "Typesafe repository" at "http://repo.typesafe.com/typesafe/releases/"

resolvers += "Sonatype Releases"  at "https://oss.sonatype.org/content/repositories/releases"

 resolvers += "JBoss repository" at "https://repository.jboss.org/nexus/content/repositories/"    
 
 resolvers += Resolver.url("scalasbt", new URL("http://scalasbt.artifactoryonline.com/scalasbt/sbt-plugin-releases"))(Resolver.ivyStylePatterns)
 
addSbtPlugin("net.virtual-void" % "sbt-dependency-graph" % "0.7.0")

addSbtPlugin("com.github.gseitz" % "sbt-release" % "0.7")

addSbtPlugin("org.scalastyle" %% "scalastyle-sbt-plugin" % "0.2.0")
 