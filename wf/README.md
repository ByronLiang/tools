# WatchFile Note

## fsnotify 库原理

调用系统命令(`syscall 方法`) Linux 的 inotify 监听文件，并使用 epoll 绑定 inotify 的信号。 当 epoll 捕捉到信息，可以过滤无用事件, 再将事件传递到应用层，再进行相关业务处理

涉及用户层(应用层)与内核层(kernel)的数据通讯

### inotify 与 epoll 区别

inotify 只能监听文件系统的事件信号。 epoll 可以监听任何信号, 如: Socket, IPC(进程通信)

## 潜在问题

### K8s 的 Docker 容器文件系统

环境： k8s 涉及Docker 容器文件系统机制

使用 Linux 原生系统调用命令，可能不适合 Docker 容器文件系统

当被监听文件长时间未产生事件(约三小时), inotify 或 epoll 可能无法正确监听文件行为。主要呈现休眠状态，丢失 fd

有必要进行心跳检测，保持监测活性

### k8s 挂载文件路径与非挂载路径心跳监测差异

使用 hostpath、EmptyDir 类型的共享卷，会将路径挂载到容器操作系统，能调用 Linux 的文件系统相关接口

### inotify 监听数量控制

避免引发系统内核异常: inotify watch limit reached

配置路径: `/proc/sys/fs/inotify`

max_user_watches: 可监听文件数

#### 查看每个进程的fd inotify 使用数量脚本

```sh
find /proc/*/fd \
-lname anon_inode:inotify \
-printf '%hinfo/%f\n' 2>/dev/null \
\
| xargs grep -c '^inotify'  \
| sort -n -t: -k2 -r
```

## 选型

### Polling 轮询

1. 减少 inotify 数量, 若需监听大量文件，需要考虑扩大 Linux inotify 最大数量

2. 需要调用系统方法: `stat` 获取监听文件信息, 并且对比文件的更改时间, 保存监听文件的基本文件信息

### inotify 与 epoll 组合

1. 保证实时性。应用层只需监听一个被绑定多个 inotify 信号的 epoll 事件信号，即可接收多个 inotify 文件信号。

2. 减小对文件遍历与识别文件是否被更新能行为, 只需对信号做识别
