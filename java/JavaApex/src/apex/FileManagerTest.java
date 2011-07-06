package apex;

import java.io.File;
import java.io.IOException;

import static org.junit.Assert.*;
import org.junit.Before;
import org.junit.Test;

public class FileManagerTest {
   FileManager _mgr;
   
   @Before
   public void setUp() {
      _mgr = new FileManager();
   }
   
   @Test
   public void testOpen() throws Exception {
      GapBuffer buf = _mgr.openBuffer(new File("testdata/test.txt"), false);
      assertEquals("hello\nthere\n", buf.allText());
   }
   
   @Test
   public void testOpenFileNotFound() {
      try {
         _mgr.openBuffer(new File("testdata/not_there.txt"), false);
         fail("Should have thrown an exception!");
      } catch (IOException e) {
         // Good
      }
   }
   
   @Test
   public void testCreateFile() throws IOException {
      GapBuffer buf = _mgr.openBuffer(new File("testdata/not_there_yet.txt"), true);
      assertNotNull(buf);
      buf.insert("hello\nthere\ntest");
      buf.write();
      assertTrue(buf.getPath().exists());
      buf.getPath().delete();
   }
   
   @Test
   public void testBackup() throws IOException {
      File file = new File("testdata/not_there_yet.txt");
      File backup = new File("testdata/not_there_yet.txt.BAK");
      try {
         GapBuffer buf = _mgr.openBuffer(file, true);
         assertNotNull(buf);
         buf.insert("hello\nthere\ntest");
         buf.write();
         assertTrue(buf.getPath().exists());
         assertFalse(backup.exists());
         buf.insert("foobar");
         buf.write();
         assertTrue(backup.exists());
         assertTrue(backup.length() < buf.getPath().length());
      } finally {
         if (file.exists()) {
            file.delete();
         }
         if (backup.exists()) {
            backup.delete();
         }
      }
   }
}
