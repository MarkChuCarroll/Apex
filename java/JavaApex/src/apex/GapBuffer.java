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
package apex;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.Stack;

/**
 * An editor buffer implementation based on the gap buffer model.
 * 
 * @author markcc
  */
public class GapBuffer {
   private char[] _pre;
   private char[] _post;
   private int _prepos;
   private int _postpos;
   private int _line;
   private int _column;
   boolean _undoing = false;
   private Stack<UndoOperation> _undo = new Stack<UndoOperation>();
   private File _path;
   
   static final int DEFAULT_SIZE = 65536;

   /**
    * Create a new gapbuffer.
    * @param size the initial size of the buffer.
    */
   public GapBuffer(File path, boolean create) throws IOException {
      _path = path;
      int initialSize = DEFAULT_SIZE;
      if (_path.exists()) {
         initialSize = (int)path.length() * 2;
      }
      if (!_path.exists() && !create) {
         throw new FileNotFoundException(path.toString());
      }
      _pre = new char[initialSize];
      _post = new char[initialSize];
      _prepos = 0;
      _postpos = 0;
      _line = 1;
      _column = 0;
      read();
   }
   
   public GapBuffer() {
      _path = null;
      _pre = new char[DEFAULT_SIZE];
      _post = new char[DEFAULT_SIZE];
      _prepos = 0;
      _postpos = 0;
      _line = 1;
      _column = 0;
   }

   /**
    * Clear the contents of a buffer.
    */
   public void clear() {
      _prepos = 0;
      _postpos = 0;
      _line = 1;
      _column = 0;
   }

   /**
    * Insert a string into the buffer at the current cursor position.
    * @param s
    */
   public void insert(String s) {
      int pos = currentPosition();
      for (int i = 0; i < s.length(); i++) {
         internalInsertChar(s.charAt(i), false);
      }
      if (!_undoing) {
         UndoOperation undo = new UndoInsert(pos, s.length());
         _undo.push(undo);
      }
   }
   
   /**
    * Insert an array of characters into the buffer at the current cursor
    * position.
    * @param chars
    */
   public void insertChars(char[] chars) {
      int pos = currentPosition();
      for (int i = 0; i < chars.length; i++) {
         internalInsertChar(chars[i], false);
      }
      if (!_undoing) {
         UndoOperation undo = new UndoInsert(pos, chars.length);
         _undo.push(undo);
      }
   }
   
   /**
    * Insert a single character at the current cursor position.
    * @param c
    */
   public void insertChar(char c) {
      internalInsertChar(c, true);
   }

   /**
    * Insert a character at the current cursor position. This is an internal method,
    * not for use by clients. It provides the capability to prevent generating an
    * undo record for the insert 
    * @param c
    * @param recordUndo if  false, don't record the insertion for undo.
    */
   private void internalInsertChar(char c, boolean recordUndo) {
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

   /**
    * Move the cursor to the first character of a line.
    * @param line
    */
   public void moveToLine(int line) {
      moveTo(0);
      while (_postpos > 0 && _line < line) {
         stepForward();
      }
   }

   /**
    * Move the cursor by a distance.
    * @param dist
    */
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

   /**
    * Move the cursor by a distance specified in lines.
    * @param numLines
    */
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

   /**
    * Copy a region of text starting at the cursor point.
    * @param len the length of the region to copy. This must be
    *   greater than 0.
    * @return the copied text.
    */
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
   
   /**
    * Retrieve the character at a position.
    * @param pos
    * @return
    */
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
   
   /**
    * Get the current length of the text in the buffer.
    * @return
    */
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
   
   public String allText() {
      StringBuilder result = new StringBuilder();
      for (int i = 0; i < _prepos; i++) {
         result.append(_pre[i]);
      }
      for (int i = 0; i < _postpos; i++) {
         result.append(_post[_postpos - i - 1]);
      }
      return result.toString();
   }
   
   public File getPath() {
      return _path;
   }
   
   /**
    * Read the contents of the file into the buffer.
    * @throws IOException
    */
   public void read() throws IOException {
      clear();
      BufferedReader in = new BufferedReader(new FileReader(_path));
      for (String s = in.readLine(); s != null; s = in.readLine()) {
         insert(s + '\n');
      }
      in.close();
   }
   
   public void write() throws IOException {
      // For writing to the buffer file, we always do backups.
      writeTo(_path, true);
   }
   
   public void writeTo(File f, boolean backup) throws IOException {
      if (backup) {
         if (f.exists()) {
            File backupFile = new File(f.getAbsolutePath() + ".BAK");
            if (backupFile.exists()) {
               backupFile.delete();
            }
            _path.renameTo(backupFile);
         }
      }
      String[] lines = allText().split("\n");
      PrintWriter out = new PrintWriter(new FileWriter(f));
      for (int i = 0; i < lines.length; i++) {
         out.println(lines[i]);
      }
      out.close();         
   }
   
   public void renameTo(File newName, boolean overwrite) throws IOException {
      if (newName.exists()) {
         if (!overwrite) {
            throw new IOException("Can't rename to an existing file without override");
         }
         newName.renameTo(new File(newName.getAbsolutePath() + ".OLD"));
      }
      _path = newName;
      write();
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
   
   
