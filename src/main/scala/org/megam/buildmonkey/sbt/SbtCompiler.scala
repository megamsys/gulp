/* 
** Copyright [2012-2013] [Megam Systems]
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
package org.megam.buildmonkey.sbt

import scala.tools.nsc.Global
import scala.tools.nsc.Settings
import scala.tools.nsc.io.AbstractFile
import scala.reflect.internal.util.NoPosition
import scala.collection.mutable
import java.io.File
import xsbti.compile.CompileProgress
import xsbti.Logger
import xsbti.F0
import sbt.Process
import sbt.ClasspathOptions

import sbt.inc.AnalysisStore
import sbt.inc.Analysis
import sbt.inc.FileBasedStore
import sbt.inc.Incremental
import sbt.compiler.IC
import sbt.compiler.CompileFailed
import java.lang.ref.SoftReference
import java.util.concurrent.atomic.AtomicReference
import org.slf4j.LoggerFactory

/**
 * @author rajthilak
 *
 */
class SbtCompiler {

  private lazy val logger = LoggerFactory.getLogger(getClass)

  
  private val sbtLogger = new xsbti.Logger {
    override def error(msg: F0[String]) = logger.error(msg())
    override def warn(msg: F0[String]) = logger.warn(msg())
    override def info(msg: F0[String]) = logger.info(msg())
    override def debug(msg: F0[String]) = logger.debug(msg())
    override def trace(exc: F0[Throwable]) = logger.error("", exc())
  }
  
  
}