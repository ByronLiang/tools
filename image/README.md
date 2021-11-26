# Image Exif

实现对JPEG文件在不进行编码下, 保留图片旋转的 EXIF 参数，但会有损 EXIF 的数据字节规范

## 单独获取指定 EXIF IFD 标签数据

`GetDefineTag` 方法能获取指定标签的数据; 当不符合指定的标签ID, 会偏移至下一个标签字节位置, 避免解析无关标签数据

### 图片编码应用

痛点: Go SDK 编码图片, 会丢失图片的EXIF 数据, 会引发图片旋转角度与原图不一致

解决: 获取图片旋转值后, 在对图片编码时, 使用 (imaging)[http://github.com/disintegration/imaging] 的 Rotate 设置图片旋转参数

## Info

Inspired by [goexif package](https://github.com/rwcarlsen/goexif) [imaging package](https://github.com/disintegration/imaging)

[EXIF格式文档](http://gvsoft.no-ip.org/exif/exif-explanation.html)
