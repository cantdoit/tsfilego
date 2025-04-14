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

package gowrapper

import (
	"fmt"
	"time"
)

// TsFile provides high-level operations for working with TsFiles
type TsFile struct {
	reader *Reader
	writer *Writer
}

// NewTsFile creates a new TsFile instance
func NewTsFile() *TsFile {
	return &TsFile{}
}

// OpenReader opens a TsFile for reading
func (tf *TsFile) OpenReader(path string, tableName string, columns []string,
	startTime, endTime *time.Time, batchSize *int) error {

	reader, err := NewReader(path, tableName, columns, startTime, endTime, batchSize)
	if err != nil {
		return fmt.Errorf("failed to open reader: %w", err)
	}
	tf.reader = reader
	return nil
}

// OpenWriter opens a TsFile for writing
func (tf *TsFile) OpenWriter(path string) error {
	writer, err := NewWriter(path)
	if err != nil {
		return fmt.Errorf("failed to open writer: %w", err)
	}
	tf.writer = writer
	return nil
}

// Read reads all data from the TsFile
func (tf *TsFile) Read() ([]map[string]interface{}, error) {
	if tf.reader == nil {
		return nil, fmt.Errorf("reader not initialized")
	}
	return tf.reader.Read()
}

// Next reads the next batch of data
func (tf *TsFile) Next() ([]map[string]interface{}, error) {
	if tf.reader == nil {
		return nil, fmt.Errorf("reader not initialized")
	}
	return tf.reader.Next()
}

// RegisterTimeseries registers a new timeseries schema
func (tf *TsFile) RegisterTimeseries(tableName, columnName, dataType string) error {
	if tf.writer == nil {
		return fmt.Errorf("writer not initialized")
	}
	return tf.writer.RegisterTimeseries(tableName, columnName, dataType)
}

// WriteRow writes a single row of data
func (tf *TsFile) WriteRow(tableName string, timestamp time.Time,
	data map[string]interface{}) error {

	if tf.writer == nil {
		return fmt.Errorf("writer not initialized")
	}
	return tf.writer.WriteRow(tableName, timestamp, data)
}

// Flush ensures all data is written to disk
func (tf *TsFile) Flush() error {
	if tf.writer == nil {
		return fmt.Errorf("writer not initialized")
	}
	return tf.writer.Flush()
}

// Close closes any open readers or writers
func (tf *TsFile) Close() error {
	var err error

	if tf.reader != nil {
		if e := tf.reader.Close(); e != nil {
			err = fmt.Errorf("reader close error: %w", e)
		}
		tf.reader = nil
	}

	if tf.writer != nil {
		if e := tf.writer.Close(); e != nil {
			if err != nil {
				err = fmt.Errorf("%v, writer close error: %w", err, e)
			} else {
				err = fmt.Errorf("writer close error: %w", e)
			}
		}
		tf.writer = nil
	}

	return err
}

// OpenForReading Helper function to create a new TsFile with reader
func OpenForReading(path string, tableName string, columns []string,
	startTime, endTime *time.Time, batchSize *int) (*TsFile, error) {

	tf := NewTsFile()
	if err := tf.OpenReader(path, tableName, columns, startTime, endTime, batchSize); err != nil {
		return nil, err
	}
	return tf, nil
}

// Helper function to create a new TsFile with writer
func OpenForWriting(path string) (*TsFile, error) {
	tf := NewTsFile()
	if err := tf.OpenWriter(path); err != nil {
		return nil, err
	}
	return tf, nil
}
