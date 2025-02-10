<!--

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at
    
        http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.

-->
# Interface Definitions

## Write Interface

### TsFileTableWriter

Used to write data to tsfile

```cpp
class TsFileTableWriter {
   public:
    /**
     * TsFileTableWriter is used to write table data into a target file with the given schema,
     * optionally limiting the memory usage.
     *
     * @param writer_file Target file where the table data will be written. Must not be null.
     *                    The caller retains ownership of this pointer.
     * @param table_schema Used to construct table structures. Defines the schema of the table
     *                     being written. Must not be null. The caller retains ownership of this pointer.
     * @param memory_threshold Optional parameter used to limit the memory size of objects.
     *                         If set to 0, no memory limit is enforced.
     */
    TsFileTableWriter(WriteFile* writer_file,
                      TableSchema* table_schema,
                      uint64_t memory_threshold = 0);
    ~TsFileTableWriter();
    /**
     * Writes the given tablet data into the target file according to the schema.
     *
     * @param tablet The tablet containing the data to be written. Must not be null.
     *               The caller retains ownership of this object.
     * @return Returns 0 on success, or a non-zero error code on failure.
     */
    int write_table(const Tablet& tablet);
    /**
     * Flushes any buffered data to the underlying storage medium, ensuring all data is written out.
     * This method ensures that all pending writes are persisted.
     *
     * @return Returns 0 on success, or a non-zero error code on failure.
     */
    int flush();
    /**
     * Closes the writer and releases any resources held by it.
     * After calling this method, no further operations should be performed on this instance.
     *
     * @return Returns 0 on success, or a non-zero error code on failure.
     */
    int close();
};
```

### TableSchema

Describe the data structure of the table schema

```cpp
class TableSchema {
public:
    /**
     * Constructs a TableSchema object with the given table name, column schemas, and column categories.
     *
     * @param table_name The name of the table. Must be a non-empty string.
     *                   This name is used to identify the table within the system.
     * @param column_schemas A vector containing pointers to MeasurementSchema objects.
     *                       Each MeasurementSchema defines the schema for one column in the table.
     *                       The caller retains ownership of these pointers.
     * @param column_categories A vector containing ColumnCategory enums that correspond to each column schema.
     *                          These categories provide additional information about how each column should be handled or optimized.
     * @note It is the responsibility of the caller to ensure that the sizes of `column_schemas` and `column_categories` match.
     */
    TableSchema(const std::string &table_name,
                const std::vector<MeasurementSchema*>
                &column_schemas,
                const std::vector<ColumnCategory> &column_categories);
};

struct MeasurementSchema {
    std::string measurement_name_; // for example: "s1"
    common::TSDataType data_type_;
    common::TSEncoding encoding_;
    common::CompressionType compression_type_;
    /**
     * @brief Constructs a MeasurementSchema object with the given parameters.
     *
     * @param measurement_name The name of the measurement. Must be a non-empty string.
     *                         This name is used to identify the measurement within the table.
     * @param data_type The data type of the measurement, such as INT32, DOUBLE, TEXT, etc.
     *                  This determines how the data will be stored and interpreted.
     * @param encoding The encoding method to use for this measurement.
     *                 Encoding can help optimize storage and retrieval performance.
     * @param compression_type The compression type to apply to this measurement.
     *                         Compression reduces storage space but may impact read/write performance.
     * @note It is the responsibility of the caller to ensure that `measurement_name` is not empty.
     */
    MeasurementSchema(const std::string &measurement_name,
                      common::TSDataType data_type, common::TSEncoding encoding,
                      common::CompressionType compression_type)
        : measurement_name_(measurement_name),
          data_type_(data_type),
          encoding_(encoding),
          compression_type_(compression_type);
};

/**
 * @brief Represents the data type of a measurement.
 *
 * This enumeration defines the supported data types for measurements in the system.
 */
enum TSDataType : uint8_t {
    BOOLEAN = 0,
    INT32 = 1,
    INT64 = 2,
    FLOAT = 3,
    DOUBLE = 4,
    TEXT = 5,
    VECTOR = 6,
    STRING = 11,
    NULL_TYPE = 254,
    INVALID_DATATYPE = 255
};

/**
 * @brief Represents the encoding method for a measurement.
 *
 * This enumeration defines the supported encoding methods that can be applied to measurements.
 */
enum TSEncoding : uint8_t {
    PLAIN = 0,
    DICTIONARY = 1,
    RLE = 2,
    DIFF = 3,
    TS_2DIFF = 4,
    BITMAP = 5,
    GORILLA_V1 = 6,
    REGULAR = 7,
    GORILLA = 8,
    ZIGZAG = 9,
    FREQ = 10,
    INVALID_ENCODING = 255
};

/**
 * @brief Represents the compression type for a measurement.
 *
 * This enumeration defines the supported compression methods that can be applied to measurements.
 */
enum CompressionType : uint8_t {
    UNCOMPRESSED = 0,
    SNAPPY = 1,
    GZIP = 2,
    LZO = 3,
    SDT = 4,
    PAA = 5,
    PLA = 6,
    LZ4 = 7,
    INVALID_COMPRESSION = 255
};
```

### Tablet

Write column memory structure

```cpp
/**
 * @brief Represents a collection of data rows with associated metadata for insertion into a table.
 *
 * This class is used to manage and organize data that will be inserted into a specific target table.
 * It handles the storage of timestamps and values, along with their associated metadata such as column names and types.
 */
class Tablet {
public:
    /**
     * @brief Constructs a Tablet object with the given parameters.
     *
     * @param insert_target_name The name of the target table where the data will be inserted.
     *                           Must be a non-empty string.
     * @param column_names A vector containing the names of the columns in the tablet.
     *                     Each name corresponds to a column in the target table.
     * @param data_types A vector containing the data types of each column.
     *                   These must match the schema of the target table.
     * @param column_categories A vector containing the categories (tag or field) of each column.
     *                          These provide additional information on how each column should be handled.
     * @param max_rows The maximum number of rows that this tablet can hold. Defaults to DEFAULT_MAX_ROWS.
     */
    Tablet(const std::string &insert_target_name,
           const std::vector<std::string> &column_names,
           const std::vector<common::TSDataType> &data_types,
           const std::vector<ColumnCategory> &column_categories,
           int max_rows = DEFAULT_MAX_ROWS);

    /**
     * @brief Adds a timestamp to the specified row.
     *
     * @param row_index The index of the row to which the timestamp will be added.
     *                  Must be less than the maximum number of rows.
     * @param timestamp The timestamp value to add.
     * @return Returns 0 on success, or a non-zero error code on failure.
     */
    int add_timestamp(uint32_t row_index, int64_t timestamp);

    /**
     * @brief Template function to add a value of type T to the specified row and column.
     *
     * @tparam T The type of the value to add.
     * @param row_index The index of the row to which the value will be added.
     *                  Must be less than the maximum number of rows.
     * @param schema_index The index of the column schema corresponding to the value being added.
     * @param val The value to add.
     * @return Returns 0 on success, or a non-zero error code on failure.
     */
    template <typename T>
    int add_value(uint32_t row_index, uint32_t schema_index, T val);

    /**
     * @brief Template function to add a value of type T to the specified row and column by name.
     *
     * @tparam T The type of the value to add.
     * @param row_index The index of the row to which the value will be added.
     *                  Must be less than the maximum number of rows.
     * @param measurement_name The name of the column to which the value will be added.
     *                         Must match one of the column names provided during construction.
     * @param val The value to add.
     * @return Returns 0 on success, or a non-zero error code on failure.
     */
    template <typename T>
    int add_value(uint32_t row_index, const std::string &measurement_name, T val);

    /**
     * @brief Initializes the Tablet object.
     *
     * This method performs any necessary setup before the tablet can be used.
     * @return Returns 0 on success, or a non-zero error code on failure.
     */
    int init();
};
```