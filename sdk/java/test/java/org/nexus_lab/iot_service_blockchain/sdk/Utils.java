package org.nexus_lab.iot_service_blockchain.sdk;

import java.util.function.Function;
import org.hyperledger.fabric.client.CloseableIterator;

class Utils {
  static <T> CloseableIterator<T> createIterator(int count, Function<Integer, T> creator) {
    return new CloseableIterator<T>() {
      private int index;

      @Override
      public boolean hasNext() {
        return index < count;
      }

      @Override
      public T next() {
        return creator.apply(index++);
      }

      @Override
      public void close() {}
    };
  }
}
