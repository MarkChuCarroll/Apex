package apex

import scala.util.Try
import scala.util.Failure
import apex.TypeSafeEquality._
import scala.util.Success

class TextOperationException(msg: String) extends BufferException(msg)

case class Selection(start: Int, end: Int)

trait TextOperation {
  /** Executes the operation. This produces a new selection and a new buffer as a result;
    * the original buffer should be left unchanged. 
    */
  def execute(buffer: Buffer, sel: Selection): Try[(Buffer, Selection)]  
}

class SequenceOperation(val ops: Seq[TextOperation]) extends TextOperation {
  def execute_on_sequence(seq: Seq[TextOperation], buffer: Buffer, sel: Selection): Try[(Buffer, Selection)] = {
    if (seq.isEmpty) {
      Success((buffer, sel))
    } else {
      seq.head.execute(buffer, sel) match {
        case Success((new_buffer, new_selection)) => execute_on_sequence(seq.tail, new_buffer, new_selection)
        case Failure(e) => Failure(e)
      }
    }
  }

  override def execute(buffer: Buffer, sel: Selection) = {
    execute_on_sequence(ops, buffer, sel)
  }
}

//class ChoiceOperation(val ops: Seq[TextOperation]) extends TextOperation {
//  def execute(sel: Selection): Try[Selection] = {
//  
//  }
//}