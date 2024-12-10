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

#include "parser/PathParser.h"
#include "parser/PathLexer.h"
#include "common/path.h"

namespace storage {

class PathNameTest : public ::testing::Test {};

TEST_F(PathNameTest, TestPathLexer) {
    antlr4::ANTLRInputStream input(std::string("root.sg1.'.d1'.s1"));
    PathLexer lexer(&input);
    antlr4::CommonTokenStream tokens(&lexer);
    tokens.fill();
    std::vector<std::string> actualTokens;
    for (const auto& token : tokens.getTokens()) {
        if (token->getType() != antlr4::Token::EOF) {
            actualTokens.push_back(token->getText());
        }
    }
    std::vector<std::string> expectedTokens = {"root", ".", "sg1", ".", "'.d1'", ".", "s1"};
    EXPECT_EQ(actualTokens, expectedTokens);
}

TEST_F(PathNameTest, TestPathNodeGenerator) {
    const std::string path_str = "root.sg1.`.d1`.s1";
    Path path(path_str, true);
    EXPECT_EQ(path.device_, "root.sg1.`.d1`");
    EXPECT_EQ(path.measurement_, "s1");
}
 

}  // namespace storage
