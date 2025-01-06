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
# Data Model

## Basic Concepts

To manage industrial IoT timing data, the measurement point data model of TsFile includes the following information

- TAG：Tag Column 
  - name（String）：Column Name 
  - dataType（TSDataType）：Data Type
- FIELD：Field Column
    - name（String）：Column Name 
    - dataType（TSDataType）：Data Type

<table>       
  <tr>             
    <th rowspan="1">concept</th>             
    <th rowspan="1">definition</th>                          
  </tr>       
  <tr>             
    <th rowspan="1">table</th>
    <th>A collection of devices with the same pattern.The storage table defined during modeling consists of three parts: identification column, time column, and physical quantity column.</th>    
  </tr>  
  <tr>
    <th rowspan="1">TAG</th>
  	<th>The unique identifier of a device, which can contain 0 to multiple tag columns in a table. The composite value formed by combining the values of the tag columns in the column order when the table was created is called the identifier, and tags with the same composite value are called the same identifier.The data type of the tag column can currently only be String, which can be left unspecified and defaults to StringThe values of the identification column can all be emptyWhen writing, all tag columns must be specified (unspecified identity columns are filled with null by default)</th>
  </tr>
  <tr>
    <th rowspan="1">Time</th>  
    <th>A table must have a time column, and data with the same identifier value is sorted by time by default.The values in the time column cannot be empty and must be in sequence.</th>
  </tr> 
  <tr>             
    <th rowspan="1">FIELD</th>  
    <th>The field column defines the measurement point names and data types for time-series data.</th>
  </tr> 
  <tr> 
    <th rowspan="1">row</th>  
    <th>A row of data in the table</th>
  </tr> 
</table>


## Example

A table is a collection of devices with the same pattern. As shown in the figure below, it is a modeling management of factory equipment, and the physical quantity collection of each device has certain commonalities (such as collecting temperature and humidity physical quantities, collecting physical quantities of the same device on the same frequency, etc.), so it can be managed on a device by device basis.

At this point, a physical device can be uniquely identified through [Region] - [Factory] - [Equipment] (orange column in the figure below, also known as device identification information). The final indicators collected by the device are [Temperature], [Humidity], [Status], and [Arrival Time] (blue column in the figure below).

![]()