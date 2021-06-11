# week0607

基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

- 思路是 监听信号量和errGroup返回的ctx的Done()，两种情况：
  - 收到系统退出信号
  - 任意一个 errGroup 的 Go内执行的函数返回error了
