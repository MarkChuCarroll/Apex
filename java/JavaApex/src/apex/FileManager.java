package apex;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class FileManager {
   private List<GapBuffer> _buffers = new ArrayList<GapBuffer>();

   /**
    * Get the buffer associated with a file.
    * 
    * @param file
    * @return
    */
   GapBuffer getBuffer(File file) {
      for (GapBuffer buf : _buffers) {
         if (file.equals(buf.getPath())) {
            return buf;
         }
      }
      return null;
   }

   /**
    * Open a new file.
    * 
    * @param f
    * @return
    * @throws IOException
    */
   GapBuffer openBuffer(File f, boolean create) throws IOException {
      GapBuffer gap = getBuffer(f);
      if (gap != null) {
         return gap;
      }
      return new GapBuffer(f, create);
   }

}
