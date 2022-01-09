# 位运算配置

## BitSet

### 说明

适用于开关类型的配置，指定位数代表指定配置项，位数长度则代表可配置项的数量；(int8, int32, int64)

### 特点

1. 只需存储一个整型位数值，能存储多个配置项

2. 调用`AddBit` `DelBit` 对指定配置项进行更新操作

3. `Exist`方法能解析位数值对应的配置项

## BitMap

### 说明

1. 对整型数字进行转换二进制, 存放位图里, 减小内存占用, 同时能快速检测当前数值是否存在，是否重复。可应用于优惠码校验等

2. 默认是按照64字节进行映射, 使用哈希值存储位图, 哈希的key=数值 / 64 ; 哈希的值是 数值 % 64 

3. 大量节省内存, 对于大数据数值存储，能减少64倍原内存的占用空间

## BitCode

1. 通过配置不同groupKey，隔离不同业务生成随机码的规则，可以配置生成随机码的元素占用位数与位置顺序, 元素位数与顺序配置可持久化到数据库

2. 对指定元素传入参数进行校验，避免超出二进制位数

3. 解析随机码的元素数据

### 优化

固定保留两个位数存放版本号，若使用两个位数存放版本号，最多可支持四个版本

可根据版本号来解析随机码。当变更随机码生成规则，有效避免历史生成的随机码无法被解析。每次变更随机码生成规则，版本号都需递增1，

### 例子

订单号由店铺ID, 国家编号与时间戳组成

- 时间戳分配39位 满足秒数级别时间戳
- 店铺ID 分配10位，满足至少1000间店铺(1-1023)
- 国家编号分配4位，满足10个国家

```go
timestampEle := NewBitCodeElement("timestamp", 39)
shopNoEle := NewBitCodeElement("shopNo", 10)
countryNoEle := NewBitCodeElement("countryNo", 4)
err := bcg.Add("orderNo", timestampEle, shopNoEle, countryNoEle)
```

`orderNo` 为业务分组，指定业务分组，并传入相应元素的时间戳、店铺ID与国家编号数据，并生成随机码

```go
data := map[string]uint64{
    "shopNo":    uint64(12),
    "countryNo": uint64(5),
    "timestamp": uint64(time.Now().Unix()),
}
code, err := bcg.GenerateCode("orderNo", data)
```

对随机码进行解析, 根据指定分组配置，可独立解析指定元素数据值

```go
code := uint64(662289952456835072)
parseMap, err := bcg.Parse("orderNo", code, "countryNo", "shopNo")
```
