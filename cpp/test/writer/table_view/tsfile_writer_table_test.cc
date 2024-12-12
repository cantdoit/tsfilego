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
#include "writer/tsfile_writer.h"

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
};

TEST_F(TsFileWriterTableTest, InitWithNullWriteFile) {
   TsFileWriter writer;
   ASSERT_EQ(writer.init(nullptr), E_INVALID_ARG);

   auto table_schema = std::make_shared<TableSchema>();

   writer.register_table(nullptr);
}