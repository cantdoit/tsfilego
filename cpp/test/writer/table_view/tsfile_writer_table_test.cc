/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * License); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License a
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
#include <gtest/gtest.h>

#include <random>

#include "common/path.h"
#include "common/record.h"
#include "common/schema.h"
#include "common/tablet.h"
#include "file/tsfile_io_writer.h"
#include "file/write_file.h"
#include "reader/qds_without_timegenerator.h"
#include "reader/tsfile_reader.h"
#include "writer/chunk_writer.h"
#include "writer/tsfile_writer.h"
using namespace storage;
using namespace common;

class TsFileWriterTableTest : public ::testing::Test {
   protected:
    void SetUp() override {
        libtsfile_init();
        tsfile_writer_ = new TsFileWriter();
        file_name_ = std::string("tsfile_writer_table_test_") +
                     generate_random_string(10) + std::string(".tsfile");
        remove(file_name_.c_str());
        int flags = O_WRONLY | O_CREAT | O_TRUNC;
#ifdef _WIN32
        flags |= O_BINARY;
#endif
        mode_t mode = 0666;
        EXPECT_EQ(tsfile_writer_->open(file_name_, flags, mode), common::E_OK);
    }
    void TearDown() override {
        delete tsfile_writer_;
        remove(file_name_.c_str());
    }
    std::string file_name_;
    TsFileWriter* tsfile_writer_ = nullptr;

   public:
    static std::string generate_random_string(int length) {
        std::random_device rd;
        std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(0, 61);

        const std::string chars =
            "0123456789"
            "abcdefghijklmnopqrstuvwxyz"
            "ABCDEFGHIJKLMNOPQRSTUVWXYZ";

        std::string random_string;

        for (int i = 0; i < length; ++i) {
            random_string += chars[dis(gen)];
        }

        return random_string;
    }

    static std::shared_ptr<TableSchema> gen_table_schema(int table_num) {
        std::vector<std::shared_ptr<MeasurementSchema>> measurement_schemas;
        std::vector<ColumnCategory> column_categories;
        int id_schema_num = 5;
        int measurement_schema_num = 5;
        for (int i = 0; i < id_schema_num; i++) {
            // TODO: support TEXT
            measurement_schemas.emplace_back(
                std::make_shared<MeasurementSchema>(
                    "id" + to_string(i), TSDataType::INT64, TSEncoding::PLAIN,
                    CompressionType::UNCOMPRESSED));
            column_categories.emplace_back(ColumnCategory::ID);
        }
        for (int i = 0; i < measurement_schema_num; i++) {
            measurement_schemas.emplace_back(
                std::make_shared<MeasurementSchema>(
                    "s" + to_string(i), TSDataType::INT64, TSEncoding::PLAIN,
                    CompressionType::UNCOMPRESSED));
            column_categories.emplace_back(ColumnCategory::MEASUREMENT);
        }
        return std::make_shared<TableSchema>("testTable" + to_string(table_num),
                                             measurement_schemas,
                                             column_categories);
    }

    static Tablet gen_tablet(const std::shared_ptr<TableSchema>& table_schema,
                             int offset, int device_num) {
        Tablet tablet(table_schema->get_table_name(),
                      table_schema->get_measurement_names(),
                      table_schema->get_data_types(),
                      table_schema->get_column_categories());
        tablet.init();

        int num_timestamp_per_device = 10;
        for (int i = 0; i < device_num; i++) {
            for (int l = 0; l < num_timestamp_per_device; l++) {
                int row_index = i * num_timestamp_per_device + l;
                tablet.add_timestamp(row_index, offset + l);
                auto column_schemas = table_schema->get_measurement_schemas();
                for (const auto& column_schema : column_schemas) {
                    switch (column_schema->data_type_) {
                        case TSDataType::INT64:
                            tablet.add_value(row_index,
                                             column_schema->measurement_name_,
                                             static_cast<int64_t>(i));
                            break;
                        case TSDataType::TEXT:
                            // TODO: support TEXT
                            tablet.add_value(row_index,
                                             column_schema->measurement_name_,
                                             static_cast<int64_t>(i));
                            break;
                        default:
                            break;
                    }
                }
            }
        }
        tablet.set_row_size(device_num * num_timestamp_per_device);
        return tablet;
    }
};

TEST_F(TsFileWriterTableTest, WriteTableTest) {
    auto table_schema = gen_table_schema(0);
    tsfile_writer_->set_generate_table_schema(true);
    tsfile_writer_->register_table(table_schema);
    auto tablet = gen_tablet(table_schema, 0, 1);
    ASSERT_EQ(tsfile_writer_->write_table(tablet), common::E_OK);
    ASSERT_EQ(tsfile_writer_->flush(), common::E_OK);
    ASSERT_EQ(tsfile_writer_->close(), common::E_OK);
}