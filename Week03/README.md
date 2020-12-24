1.基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。
#分析
- http server 与 信号的注册和处理 都会阻塞，使用group.Go启动goroutine处理

- http server 启动一个done的chan做监听处理，开放close接口做关闭http server 的控制
- linux signal 信号本身是chan阻塞，有信号就可以解除
- 测试中，errgroup.WithContext 返回的 context.Context 如果不进行 case <-ctx.Done() 处理，退出一个还是会进行阻塞，所以在http server和linux signal都做了上下文取消的监听处理，目前测试是正常的，不知道是不是我的代码有误
