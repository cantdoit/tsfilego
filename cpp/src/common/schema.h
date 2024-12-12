/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * License); you may not use this file except in compliance
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

#ifndef COMMON_SCHEMA_H
#define COMMON_SCHEMA_H

#include <map>  // use unordered_map instead
#include <memory>
#include <string>

#include "common/db_common.h"
#include "writer/time_chunk_writer.h"
#include "writer/value_chunk_writer.h"

namespace storage {
class ChunkWriter;
}

namespace storage {

/* schema information for one measurement */
struct MeasurementSchema {
    std::string measurement_name_;  // for example: "s1"
    common::TSDataType data_type_;
    common::TSEncoding encoding_;
    common::CompressionType compression_type_;
    storage::ChunkWriter *chunk_writer_;
    ValueChunkWriter *value_chunk_writer_;

    MeasurementSchema()
        : measurement_name_(),
          data_type_(common::INVALID_DATATYPE),
          encoding_(common::INVALID_ENCODING),
          compression_type_(common::INVALID_COMPRESSION),
          chunk_writer_(nullptr),
          value_chunk_writer_(nullptr) {}

    MeasurementSchema(const std::string &measurement_name,
                      common::TSDataType data_type, common::TSEncoding encoding,
                      common::CompressionType compression_type)
        : measurement_name_(measurement_name),
          data_type_(data_type),
          encoding_(encoding),
          compression_type_(compression_type),
          chunk_writer_(nullptr),
          value_chunk_writer_(nullptr) {}
};

typedef std::map<std::string, MeasurementSchema *> MeasurementSchemaMap;
typedef std::map<std::string, MeasurementSchema *>::iterator
    MeasurementSchemaMapIter;
typedef std::pair<MeasurementSchemaMapIter, bool>
    MeasurementSchemaMapInsertResult;

/* schema information for a device */
struct MeasurementSchemaGroup {
    // measurement_name -> MeasurementSchema
    MeasurementSchemaMap measurement_schema_map_;
    bool is_aligned_ = false;
    TimeChunkWriter *time_chunk_writer_ = nullptr;
};

enum class ColumnCategory { ID, MEASUREMENT };

class TableSchema {
   public:
    TableSchema() = default;
    TableSchema(const std::string &table_name,
                const std::vector<std::shared_ptr<MeasurementSchema>>
                    &measurement_schemas,
                const std::vector<ColumnCategory> &column_categories)
        : table_name_(table_name),
          measurementSchemas_(measurement_schemas),
          columnCategories_(column_categories) {}

    TableSchema(TableSchema &&other) noexcept
        : table_name_(std::move(other.table_name_)),
          measurementSchemas_(std::move(other.measurementSchemas_)),
          columnCategories_(std::move(other.columnCategories_)) {}

    TableSchema(const TableSchema &other) = default;

    ~TableSchema() = default;

    std::string get_table_name() { return table_name_; }

   private:
    std::string table_name_;
    std::vector<std::shared_ptr<MeasurementSchema>> measurementSchemas_;
    std::vector<ColumnCategory> columnCategories_;
};

}  // end namespace storage
#endif  // COMMON_SCHEMA_H
