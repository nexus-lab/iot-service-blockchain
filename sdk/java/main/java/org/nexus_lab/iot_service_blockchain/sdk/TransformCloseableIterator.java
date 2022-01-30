package org.nexus_lab.iot_service_blockchain.sdk;

import java.util.NoSuchElementException;
import org.hyperledger.fabric.client.CloseableIterator;

/**
 * A class that transform {@link org.hyperledger.fabric.client.CloseableIterator} of one type to
 * {@link org.hyperledger.fabric.client.CloseableIterator} of another type
 */
public abstract class TransformCloseableIterator<U, V> implements CloseableIterator<V> {
  private final CloseableIterator<U> sourceIterator;

  private U next;
  private boolean isClosed;
  private boolean isChecked;

  /** @param iterator source iterator to be transformed */
  public TransformCloseableIterator(CloseableIterator<U> iterator) {
    this.sourceIterator = iterator;
  }

  @Override
  public boolean hasNext() {
    this.next = null;

    if (this.isClosed) {
      return false;
    }

    while (this.sourceIterator.hasNext()) {
      U next = this.sourceIterator.next();
      if (this.canTransform(next)) {
        this.next = next;
        break;
      }
    }

    this.isChecked = true;
    return this.next != null;
  }

  @Override
  public V next() {
    // make sure hasNext() is always called before next()
    if (!this.isChecked) {
      this.hasNext();
    }
    this.isChecked = false;

    if (this.next == null) {
      throw new NoSuchElementException();
    }

    return this.transform(this.next);
  }

  @Override
  public void close() {
    this.isClosed = true;
    this.sourceIterator.close();
  }

  /**
   * Check if the next element of source type can be transformed into the element of target type
   *
   * @param element next element from the source iterator
   * @return if element can be transformed
   */
  public abstract boolean canTransform(U element);

  /**
   * Transform the next element of source type into the element of target type
   *
   * @param element next element from the source iterator
   * @return the transformed element
   */
  public abstract V transform(U element);
}
