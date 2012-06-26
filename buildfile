# Generated by Buildr 1.4.6, change to your liking
require 'buildr/scala'

repositories.remote << 'http://repo1.maven.org/maven2/'


# Version number for this release
VERSION_NUMBER = "0.0.1"
# Group identifier for your projects
GROUP = "Apex"
COPYRIGHT = "Mark C. Chu-Carroll <markcc@gmail.com>"

# Specify Maven 2.0 remote repositories here, like this:
repositories.remote << "http://www.ibiblio.org/maven2/"

desc "The Apex project"
define "Apex" do
  project.version = VERSION_NUMBER
  project.group = GROUP
  manifest["Implementation-Vendor"] = COPYRIGHT
  compile.with "org.scala-tools.testing:specs_2.9.1:jar:1.6.9"
end
