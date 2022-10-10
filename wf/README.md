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

当被监听文件长时间未产生事件(约三小时), inotify 或 epoll 可能无法正确监听文件行为。主要呈现休眠状态

有必要进行心跳检测，保持监测活性

### k8s 挂载文件路径与非挂载路径心跳监测差异

### inotify 监听数量控制

避免引发系统内核异常: inotify watch limit reached

配置路径: `/proc/sys/fs/inotify`

max_user_watches: 可监听文件数




