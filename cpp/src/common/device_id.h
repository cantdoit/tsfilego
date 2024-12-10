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

class IDeviceID {
   public:
    virtual ~IDeviceID() = default;
    virtual int Serialize(std::ostream& output_stream) const = 0;
    virtual int Serialize(std::vector<uint8_t>& byte_buffer) const = 0;
    virtual std::vector<uint8_t> GetBytes() const = 0;
    virtual bool IsEmpty() const = 0;
    virtual bool IsTableModel() const = 0;
    virtual std::string GetTableName() const = 0;
    virtual int SegmentNum() const = 0;
    virtual std::string Segment(int i) const = 0;
    virtual int SerializedSize() const = 0;
    virtual bool StartWith(const std::string& prefix,
                           bool match_entire_segment = false) const = 0;
    virtual std::vector<std::string> GetSegments() const = 0;
    virtual bool MatchDatabaseName(const std::string& database_name) const = 0;
    virtual int CompareTo(const IDeviceID& other) const = 0;
};

class StringArrayDeviceID : public IDeviceID {
   public:
    explicit StringArrayDeviceID(const std::vector<std::string>& segments)
        : segments_(Formalize(segments)) {}

    explicit StringArrayDeviceID(const std::string& device_id_string)
        : segments_(SplitDeviceIdString(device_id_string)) {}

    int Serialize(std::ostream& output_stream) const override {
        int cnt = 0;
        uint32_t length = static_cast<uint32_t>(segments_.size());
        output_stream.write(reinterpret_cast<const char*>(&length),
                            sizeof(length));
        cnt += sizeof(length);
        for (const auto& segment : segments_) {
            uint32_t size = static_cast<uint32_t>(segment.size());
            output_stream.write(reinterpret_cast<const char*>(&size),
                                sizeof(size));
            output_stream.write(segment.data(), size);
            cnt += sizeof(size) + size;
        }
        return cnt;
    }

    int Serialize(std::vector<uint8_t>& byte_buffer) const override {
        std::ostringstream stream;
        int size = Serialize(stream);
        std::string str = stream.str();
        byte_buffer.assign(str.begin(), str.end());
        return size;
    }

    std::vector<uint8_t> GetBytes() const override {
        std::vector<uint8_t> buffer;
        Serialize(buffer);
        return buffer;
    }

    bool IsEmpty() const override { return segments_.empty(); }

    bool IsTableModel() const override {
        return !segments_.empty() &&
               segments_[0].find(".") == std::string::npos;
    }

    std::string GetTableName() const override {
        return segments_.empty() ? "" : segments_[0];
    }

    int SegmentNum() const override {
        return static_cast<int>(segments_.size());
    }

    std::string Segment(int i) const override {
        if (i < 0 || i >= static_cast<int>(segments_.size())) {
            throw std::out_of_range("Segment index out of range");
        }
        return segments_[i];
    }

    int SerializedSize() const override {
        int size = sizeof(uint32_t);
        for (const auto& segment : segments_) {
            size += sizeof(uint32_t) + static_cast<int>(segment.size());
        }
        return size;
    }

    bool StartWith(const std::string& prefix,
                   bool match_entire_segment = false) const override {
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

    std::vector<std::string> GetSegments() const override { return segments_; }

    bool MatchDatabaseName(const std::string& database_name) const override {
        std::string table_name = GetTableName();
        return table_name.find(database_name) == 0;
    }

    int CompareTo(const IDeviceID& other) const override {
        const auto& other_segments = other.GetSegments();
        return std::lexicographical_compare(segments_.begin(), segments_.end(),
                                            other_segments.begin(),
                                            other_segments.end())
                   ? -1
                   : (segments_ == other_segments ? 0 : 1);
    }

   private:
    std::vector<std::string> segments_;

    std::vector<std::string> Formalize(
        const std::vector<std::string>& segments) {
        auto it =
            std::find_if(segments.rbegin(), segments.rend(),
                         [](const std::string& seg) { return !seg.empty(); });
        return std::vector<std::string>(segments.begin(), it.base());
    }

    std::vector<std::string> SplitDeviceIdString(
        const std::string& device_id_string) {
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
        : deviceID_(deviceID), tableName_(), segments_() {}

    bool operator==(const PlainDeviceID& other) const {
        return deviceID_ == other.deviceID_;
    }

    bool operator!=(const PlainDeviceID& other) const {
        return !(*this == other);
    }

    int Serialize(std::ostream& output_stream) const override {
        uint32_t length = static_cast<uint32_t>(deviceID_.size());
        output_stream.write(reinterpret_cast<const char*>(&length),
                            sizeof(length));
        output_stream.write(deviceID_.data(), deviceID_.size());
        return sizeof(length) + deviceID_.size();
    }

    int Serialize(std::vector<uint8_t>& byte_buffer) const override {
        std::ostringstream stream;
        int size = Serialize(stream);
        std::string str = stream.str();
        byte_buffer.assign(str.begin(), str.end());
        return size;
    }

    std::vector<uint8_t> GetBytes() const override {
        std::vector<uint8_t> buffer;
        Serialize(buffer);
        return buffer;
    }

    bool IsEmpty() const override { return deviceID_.empty(); }

    bool IsTableModel() const override { return false; }

    std::string GetTableName() const override {
        if (!tableName_.empty()) {
            return tableName_;
        }

        size_t lastSeparatorPos = deviceID_.find_last_of('.');
        if (lastSeparatorPos == std::string::npos) {
            tableName_ = deviceID_;  // Use entire deviceID as tableName
        } else {
            tableName_ = deviceID_.substr(0, lastSeparatorPos);
        }
        return tableName_;
    }

    int SegmentNum() const override {
        if (!segments_.empty()) {
            return static_cast<int>(segments_.size());
        }
        SplitSegments();
        return static_cast<int>(segments_.size());
    }

    std::string Segment(int i) const override {
        if (i < 0 || i >= SegmentNum()) {
            throw std::out_of_range("Segment index out of range");
        }
        return segments_[i];
    }

    int CompareTo(const IDeviceID& other) const override {
        const PlainDeviceID* otherPlain =
            dynamic_cast<const PlainDeviceID*>(&other);
        if (!otherPlain) {
            throw std::invalid_argument("Incompatible IDeviceID type");
        }
        return deviceID_.compare(otherPlain->deviceID_);
    }

   private:
    std::string deviceID_;
    mutable std::string tableName_;
    mutable std::vector<std::string> segments_;

    void SplitSegments() const {
        std::istringstream stream(deviceID_);
        std::string segment;
        while (std::getline(stream, segment, '.')) {
            segments_.push_back(segment);
        }
    }
};

//class PlainDeviceIDFactory {
//   public:
//    static std::unique_ptr<IDeviceID> Create(
//        const std::string& deviceIdString) {
//        return std::make_unique<PlainDeviceID>(deviceIdString);
//    }
//
//    static std::unique_ptr<IDeviceID> Create(
//        const std::vector<std::string>& segments) {
//        return std::make_unique<PlainDeviceID>(JoinSegments(segments));
//    }
//
//   private:
//    static std::string JoinSegments(const std::vector<std::string>& segments) {
//        return std::accumulate(segments.begin(), segments.end(), std::string(),
//                               [](const std::string& a, const std::string& b) {
//                                   return a.empty() ? b : a + "." + b;
//                               });
//    }
//};

#endif