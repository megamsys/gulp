/* 
** Copyright [2012] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
*/
/**
 * @author rajthilak
 *
 */

import org.specs2._
import scalaz._
import Scalaz._
import org.specs2.mutable._
import org.specs2.Specification
import org.megam.common._
import com.twitter.zk._
import com.twitter.util.{ Duration, Future, Promise, TimeoutException, Timer, Return, Await }
import org.apache.zookeeper.data.{ ACL, Stat }
import org.apache.zookeeper.KeeperException

class SBTSpecs extends Specification {

  def is =
    "Specs".title ^ end ^
      """
  SBT Monkey 
  """ ^ end ^
      "The SBT Monkey Should" ^
      //"Correctly invoke a build from a prescribed path" ! SBTBuild().createSucceeds ^     
      end

  trait TestContext {
    println("Setting up SBT Monkey)
    val path = "~/temp/playf"
   
  }

  case class SBTBuild() extends TestContext {

    def createSucceeds = //
     }

}





