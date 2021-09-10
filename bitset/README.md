# 位运算配置

## 说明

适用于开关类型的配置，指定位数代表指定配置项，位数长度则代表可配置项的数量；(int8, int32, int64)

## 特点

1. 只需存储一个整型位数值，能存储多个配置项

2. 调用`AddBit` `DelBit` 对指定配置项进行更新操作

3. `Exist`方法能解析位数值对应的配置项