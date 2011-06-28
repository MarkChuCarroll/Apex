// Copyright 2011 Mark C. Chu-Carroll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file contains the code that implements file management. Each
// file managed by the system is associated with a buffer. Operations
// to save and load files are handled by the file manager.

import java.io.File;
import scala.collection.mutable.HashMap

class FileManager {
  val _fileToBuffer = new HashMap[File, Buffer]()
  val _bufferToFile = new HashMap[Buffer, File]()

  def getBuffer(f : File) : GapBuffer = _bufferToFile.get(f)
    
  def getFile(b : GapBuffer) : File = _fileToBuffer.get(b)

  def saveFile(f : File) = {
    val buf = _bufferToFile.get(f)
    if (buf == null) {
      throw ...
    }
    val contents = 
  }
    
  def saveBuffer(b : GapBuffer) {
    if (file.exists()) {
      // if there's already a backup file, delete it, and then rename
      // the current version to the backup.
      val backup = new File(file.getPath() + ".BAK")
      if (backup.exists()) {
        backup.delete()
      }
      file.renameTo(backup)
    }
    // Convert the buffer to an array of lines to write to the file;
	 // we use a PrintWriter, because that will take care of
	 // OS-specific end-of-line canonicalization.
    val buffer_contents = new Array[Char](length())
    for (i <- 0 until length()) {
      buffer_contents(i) = char_at(i)
    }
    val bufferLines = new String(buffer_contents).split('\n')
    val out = new PrintWriter(new FileWriter(file))
    bufferLines foreach(line => out.println(line))
    out.close()    

  }

  /**
   * Open a file.
   * @param f the pathname of the file to open.
   * @param create a flag indicating whether or not the file should
   *   be created if no file with the name exists.
   * If the file doesn't exist, and create == false, then an IOException
   * will be raised.
   */
  def open(f : File, create : Bool) : GapBuffer

  /**
   * Close a buffer specified by file.
   * @param f the pathname of the file whose buffer should be closed.
   * @param force a flag indicating whether or not the file should be
   *   closed if it's been modified since the last save.
   * @return a flag indicating whether or not the file was closed. If
   * force == true, the buffer will always be closed, and so this
   * will always return true.
   */
  def closeFile(f : File, force : Bool) : Bool

  /**
   * Close a buffer
   * @param f the pathname of the file whose buffer should be closed.
   * @param force a flag indicating whether or not the file should be
   *   closed if it's been modified since the last save.
   * @return a flag indicating whether or not the file was closed. If
   * force == true, the buffer will always be closed, and so this
   * will always return true.
   */
  def closeBuffer(b : GapBuffer, force : Bool) : Bool
}
