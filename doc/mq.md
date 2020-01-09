## rabbitMQ

安装

[https://www.cnblogs.com/vaiyanzi/p/9531607.html](https://www.cnblogs.com/vaiyanzi/p/9531607.html)

[https://www.cnblogs.com/JustinLau/p/11738511.html](https://www.cnblogs.com/JustinLau/p/11738511.html)

使用

[https://juejin.im/post/5dd8cd7ae51d4523501f7331](https://juejin.im/post/5dd8cd7ae51d4523501f7331)

```bash
# 进入sbin
λ cd C:\Program Files\RabbitMQ Server\rabbitmq_server-3.8.1\sbin

# 启用插件
λ rabbitmq-plugins enable rabbitmq_management

# 查看user_list
λ rabbitmqctl list_users
Listing users ...
user    tags
guest   [administrator]

# 添加用户
λ rabbitmqctl add_user rabbitAdmin rabbitAdmin

# 设置用户权限
λ rabbitmqctl set_user_tags rabbitAdmin administrator
Setting tags for user "rabbitAdmin" to [administrator] ...
```

访问 http://localhost:15672/


https://www.jianshu.com/p/1f9915818992

[goland license](https://blog.csdn.net/xinghuo0007/article/details/99899751)
