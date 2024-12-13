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
        tsfile_writer_ = new TsFileWriter();
        libtsfile_init();
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
    TsFileWriter *tsfile_writer_ = nullptr;

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

    static std::string field_to_string(storage::Field *value) {
        if (value->type_ == common::TEXT) {
            return std::string(value->value_.sval_);
        } else {
            std::stringstream ss;
            switch (value->type_) {
                case common::BOOLEAN:
                    ss << (value->value_.bval_ ? "true" : "false");
                    break;
                case common::INT32:
                    ss << value->value_.ival_;
                    break;
                case common::INT64:
                    ss << value->value_.lval_;
                    break;
                case common::FLOAT:
                    ss << value->value_.fval_;
                    break;
                case common::DOUBLE:
                    ss << value->value_.dval_;
                    break;
                case common::NULL_TYPE:
                    ss << "NULL";
                    break;
                default:
                    ASSERT(false);
                    break;
            }
            return ss.str();
        }
    }

    static std::shared_ptr<TableSchema> gen_table_schema(int table_num) {
        std::vector<std::shared_ptr<MeasurementSchema>> measurement_schemas;
        std::vector<ColumnCategory> column_categories;
        int id_schema_num = 5;
        int measurement_schema_num = 5;
        for (int i = 0; i < id_schema_num; i++) {
            measurement_schemas.emplace_back(new MeasurementSchema(
                "id" + to_string(i), TSDataType::TEXT, TSEncoding::PLAIN,
                CompressionType::UNCOMPRESSED));
            column_categories.emplace_back(ColumnCategory::ID);
        }
        for (int i = 0; i < measurement_schema_num; i++) {
            measurement_schemas.emplace_back(new MeasurementSchema(
                "s" + to_string(i), TSDataType::INT64, TSEncoding::PLAIN,
                CompressionType::UNCOMPRESSED));
            column_categories.emplace_back(ColumnCategory::MEASUREMENT);
        }
        return std::make_shared<TableSchema>("testTable" + to_string(table_num),
                                             measurement_schemas,
                                             column_categories);
    }

    Tablet gen_tablet(std::shared_ptr<TableSchema> table_schema, int offset,
                      int device_num) {
        Tablet tablet(table_schema->get_table_name(),
                      table_schema->get_measurement_names(),
                      table_schema->get_data_types(),
                      table_schema->get_column_categories());
        tablet.init();

        return tablet;
    }

    //   public static Tablet genTablet(TableSchema tableSchema, int offset, int
    //   deviceNum) {
    //        Tablet tablet =
    //            new Tablet(
    //                tableSchema.getTableName(),
    //                IMeasurementSchema.getMeasurementNameList(tableSchema.getColumnSchemas()),
    //                IMeasurementSchema.getDataTypeList(tableSchema.getColumnSchemas()),
    //                tableSchema.getColumnTypes());
    //
    //        for (int i = 0; i < deviceNum; i++) {
    //            for (int l = 0; l < numTimestampPerDevice; l++) {
    //                int rowIndex = i * numTimestampPerDevice + l;
    //                tablet.addTimestamp(rowIndex, offset + l);
    //                List<IMeasurementSchema> columnSchemas =
    //                tableSchema.getColumnSchemas(); for (int j = 0; j <
    //                columnSchemas.size(); j++) {
    //                    IMeasurementSchema columnSchema =
    //                    columnSchemas.get(j); tablet.addValue(
    //                        columnSchema.getMeasurementName(), rowIndex,
    //                        getValue(columnSchema.getType(), i));
    //                }
    //            }
    //        }
    //        tablet.setRowSize(deviceNum * numTimestampPerDevice);
    //        return tablet;
    //    }
};

TEST_F(TsFileWriterTableTest, InitWithNullWriteFile) {
    auto table_schema = std::make_shared<TableSchema>();
    tsfile_writer_->set_generate_table_schema(true);
    tsfile_writer_->register_table(gen_table_schema(0));
}