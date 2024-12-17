/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.tsfile.utils;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

import static org.apache.tsfile.utils.RamUsageEstimator.shallowSizeOfInstance;
import static org.apache.tsfile.utils.RamUsageEstimator.sizeOf;

/**
 * This class represents a pooled binary object for application layer. It is designed to improve
 * allocation performance and reduce GC overhead by reusing binary objects from a pool instead of
 * creating new instances each time. WARNING: The actual length of the binary may not equal to the
 * length of the underlying byte array. Always use getLength() instead of getValue().length to get
 * the correct length.
 */
public class PooledBinary extends Binary {

  private static final long INSTANCE_SIZE = shallowSizeOfInstance(PooledBinary.class);
  private static final long serialVersionUID = 6394197743397020735L;

  private int length;

  private int arenaIndex = -1;

  /** if the bytes v is modified, the modification is visible to this binary. */
  public PooledBinary(byte[] v) {
    super(v);
    this.length = values.length;
  }

  public PooledBinary(String s, Charset charset) {
    super(s, charset);
    this.length = values.length;
  }

  public PooledBinary(byte[] v, int length, int arenaIndex) {
    super(v);
    this.length = length;
    this.arenaIndex = arenaIndex;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    PooledBinary otherBinary = (PooledBinary) o;

    if (length != otherBinary.length) {
      return false;
    }

    byte[] v0 = values;
    byte[] v1 = otherBinary.values;

    for (int i = 0; i < length; i++) {
      if (v0[i] != v1[i]) {
        return false;
      }
    }

    return true;
  }

  @Override
  public int hashCode() {
    // copied from Arrays.hashCode
    if (values == null) return 0;

    int result = 1;
    for (int i = 0; i < length; i++) {
      result = 31 * result + values[i];
    }

    return result;
  }

  @Override
  public int getLength() {
    return this.length;
  }

  @Override
  public String getStringValue(Charset charset) {
    return new String(values, 0, length, charset);
  }

  @Override
  public String toString() {
    // use UTF_8 by default since toString do not provide parameter
    return getStringValue(StandardCharsets.UTF_8);
  }

  @Override
  public Pair<byte[], Integer> getValuesAndLength() {
    return new Pair<>(values, length);
  }

  @Override
  public void setValues(byte[] values) {
    super.setValues(values);
    this.length = values.length;
  }

  public void setValues(byte[] values, int length) {
    super.setValues(values);
    this.length = length;
  }

  public int getArenaIndex() {
    return this.arenaIndex;
  }

  @Override
  public long ramBytesUsed() {
    return INSTANCE_SIZE + sizeOf(values);
  }

  @Override
  public long ramShallowBytesUsed() {
    return INSTANCE_SIZE;
  }
}
