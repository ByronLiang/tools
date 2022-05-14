# 指定位置文件行扫描

通过记录读取每行数据的偏移，能实现直接定位指定位置继续完成剩余文件行读取

## 配置

- 文件读取起始偏移量: `position`

- 文件每行内容读取的数据操作: `ScanHandle`

- 文件读取终止, 对当前文件读取偏移量进行回调操作: `ExitHook`

## 参考

[read file with specific-line-number-using-scanner](https://stackoverflow.com/questions/34654514/how-to-read-a-file-starting-from-a-specific-line-number-using-scanner)
