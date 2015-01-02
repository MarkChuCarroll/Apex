package apex

import scala.util.Failure
import scala.util.Success
import scala.util.Try

/** A utility class which represents one of the two contiguous text segments of
  * a gap buffer. The main behavior of this resembles a stack, where you push or
  * pop off of the active end. 
  */
trait BufferSegment {
  // Convenience for generating better error messages
  def name: String
  
  def length: Int
  
  def push(c: Char): Unit
  
  def pop: Try[Char]

  def at(idx: Int): Try[Char]

  def all: String
  
  def isEmpty: Boolean

  def reset: Unit
}


class ArrayBufferSegment(val name: String) extends BufferSegment {
  val initLength = 65536

  /** Returns the current number of characters in this buffer.
    */
  var length: Int = 0
  
  private var chars: Array[Char] = new Array[Char](initLength)
  
  /** Pushes a character onto the end of the buffer.
    */
  def push(c: Char) {
    if (length >= chars.length - 1) {
      val newChars = new Array[Char](chars.length * 2)
      for (i <- 0 until length) {
        newChars(i) = chars(i)
        chars = newChars
      }
    }
    chars(length) = c
    length = length + 1
  }
  
  /** Pops character off the active end of the buffer. Throws an
    * exception if the buffer is empty. 
    */
  def pop: Try[Char] = {
    if (length > 0) {
      length = length - 1
      Success(chars(length))
    } else {
      Failure(new BufferStackException(s"Text segment $name empty"))
    } 
  }
  
  /** Get the character at a position.
    * @param idx the character position 
    * @return Some of the character, or else None.
    */ 
  def at(idx: Int): Try[Char] = {
    if (idx >= length) {
      Failure(new BufferStackException("Character not found at index"))
    } else {
      Success(chars(idx))
    }
  }
  
  /** Get the entire buffer 
    */
  def all: String = new String(chars.slice(0, length))
  
  def isEmpty: Boolean = length == 0

  /** Reset the buffer, deleting all of its contents.
    */
  def reset {
    length = 0
  }
}


class ListBufferSegment(val name: String) extends BufferSegment {
  
  private var contents: List[Char] = Nil

  def length: Int = contents.length
  
  /** Pushes a character onto the end of the buffer.
    */
  def push(c: Char) {
    contents = (c :: contents)
  }
  
  /** Pops character off the active end of the buffer. Throws an
    * exception if the buffer is empty. 
    */
  def pop: Try[Char] = {
    contents match {
      case (c::rest) => {
        contents = rest
        Success(c) 
      }
      case Nil => {
        Failure(new BufferStackException(s"Text segment $name empty"))        
      }
    }
  }
  
  /** Get the character at a position.
    * @param idx the character position 
    * @return Some of the character, or else None.
    */ 
  def at(idx: Int): Try[Char] = {    
    if (idx >= length) {
      Failure(new BufferStackException("Character not found at index"))
    } else {
      Success(contents(length - idx - 1))
    }
  }
  
  /** Get the entire buffer 
    */
  def all: String = new String(contents.reverse.toArray)
  
  def isEmpty = contents.isEmpty

  /** Reset the buffer, deleting all of its contents.
    */
  def reset {
    contents = Nil
  }
}