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
    std::map<std::string, std::string> props_;

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

    int serialize_to(common::ByteStream &out) {
        int ret = common::E_OK;
        if (RET_FAIL(
                common::SerializationUtil::write_str(measurement_name_, out))) {
        } else if (RET_FAIL(
                       common::SerializationUtil::write_i32(data_type_, out))) {
        } else if (RET_FAIL(
                       common::SerializationUtil::write_i32(encoding_, out))) {
        } else if (RET_FAIL(common::SerializationUtil::write_i32(
                       compression_type_, out))) {
        }
        if (ret == common::E_OK) {
            if (RET_FAIL(common::SerializationUtil::write_str(measurement_name_,
                                                              out))) {
                for (auto &prop : props_) {
                    if (RET_FAIL(common::SerializationUtil::write_str(
                            prop.first, out))) {
                    } else if (RET_FAIL(common::SerializationUtil::write_str(
                                   prop.second, out))) {
                    }
                    if (ret != common::E_OK) break;
                }
            }
        }
        return ret;
    }
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
          measurement_schemas_(measurement_schemas),
          column_categories_(column_categories) {
        int idx = 0;
        for (auto &measurement_schema : measurement_schemas_) {
            column_pos_index_.insert(
                std::make_pair(measurement_schema->measurement_name_, idx++));
        }
    }

    TableSchema(TableSchema &&other) noexcept
        : table_name_(std::move(other.table_name_)),
          measurement_schemas_(std::move(other.measurement_schemas_)),
          column_categories_(std::move(other.column_categories_)) {}

    TableSchema(const TableSchema &other) = default;

    int serialize_to(common::ByteStream &out) {
        int ret = common::E_OK;
        if (RET_FAIL(common::SerializationUtil::write_var_uint(
                measurement_schemas_.size(), out))) {
        } else {
            for (size_t i = 0; IS_SUCC(ret) && i < measurement_schemas_.size();
                 i++) {
                auto column_schema = measurement_schemas_[i];
                auto column_category = column_categories_[i];
                // column_schema-
                common::SerializationUtil::write_i32(
                    static_cast<int32_t>(column_category), out);
            }
        }
        return ret;
    }

    ~TableSchema() = default;

    std::string get_table_name() { return table_name_; }

    auto get_measurement_names() const {
        std::vector<std::string> ret;
        for (const auto &measurement_schema : measurement_schemas_) {
            ret.emplace_back(measurement_schema->measurement_name_);
        }
        return ret;
    }

    int find_column_index(std::string column_name) {
        std::string lower_case_column_name = to_lower(column_name);
        auto it = column_pos_index_.find(lower_case_column_name);
        if (it != column_pos_index_.end()) {
            return it->second;
        } else {
            int index = -1;
            for (size_t i = 0; i < measurement_schemas_.size(); ++i) {
                if (to_lower(measurement_schemas_[i]->measurement_name_) ==
                    lower_case_column_name) {
                    index = static_cast<int>(i);
                    break;
                }
            }
            column_pos_index_[lower_case_column_name] = index;
            return index;
        }
    }

    void update(ChunkGroupMeta *chunk_group_meta) {
        for (auto iter = chunk_group_meta->chunk_meta_list_.begin();
             iter != chunk_group_meta->chunk_meta_list_.end(); iter++) {
            auto &chunk_meta = iter.get();
            int column_idx =
                find_column_index(chunk_meta->measurement_name_.to_string());
            if (column_idx == -1) {
                auto measurement_schema = std::make_shared<MeasurementSchema>(
                    chunk_meta->measurement_name_.to_string(),
                    chunk_meta->data_type_, chunk_meta->encoding_,
                    chunk_meta->compression_type_);
                measurement_schemas_.emplace_back(measurement_schema);
                column_categories_.emplace_back(ColumnCategory::MEASUREMENT);
                column_pos_index_.insert(
                    std::make_pair(chunk_meta->measurement_name_.to_string(),
                                   measurement_schemas_.size() - 1));
            } else {
                auto origin_measurement_schema =
                    measurement_schemas_.at(column_idx);
                if (origin_measurement_schema->data_type_ !=
                    chunk_meta->data_type_) {
                    origin_measurement_schema->data_type_ =
                        common::TSDataType::STRING;
                }
            }
        }
    }

    auto get_data_types() const {
        std::vector<common::TSDataType> ret;
        for (const auto &measurement_schema : measurement_schemas_) {
            ret.emplace_back(measurement_schema->data_type_);
        }
        return ret;
    }

    auto get_column_categories() const { return column_categories_; }

    auto get_measurement_schemas() const { return measurement_schemas_; }

   private:
    std::string to_lower(const std::string &str) {
        std::string result;
        std::transform(str.begin(), str.end(), std::back_inserter(result),
                       [](unsigned char c) { return std::tolower(c); });
        return result;
    }
    std::string table_name_;
    std::vector<std::shared_ptr<MeasurementSchema>> measurement_schemas_;
    std::vector<ColumnCategory> column_categories_;
    std::map<std::string, int> column_pos_index_;
};

struct Schema {
    typedef std::unordered_map<std::string, std::shared_ptr<TableSchema>>
        TableSchemasMap;
    TableSchemasMap table_schema_map_;

    void update_table_schema(ChunkGroupMeta *chunk_group_meta) {
        std::shared_ptr<IDeviceID> device_id = chunk_group_meta->device_name_;
        auto table_name = device_id->get_table_name();
        if (table_schema_map_.find(table_name) == table_schema_map_.end()) {
            table_schema_map_[table_name] = std::make_shared<TableSchema>();
        }
        table_schema_map_[table_name]->update(chunk_group_meta);
    }
};

}  // end namespace storage
#endif  // COMMON_SCHEMA_H
