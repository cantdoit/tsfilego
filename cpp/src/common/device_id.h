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

#ifndef COMMON_DEVICE_ID_H
#define COMMON_DEVICE_ID_H

#include <algorithm>
#include <cstdint>
#include <cstring>
#include <functional>
#include <iostream>
#include <memory>
#include <numeric>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

#include "common/allocator/byte_stream.h"
#include "utils/errno_define.h"

class IDeviceID {
   public:
    virtual ~IDeviceID() = default;
    virtual int serialize(common::ByteStream& write_stream) { return 0; }
    virtual std::vector<uint8_t> get_bytes() { return {}; }
    virtual bool is_empty() { return false; }
    virtual bool is_table_model() { return false; }
    virtual std::string get_table_name() { return ""; }
    virtual int segment_num() { return 0; }
    virtual std::string segment(int i) { return ""; }
    virtual int serialized_size() { return 0; }
    virtual bool start_with(const std::string& prefix,
                            bool match_entire_segment = false) {
        return false;
    }
    virtual std::vector<std::string> get_segments() const { return {}; }
    virtual bool match_database_name(const std::string& database_name) {
        return false;
    }
    virtual std::string get_device_name() const { return ""; };
    virtual bool operator<(const IDeviceID& other) { return 0; }
    virtual bool operator==(const IDeviceID& other) { return false; }
    virtual bool operator!=(const IDeviceID& other) { return false; }
};

struct IDeviceIDComparator {
    bool operator()(const std::shared_ptr<IDeviceID>& lhs,
                    const std::shared_ptr<IDeviceID>& rhs) const {
        return *lhs < *rhs;
    }
};

class StringArrayDeviceID : public IDeviceID {
   public:
    explicit StringArrayDeviceID(const std::vector<std::string>& segments)
        : segments_(formalize(segments)) {}

    explicit StringArrayDeviceID(const std::string& device_id_string)
        : segments_(split_device_id_string(device_id_string)) {}

    ~StringArrayDeviceID() {}

    std::string get_device_name() const override {
        return std::accumulate(std::next(segments_.begin()), segments_.end(),
                               segments_.front(),
                               [](std::string a, const std::string& b) {
                                   return std::move(a) + "." + b;
                               });
    };

    int serialize(common::ByteStream& write_stream) override {
        int ret = common::E_OK;
        if (RET_FAIL(common::SerializationUtil::write_var_int(segment_num(),
                                                              write_stream))) {
            return ret;
        }
        for (const auto& segment : segments_) {
            if (RET_FAIL(common::SerializationUtil::write_var_int(
                    segment.size(), write_stream))) {
                return ret;
            } else if (RET_FAIL(write_stream.write_buf(segment.c_str(),
                                                       segment.size()))) {
                return ret;
            }
        }
        return ret;
    }

    bool is_empty() override { return segments_.empty(); }

    bool is_table_model() override {
        return !segments_.empty() &&
               segments_[0].find('.') == std::string::npos;
    }

    std::string get_table_name() override {
        return segments_.empty() ? "" : segments_[0];
    }

    int segment_num() override { return static_cast<int>(segments_.size()); }

    std::string segment(int i) override {
        if (i < 0 || i >= static_cast<int>(segments_.size())) {
            throw std::out_of_range("segment index out of range");
        }
        return segments_[i];
    }

    int serialized_size() override {
        int size = sizeof(uint32_t);
        for (const auto& segment : segments_) {
            size += sizeof(uint32_t) + static_cast<int>(segment.size());
        }
        return size;
    }

    bool start_with(const std::string& prefix,
                    bool match_entire_segment = false) override {
        size_t matched_pos = 0;
        for (const auto& segment : segments_) {
            if (segment.compare(0, prefix.size() - matched_pos, prefix,
                                matched_pos) == 0) {
                return true;
            }
            matched_pos += segment.size() + 1;
            if (matched_pos >= prefix.size()) {
                return false;
            }
        }
        return false;
    }

    std::vector<std::string> get_segments() const override { return segments_; }

    bool match_database_name(const std::string& database_name) override {
        std::string table_name = get_table_name();
        return table_name.find(database_name) == 0;
    }

    virtual bool operator<(const IDeviceID& other) override {
        auto other_segments = other.get_segments();
        return std::lexicographical_compare(segments_.begin(), segments_.end(),
                                            other_segments.begin(),
                                            other_segments.end());
    }

    virtual bool operator==(const IDeviceID& other) override {
        auto other_segments = other.get_segments();
        return (segments_.size() == other_segments.size()) &&
               std::equal(segments_.begin(), segments_.end(),
                          other_segments.begin());
    }

    virtual bool operator!=(const IDeviceID& other) override {
        return !(*this == other);
    }

   private:
    std::vector<std::string> segments_;

    static std::vector<std::string> formalize(
        const std::vector<std::string>& segments) {
        auto it =
            std::find_if(segments.rbegin(), segments.rend(),
                         [](const std::string& seg) { return !seg.empty(); });
        return std::vector<std::string>(segments.begin(), it.base());
    }

    static std::vector<std::string> split_device_id_string(
        std::basic_string<char> device_id_string) {
        std::vector<std::string> splits;
        std::istringstream stream(device_id_string);
        std::string segment;
        while (std::getline(stream, segment, '.')) {
            splits.push_back(segment);
        }
        return splits;
    }
};

class PlainDeviceID : public IDeviceID {
   public:
    explicit PlainDeviceID(const std::string& deviceID)
        : device_id_(deviceID), tableName_(), segments_() {}

    ~PlainDeviceID() {}

    bool operator==(const IDeviceID& other) override {
        return device_id_ == other.get_device_name();
    }

    bool operator!=(const IDeviceID& other) override {
        return device_id_ == other.get_device_name();
    }

    int serialize(common::ByteStream& write_stream) override {
        int ret = common::E_OK;
        if (RET_FAIL(common::SerializationUtil::write_var_int(device_id_.size(),
                                                              write_stream))) {
            return ret;
        } else if (RET_FAIL(write_stream.write_buf(device_id_.c_str(),
                                                   device_id_.size()))) {
            return ret;
        }
        return ret;
    }

    std::string get_device_name() const override { return device_id_; };

    bool is_empty() override { return device_id_.empty(); }

    bool is_table_model() override { return false; }

    std::string get_table_name() override {
        if (!tableName_.empty()) {
            return tableName_;
        }

        size_t lastSeparatorPos = device_id_.find_last_of('.');
        if (lastSeparatorPos == std::string::npos) {
            tableName_ = device_id_;  // Use entire deviceID as tableName
        } else {
            tableName_ = device_id_.substr(0, lastSeparatorPos);
        }
        return tableName_;
    }

    int segment_num() override {
        if (!segments_.empty()) {
            return static_cast<int>(segments_.size());
        }
        split_segments();
        return static_cast<int>(segments_.size());
    }

    std::string segment(int i) override {
        if (i < 0 || i >= segment_num()) {
            throw std::out_of_range("segment index out of range");
        }
        return segments_[i];
    }

    bool operator<(const IDeviceID& other) override {
        return device_id_ < other.get_device_name();
    }

   private:
    std::string device_id_;
    std::string tableName_;
    std::vector<std::string> segments_;

    void split_segments() {
        std::istringstream stream(device_id_);
        std::string segment;
        while (std::getline(stream, segment, '.')) {
            segments_.push_back(segment);
        }
    }
};

class PlainDeviceIDFactory {
   public:
    static std::shared_ptr<IDeviceID> create(
        const std::string& deviceIdString) {
        return std::make_shared<PlainDeviceID>(deviceIdString);
    }

    static std::shared_ptr<IDeviceID> create(
        const std::vector<std::string>& segments) {
        return std::make_shared<PlainDeviceID>(join_segments(segments));
    }

   private:
    static std::string join_segments(const std::vector<std::string>& segments) {
        return std::accumulate(segments.begin(), segments.end(), std::string(),
                               [](const std::string& a, const std::string& b) {
                                   return a.empty() ? b : a + "." + b;
                               });
    }
};

#endif