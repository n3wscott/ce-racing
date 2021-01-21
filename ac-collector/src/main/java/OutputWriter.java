import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.concurrent.ConcurrentLinkedQueue;

public class OutputWriter implements Runnable {

  private final ConcurrentLinkedQueue<byte[]> linkedQueue;
  private final File datOutputFile;

  public OutputWriter(ConcurrentLinkedQueue<byte[]> linkedQueue, String fileOutput) {
    this.linkedQueue = linkedQueue;
    this.datOutputFile = new File(fileOutput);
  }

  @Override
  public void run() {
    try (FileOutputStream stream = new FileOutputStream(datOutputFile)) {
      while (!Thread.currentThread().isInterrupted()) {
        byte[] data = this.linkedQueue.poll();
        if (data == null) {
          continue;
        }
        stream.write(data);
        stream.flush();
      }
    } catch (IOException e) {
      e.printStackTrace();
    }
  }

}
