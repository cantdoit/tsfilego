all things needed to be understood is marked with (whut)





tsfile writing

the order of processing is as follows when writing a File

1. Open a new File in tsfile_writer (making checks for name conficts etc)

   1. calls the create function from write_file
2. Calls tsfile_io_writer to write magic string(whut) and version
3. Starts a chunk group flush process. Each device or logical group of data is written as a "chunk group" (e.g., a collection of measurements for one device)

   1. in tsfile_writer:727 where a chunk group header is written to File (the device)
4. After the ehader is writted it write the data (measurements) into the chunk

   1. defined my Bytesream and ColumnDesc
5. write_file then uses {sync} to write to disk



#### `TsFileWriter`

1. **open:**
   * `WriteFile (create)`
   * → uses `write_file.cc::WriteFile::create`
   * `TsFileIOWriter (init)`
   * → uses `tsfile_io_writer.cc::TsFileIOWriter::init`
2. **start\_file:**
   * `TsFileIOWriter (start_file)`
   * → uses `tsfile_io_writer.cc::TsFileIOWriter::start_file`
3. **start\_flush\_chunk\_group:**
   * `TsFileIOWriter (start_flush_chunk_group)`
   * → uses `tsfile_io_writer.cc::TsFileIOWriter::start_flush_chunk_group`
4. **start\_flush\_chunk:**
   * `TsFileIOWriter (start_flush_chunk)`
   * → uses `tsfile_io_writer.cc::TsFileIOWriter::start_flush_chunk`
5. **sync:**
   * `WriteFile (sync)`
   * → uses `write_file.cc::WriteFile::sync`
6. **close:**
   * `WriteFile (close)`
   * → uses `write_file.cc::WriteFile::close`
