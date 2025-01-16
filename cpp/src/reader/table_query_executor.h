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
#ifndef READER_TABLE_QUERY_EXECUTOR_H
#define READER_TABLE_QUERY_EXECUTOR_H

#include "block/device_ordered_tsblock_reader.h"
#include "block/tsblock_reader.h"
#include "common/schema.h"
#include "expression.h"
#include "ichunk_reader.h"
#include "imeta_data_querier.h"
#include "task/device_task_iterator.h"
#include "utils/errno_define.h"

namespace storage {

class ColumnMapping {
   public:
    int add(const std::string &column_name, int index, TableSchema &schema) {
        int columnIndex = schema.find_column_index(column_name);
        if (columnIndex < 0) {
            return common::E_COLUMN_NOT_EXIST;
        }

        ColumnCategory columnCategory =
            schema.get_column_categories()[columnIndex];
        columnPosMap[column_name].push_back(index);

        if (columnCategory == ColumnCategory::TAG) {
            tag_columns_.insert(column_name);
        } else {
            field_columns_.insert(column_name);
        }

        return common::E_OK;
    }

    int add(const Expression &measurementFilter) {
        // TODO: get measurements in the filter and add them to
        // field_columns_
        return common::E_OK;
    }

    const std::vector<int> &getColumnPos(const std::string &column_name) const {
        static const std::vector<int> empty;
        auto it = columnPosMap.find(column_name);
        return it != columnPosMap.end() ? it->second : empty;
    }

    bool is_tag(const std::string &column_name) const {
        return tag_columns_.find(column_name) != tag_columns_.end();
    }

    bool is_field(const std::string &column_name) const {
        return field_columns_.find(column_name) != field_columns_.end();
    }

    const std::unordered_set<std::string> &get_id_columns() const {
        return tag_columns_;
    }

    const std::unordered_set<std::string> &get_measurement_columns() const {
        return field_columns_;
    }

   private:
    std::unordered_map<std::string, std::vector<int>> columnPosMap;
    std::unordered_set<std::string> tag_columns_;
    std::unordered_set<std::string> field_columns_;
};

class TableQueryExecutor {
   public:
    enum class TableQueryOrdering { TIME, DEVICE };

    TableQueryExecutor(std::shared_ptr<IMetadataQuerier> meta_data_querier,
                       std::shared_ptr<IChunkReader> chunk_reader_,
                       TableQueryOrdering table_query_ordering,
                       int block_size = 1024)
        : meta_data_querier_(std::move(meta_data_querier)),
          chunk_reader_(std::move(chunk_reader_)),
          table_query_ordering_(table_query_ordering),
          block_size_(block_size) {}

    int query(const std::string &table_name,
              const std::vector<std::string> &columns,
              const Expression &time_filter, const Expression &id_filter,
              const Expression &measurement_filter,
              std::unique_ptr<TsBlockReader> &ret_reader) {
        int ret = common::E_OK;
        TsFileMeta* file_metadata;
        if (RET_FAIL(meta_data_querier_->get_whole_file_metadata(file_metadata))) {
            return ret;
        }
        common::PageArena pa;  // TODO: Optimize the memory allocation, use pa
                               // only to alloc String is not good
        pa.init(512, common::MOD_TSFILE_READER);
        common::String table_name_str;
        table_name_str.dup_from(table_name, pa);
        MetaIndexNode *table_root = nullptr;
        std::shared_ptr<TableSchema> table_schema;
        if (RET_FAIL(file_metadata->get_table_metaindex_node(table_name_str,
                                                            table_root))) {
        } else if (RET_FAIL(file_metadata->get_table_schema(table_name,
                                                           table_schema))) {
        }

        if (IS_FAIL(ret)) {
            ret_reader = std::unique_ptr<EmptyTsBlockReader>();
            return ret;
        }

        ColumnMapping columnMapping;
        for (size_t i = 0; i < columns.size(); ++i) {
            columnMapping.add(columns[i], static_cast<int>(i), *table_schema);
        }
        columnMapping.add(measurement_filter);

        auto deviceTaskIterator = std::unique_ptr<DeviceTaskIterator>(new DeviceTaskIterator(
            columns, *table_root, columnMapping, *meta_data_querier_, id_filter,
            *table_schema));

        switch (table_query_ordering_) {
            case TableQueryOrdering::DEVICE:
                ret_reader = std::make_unique<DeviceOrderedTsBlockReader>(
                    std::move(deviceTaskIterator), *meta_data_querier_,
                    *chunk_reader_, time_filter, measurement_filter,
                    block_size_);
            case TableQueryOrdering::TIME:
            default:
                ret = common::E_UNSUPPORTED_ORDER;
        }

        return ret;
    }

   private:
    std::shared_ptr<IMetadataQuerier> meta_data_querier_;
    std::shared_ptr<IChunkReader> chunk_reader_;
    TableQueryOrdering table_query_ordering_;
    int block_size_;
};

}  // namespace storage

#endif  // READER_TABLE_QUERY_EXECUTOR_H