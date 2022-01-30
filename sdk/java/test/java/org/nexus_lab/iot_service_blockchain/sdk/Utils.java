package org.nexus_lab.iot_service_blockchain.sdk;

import java.util.function.Function;
import org.hyperledger.fabric.client.CloseableIterator;

public class Utils {
  public static <T> CloseableIterator<T> createIterator(int count, Function<Integer, T> creator) {
    return new CloseableIterator<T>() {
      private int i;

      @Override
      public boolean hasNext() {
        return i < count;
      }

      @Override
      public T next() {
        return creator.apply(i++);
      }

      @Override
      public void close() {}
    };
  }
}
