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

#ifndef COMMON_TABLET_H
#define COMMON_TABLET_H

#include <vector>

#include "common/config/config.h"
#include "common/container/bit_map.h"
#include "common/db_common.h"
#include "device_id.h"
#include "schema.h"

namespace storage {

class TabletRowIterator;
class TabletColIterator;

class Tablet {
   public:
    static const int DEFAULT_MAX_ROWS = 1024;

   public:
    Tablet(const std::string &device_name,
           std::vector<MeasurementSchema> *schema_vec,
           int max_rows = DEFAULT_MAX_ROWS)
        : max_rows_(max_rows),
          cur_row_size_(0),
          insert_target_name_(device_name),
          schema_vec_(schema_vec),
          timestamps_(nullptr),
          value_matrix_(nullptr),
          bitmaps_(nullptr) {
        ASSERT(device_name.size() >= 1);
        ASSERT(schema_vec != nullptr);
        ASSERT(max_rows > 0 && max_rows < (1 << 30));
        if (max_rows < 0) {
            ASSERT(false);
            max_rows_ = DEFAULT_MAX_ROWS;
        }
    }

    Tablet(const std::string &insert_target_name,
           const std::vector<std::string> &column_names,
           const std::vector<common::TSDataType> &data_types,
           const std::vector<ColumnCategory> &column_categories,
           int max_rows = DEFAULT_MAX_ROWS, bool has_column_categories = true)
        : max_rows_(max_rows),
          cur_row_size_(0),
          insert_target_name_(insert_target_name),
          timestamps_(nullptr),
          value_matrix_(nullptr),
          bitmaps_(nullptr) {
        schema_vec_ = new std::vector<MeasurementSchema>();
        for (size_t i = 0; i < column_names.size(); i++) {
            schema_vec_->emplace_back(
                MeasurementSchema(column_names[i], data_types[i], common::PLAIN,
                                  common::UNCOMPRESSED));
        }
        if (has_column_categories) {
            set_column_categories(column_categories);
        }
    }

    ~Tablet() { destroy(); }

    int init();
    void destroy();
    size_t get_column_count() const { return schema_vec_->size(); }
    int get_cur_row_size() const { return cur_row_size_; }
    void set_row_size(int row_size) { cur_row_size_ = row_size; }

    int set_timestamp(int row_index, int64_t timestamp);

    int set_value(int row_index, uint32_t schema_index, bool val);
    int set_value(int row_index, uint32_t schema_index, int32_t val);
    int set_value(int row_index, uint32_t schema_index, int64_t val);
    int set_value(int row_index, uint32_t schema_index, float val);
    int set_value(int row_index, uint32_t schema_index, double val);
    // int set_value(int row_index, int schema_index, double val);
    void *get_value(int row_index, uint32_t schema_index,
                    common::TSDataType &data_type) const;

    int set_value(int row_index, const std::string &measurement_name, bool val);
    int set_value(int row_index, const std::string &measurement_name,
                  int32_t val);
    int set_value(int row_index, const std::string &measurement_name,
                  int64_t val);
    int set_value(int row_index, const std::string &measurement_name,
                  float val);
    int set_value(int row_index, const std::string &measurement_name,
                  double val);
    // int set_value(int row_index, const std::string &measurement_name, double
    // val);
    void set_column_categories(
        const std::vector<ColumnCategory> &column_categories);
    std::unique_ptr<IDeviceID> get_device_id(int i) const;

    friend class TabletColIterator;
    friend class TsFileWriter;
    friend struct MeasurementNamesFromTablet;

   private:
    typedef std::map<std::string, int>::iterator SchemaMapIterator;

   private:
    int max_rows_;
    int cur_row_size_;
    std::string insert_target_name_;
    std::vector<MeasurementSchema> *schema_vec_;
    std::map<std::string, int> schema_map_;
    int64_t *timestamps_;
    void **value_matrix_;
    common::BitMap *bitmaps_;
    std::vector<ColumnCategory> column_categories_;
    std::vector<int> id_column_indexes_;
};

}  // end namespace storage
#endif  // COMMON_TABLET_H
