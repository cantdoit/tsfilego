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
#include "tsfile_reader.h"

#include "common/schema.h"
#include "filter/time_operator.h"
#include "tsfile_executor.h"

using namespace common;
using namespace storage;

namespace storage {

TsFileReader::TsFileReader() : read_file_(nullptr), tsfile_executor_(nullptr) {}

TsFileReader::~TsFileReader() { close(); }

int TsFileReader::open(const std::string &file_path) {
    int ret = E_OK;
    read_file_ = new storage::ReadFile;
    tsfile_executor_ = new storage::TsFileExecutor();
    if (RET_FAIL(read_file_->open(file_path))) {
        std::cout << "filed to open file " << ret << std::endl;
    } else if (RET_FAIL(tsfile_executor_->init(read_file_))) {
        std::cout << "filed to init " << ret << std::endl;
    }
    return ret;
}

int TsFileReader::close() {
    int ret = E_OK;
    if (tsfile_executor_ != nullptr) {
        delete tsfile_executor_;
        tsfile_executor_ = nullptr;
    }
    if (read_file_ != nullptr) {
        read_file_->close();
        delete read_file_;
        read_file_ = nullptr;
    }
    return ret;  // TO DO
}

int TsFileReader::query(QueryExpression *qe, ResultSet *&ret_qds) {
    return tsfile_executor_->execute(qe, ret_qds);
}

int TsFileReader::query(std::vector<std::string> &path_list, int64_t start_time,
                        int64_t end_time, ResultSet *&result_set) {
    Filter *time_filter = new TimeBetween(start_time, end_time, false);
    Expression *exp =
        new storage::Expression(storage::GLOBALTIME_EXPR, time_filter);
    std::vector<Path> path_list_vec;
    for (auto path : path_list) {
        uint32_t last_point_pos = path.find_last_of('.');
        if (last_point_pos <= 0) {
            return E_INVALID_PATH;
        }
        std::string device_name = path.substr(0, last_point_pos);
        std::string measurement_name =
            path.substr(last_point_pos + 1, path.size() - last_point_pos);
        path_list_vec.emplace_back(Path(device_name, measurement_name));
    }
    QueryExpression *query_expression =
        QueryExpression::create(path_list_vec, exp);
    return tsfile_executor_->execute(query_expression, result_set);
}

void TsFileReader::destroy_query_data_set(storage::ResultSet *qds) {
    tsfile_executor_->destroy_query_data_set(qds);
}

std::vector<std::string> TsFileReader::get_all_devices() {
    TsFileMeta *tsfile_meta = tsfile_executor_->get_tsfile_meta();
    std::vector<std::string> device_ids;
    if (tsfile_meta != nullptr) {
        device_ids.reserve(tsfile_meta->index_node_->children_.size());
        for (const auto &meta_index_entry :
             tsfile_meta->index_node_->children_) {
            device_ids.push_back(meta_index_entry->name_.to_std_string());
        }
    }
    return device_ids;
}

int TsFileReader::get_timeseries_schema(
    const std::string &device_id, std::vector<MeasurementSchema> &result) {
    int ret = E_OK;
    std::vector<ITimeseriesIndex *> timeseries_indexs;
    PageArena pa;
    pa.init(512, MOD_TSFILE_READER);
    if (RET_FAIL(tsfile_executor_->get_tsfile_io_reader()
                     ->get_device_timeseries_meta_without_chunk_meta(
                         device_id, timeseries_indexs, pa))) {
    } else {
        for (auto timeseries_index : timeseries_indexs) {
            MeasurementSchema ms(
                timeseries_index->get_measurement_name().to_std_string(),
                timeseries_index->get_data_type());
            result.push_back(ms);
        }
    }
    return E_OK;
}

ResultSet *TsFileReader::read_timeseries(
    const std::string &device_name, std::vector<std::string> measurement_name) {
    return nullptr;
}

}  // namespace storage
