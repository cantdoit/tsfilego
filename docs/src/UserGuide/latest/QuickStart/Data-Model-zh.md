<!--

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at
    
        http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.

-->
# 数据模型

## 基础概念

为管理工业物联网时序数据，TsFile 的测点数据模型包含如下信息

- TAG：标识列
  - name（String）：列名
  - dataType（TSDataType）：数据类型
- FIELD：物理量列
  - name（String）：列名
  - dataType（TSDataType）：数据类型

<table>       
  <tr>             
    <th rowspan="1">概念</th>             
    <th rowspan="1">定义</th>                          
  </tr>       
  <tr>             
    <th rowspan="1">表</th>
      <th>一类具有相同模式的设备的集合。建模时定义的存储表由标识列、时间列和物理量列三部分组成。</th>    
  </tr>  
  <tr>
    <th rowspan="1">设备标识列</th>
  	<th>设备唯一标识，一个表内可包含0至多个标识列，标识列的值按建表时的列顺序组合形成的复合值称为标识，复合值相同的标识为同一标识。标识列的数据类型目前只能为String，可以不指定，默认为String标识列的值可以全为空写入时必须指定所有标识列（未指定的标识列默认使用 null 填充）</th>
  </tr>
  <tr>
    <th rowspan="1">时间列</th>  
    <th>一个表必须有一列时间列，相同标识取值的数据默认按时间排序。时间列的值不能为空，必须顺序的。</th>
  </tr> 
  <tr>             
    <th rowspan="1">物理量列</th>  
    <th>测点列定义了时序数据的测点名称、数据类型。</th>
  </tr> 
  <tr> 
    <th rowspan="1">行</th>  
    <th>表中的一行数据</th>
  </tr> 
</table>


## 示例

表是一类具有相同模式的设备的集合。如下图所示，是一个工厂设备的建模管理，每个设备的物理量采集都具备一定共性（如都采集温度和湿度物理量、同一设备的物理量同频采集等），因此可以以设备为单位进行管理。

此时通过【地区】-【工厂】-【设备】（下图橙色列，又称设备标识信息）可以唯一确定一个实体设备，设备最终采集的指标为【温度】、【湿度】、【状态】、【到达时间】（下图蓝色列）。

![]()