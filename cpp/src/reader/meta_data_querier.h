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

#ifndef READER_META_DATA_QUERIER_H
#define READER_META_DATA_QUERIER_H

#include "common/cache/lru_cache.h"
#include "common/device_id.h"
#include "file/tsfile_io_reader.h"
#include "imeta_data_querier.h"

namespace storage {

class MetadataQuerier : IMetadataQuerier {
   public:
    static constexpr int CACHED_ENTRY_NUMBER = 1000;

    enum class LocateStatus { BEFORE, IN, AFTER };

    explicit MetadataQuerier(std::shared_ptr<TsFileIOReader> tsfile_io_reader)
        : tsfile_io_reader_(std::move(tsfile_io_reader)) {
        file_metadata_ = tsfile_io_reader_->get_tsfile_meta();
        device_chunk_meta_cache_ = std::make_unique<
            common::Cache<std::pair<IDeviceID, std::string>,
                          std::vector<std::shared_ptr<ChunkMeta>>, std::mutex>>(
            CACHED_ENTRY_NUMBER, CACHED_ENTRY_NUMBER / 10);
    }

    std::vector<std::shared_ptr<ChunkMeta>> get_chunk_metadata_list(
        const Path& timeseriesPath);

    std::vector<std::vector<std::shared_ptr<ChunkMeta>>> get_chunk_meta_lists(
        const IDeviceID& device_id,
        const std::set<std::string>& measurement_names,
        const MetaIndexNode& measurement_node);

    TsFileMeta get_whole_file_meta_data() const;

    void load_chunk_meta_datas(const std::vector<Path>& paths);

    common::TSDataType get_data_type(const Path& path) const;

    std::vector<TimeRange> convert_space_to_time_partition(
        const std::vector<Path>& paths, long spacePartitionStartPos,
        long spacePartitionEndPos);

    std::unique_ptr<

    void clear();

   private:
    std::shared_ptr<TsFileIOReader> tsfile_io_reader_;
    TsFileMeta* file_metadata_;
    std::unique_ptr<
        common::Cache<std::pair<IDeviceID, std::string>,
                      std::vector<std::shared_ptr<ChunkMeta>>, std::mutex>>
        device_chunk_meta_cache_;

    int load_chunk_meta(const std::pair<IDeviceID, std::string>& key,
                        std::vector<ChunkMeta*>& chunk_meta_list);

    static LocateStatus checkLocateStatus(
        const std::shared_ptr<ChunkMeta>& chunk_meta, long start, long end);
};

}  // end namespace storage
#endif  // READER_META_DATA_QUERIER_H
