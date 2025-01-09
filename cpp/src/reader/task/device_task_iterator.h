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
#ifndef READER_TASK_DEVICE_TASK_ITERATOR_H
#define READER_TASK_DEVICE_TASK_ITERATOR_H

#include "common/device_id.h"
#include "table_query_executor.h"
#include "imeta_data_querier.h"
namespace storage {
class DeviceTaskIterator {
   public:
    explicit DeviceTaskIterator(std::vector<std::string> column_names,
                                MetaIndexNode index_root,
                                ColumnMapping column_mapping,
                                IMetadataQuerier &metadata_querier,
                                Filter id_filter, TableSchema table_schema)
        : column_names_(column_names),
          column_mapping_(column_mapping),
          device_meta_iterator_(
              metadata_querier.deviceIterator(index_root, id_filter)),
          table_schema_(table_schema) {}
    ~DeviceTaskIterator();
    
    bool has_next() const { return device_meta_iterator_->has_next(); }

    DeviceQueryTask next() {
        auto &pair = device_meta_iterator_->next();
        return DeviceQueryTask(pair.first, column_names_, column_mapping_,
                               pair.second, table_schema_);
    }
   private:
    std::vector<std::string> column_names_;
    ColumnMapping column_mapping_;
    std::unique_ptr<std::vector<std::pair<IDeviceID, MetaIndexNode>>::iterator>
        device_meta_iterator_;
    TableSchema table_schema_;
};

}  // namespace storage

#endif  // READER_TASK_DEVICE_TASK_ITERATOR_H