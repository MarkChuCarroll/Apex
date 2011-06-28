package apex;

import java.util.Stack;

public class GapBuffer {
   private char[] _pre;
   private char[] _post;
   private int _prepos;
   private int _postpos;
   private int _line;
   private int _column;
   boolean _undoing = false;
   private Stack<UndoOperation> _undo = new Stack<UndoOperation>();

   public GapBuffer(int size) {
      _pre = new char[size];
      _post = new char[size];
      _prepos = 0;
      _postpos = 0;
      _line = 1;
      _column = 0;
   }

   public void clear() {
      _prepos = 0;
      _postpos = 0;
      _line = 1;
      _column = 0;
   }

   public void insert(String s) {
      int pos = currentPosition();
      for (int i = 0; i < s.length(); i++) {
         insertChar(s.charAt(i), false);
      }
      if (!_undoing) {
         UndoOperation undo = new UndoInsert(pos, s.length());
         _undo.push(undo);
      }
   }
   
   public void insertChars(char[] chars) {
      int pos = currentPosition();
      for (int i = 0; i < chars.length; i++) {
         insertChar(chars[i]);
      }
      if (!_undoing) {
         UndoOperation undo = new UndoInsert(pos, chars.length);
         _undo.push(undo);
      }
   }
   
   public void insertChar(char c) {
      insertChar(c, true);
   }

   public void insertChar(char c, boolean recordUndo) {
      if (recordUndo && !_undoing) {
         UndoOperation undo = new UndoInsert(currentPosition(), 1);
         _undo.push(undo);
      }
      pushPre(c);
      forwardUpdatePosition(c);
   }

   /**
    * Move the cursor to a position.
    * 
    * @param pos the position to move to.
    */
   public void moveTo(int pos) {
      int distance = pos - currentPosition();
      moveBy(distance);
   }

   public int moveToLine(int line) {
      moveTo(0);
      while (_postpos > 0 && _line < line) {
         stepForward();
      }
      return _prepos;
   }

   public void moveBy(int dist) {
      if (dist > 0) {
         if (dist > _postpos) {
            dist = _postpos;
         }
         for (int i = 0; i < dist; i++) {
            stepForward();
         }
      } else {
         dist = -dist;
         if (dist > _prepos) {
            dist = _prepos;
         }
         for (int i = 0; i < dist; i++) {
            stepBackward();
         }
      }
   }

   public void moveByLine(int numLines) {
      int targetLine = _line + numLines;
      moveToLine(targetLine);
   }

   /**
    * Cut a region of text starting at the cursor point.
    * @param len the length of the region to cut. This must be
    *   greater than 0.
    * @return the cut text.
    */
   public char[] cut(int len) {
      if (len < 0) {
         throw new IllegalArgumentException("cut length must not be negative");
      }
      if (len >_postpos) {
         len = _postpos;
      }
      int pos = currentPosition();
      char[] result = new char[len];
      for (int i = 0; i < len; i++) {
         result[i] = popPost();
      }
      if (!_undoing) {
         UndoOperation undo = new UndoDelete(pos, result);
         _undo.push(undo);
      }
      return result;
   }

   public char[] copy(int len) {
      if (len < 0) {
         throw new IllegalArgumentException("copy length must not be negative");
      }
      if (len > _postpos) {
         len = _postpos;
      }
      char[] result = new char[len];
      for (int i = 0; i < len; i++) {
         result[i] = charAt(i + currentPosition());
      }
      return result;      
   }
   
   public char charAt(int pos) {
      if (pos < 0) {
         throw new IllegalArgumentException("character index must be >= 0");
      }
      if (pos > length()) {
         throw new IllegalArgumentException("character index past buffer end");
      }
      if (pos < _prepos) {
         return _pre[pos];
      } else {
         return _post[_postpos - (pos - _prepos) - 1];
      }
   }
   
   public int length() {
      return _prepos + _postpos;
   }

   public int currentPosition() {
      return _prepos;
   }
   
   public int currentLine() {
      return _line;
   }
   
   public int currentColumn() {
      return _column;
   }

   public void insertAt(int pos, String s) {
      moveTo(pos);
      insert(s);
   }
   
   public void insertCharsAt(int pos, char[] chars) {
      moveTo(pos);
      for (int i = 0; i < chars.length; i++) {
         insertChar(chars[i]);
      }
      if (!_undoing) {
         UndoOperation undo = new UndoInsert(pos, chars.length);
         _undo.push(undo);
      }
   }
   
   public char[] cutRange(int pos, int len) {
      moveTo(pos);
      return cut(len);
   }

   public char[] copyRange(int pos, int len) {
      moveTo(pos);
      return copy(len);
   }

   public int[] getLineAndColumnOf(int pos) {
      int currentPos = currentPosition();
      moveTo(pos);
      int[] result = new int[] { currentLine(), currentColumn() };
      moveTo(currentPos);
      return result;
   }

   /**
    * Get the character index of the first character of the Nth
    * line. 
    * @param target
    * @return the index of the line, -1 if there aren't enough lines
    *    in the file.
    */
   public int getPositionOfLine(int target) {
      if (target == 1) {
         return 0;
      }
      int bufferLength = length();
      int line = 1;
      int idx = 0;
      while(idx < bufferLength && line < target) {
         if (charAt(idx) == '\n') {
            line++;
            if (line == target) {
               return idx + 1;
            }
         }
         idx++;
      }
      return -1;
   }

   /**
    * Run a single undo step.
    * @return true if the undo completed successfully;
    *   false if there were no further undos on the stack.
    */
   public boolean undo() {
      if (!_undo.isEmpty()) {
         _undo.pop().execute();
         return true;
      } else {
         return false;
      }
   }

   void pushPre(char c) {
      if (_prepos == _pre.length) {
         expandCapacity();
      }
      _pre[_prepos] = c;
      _prepos++;
   }

   private void expandCapacity() {
      char[] newPre = new char[2*_pre.length];
      for (int i = 0; i < _prepos; i++) {
         newPre[i] = _pre[i];
      }
      _pre = newPre;
      char[] newPost = new char[2*_post.length];
      for (int i = 0; i < _postpos; i++) {
         newPost[i] = _post[i];
      }
      _post = newPost;
   }

   void pushPost(char c) {
      if (_postpos == _post.length) {
         expandCapacity();
      }
      _post[_postpos] = c;
      _postpos++;

   }

   char popPre() {
      _prepos--;
      return _pre[_prepos];
   }

   char popPost() {
      _postpos--;
      return _post[_postpos];
   }

   void stepForward() {
      char c = popPost();
      pushPre(c);
      forwardUpdatePosition(c);
   }

   void stepBackward() {
      char c = popPre();
      pushPost(c);
      reverseUpdatePosition(c);
   }

   void reverseUpdatePosition(char c) {
      if (c == '\n') { // we've stepped backwards over a newline, so
         // we need to update.
         int i = 0;
         while (i < _prepos && _pre[_prepos - i - 1 ] != '\n') {
            i++;
         }
         _column = i;
         _line--;
      } else {
         _column--;
      }
   }

   /**
    * Update the position after adding a character to the region before the
    * cursor. This can calculate the cursor position after either a character
    * insert, or stepping the cursor forward.
    * 
    * @param c
    */
   void forwardUpdatePosition(char c) {      
      if (c == '\n') {
         _line++;
         _column = 0;
      } else {
         _column++;
      }
   }
   
   public String debugString() {
      StringBuilder result = new StringBuilder();
      result.append("{");
      for (int i = 0; i < _prepos; i++) {
         result.append(_pre[i]);
      }
      result.append("}GAP{");
      for (int i = 0; i < _postpos; i++) {
         result.append(_post[_postpos - i - 1]);
      }
      result.append("}");
      return result.toString();
   }
   
   interface UndoOperation { 
      void execute();
   }
   
   private class UndoInsert implements UndoOperation {
      private int _pos;
      private int _len;

      public UndoInsert(int pos, int len) {
         this._pos = pos;
         this._len = len;
      }

      @Override
      public void execute() {
         GapBuffer.this._undoing = true;
         GapBuffer.this.cutRange(_pos, _len);
         GapBuffer.this._undoing = false;
      }
   }
      
   private class UndoDelete implements UndoOperation {
      private int _pos;
      private char[] _deleted;

      public UndoDelete(int pos, char[] deleted) {
         this._deleted = deleted;
         this._pos = pos;
      }
         
      @Override
      public void execute() {
         GapBuffer.this._undoing = true;
         GapBuffer.this.insertCharsAt(_pos, _deleted);
         GapBuffer.this._undoing = false;
      }
      
   }
}
   
   
