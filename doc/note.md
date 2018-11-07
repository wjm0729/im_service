[GoBelieve 官网](http://developer.gobelieve.io/)
    
[Gobelieve架构解析](https://www.jianshu.com/p/8121d6e85282)

GoBelieve service部署常见问题总结 [原文链接](https://www.cnblogs.com/nuanshou/p/9071135.html)
    
    问题1：
    大家好，我按照文档的步骤编译im时(make install)，出现 can't load package: package main:
    app_route.go:1:1: expected 'package', found 'ILLEGAL'
    
    答：在linux系统下下载代码编译
    
    
    问题2：有没有配套的测试脚本使用？
    
    答：用vagrant工程里面的python脚本测试，https://github.com/GoBelieveIO/gobelieve_vagrant
    
    
    问题3：客户端的tcp连接超时时间在哪调？
    
    答：没有手工设置超时，使用系统默认的连接超时
    
    
    问题4：是自动重连的吗？间隔是多少   能不能修改？能不能自己设置间隔？
    
    答：自动重连不是固定间隔的，现在没有一个简单的方法设置
     
    
    问题5：im文件夹里面的几个benchmark都是可以单独测的吗？
    
    答：可以的，不过最好先看下代码
     
    
    问题6：请教下，appid和uid分别代表什么标识呀？为什么要两个来表示一台设备
    client.appid client.uid
    
    答：一台设备可能会有两个不同的应用的
     
    
    问题7：
    I0409 19:03:04.426812   20082 im.go:230] dispatch app message:MSG_SYNC_NOTIFY
    I0409 19:03:04.470103   20082 client.go:154] msg cmd:MSG_SYNC
    这个是后端有自己的触发的地方吗？还是完成靠客户端？这二个时间怎么这个近~ 我还以为后端有触发的地方，那算是拉取方式了？
    
    答：客户端，服务器下发新消息通知触发客户端这个动作的
     
    
    问题8：问下 要是自己搭建的话 离线消息推送怎么做呢
    
    答：所有的离线消推送消息都会派送到redis的队列中， 你写一个redis队列的消费者，去调用第三方的推送服务
    
    
    问题9：请问imr横向扩展部署有试过么 可以部署多少个？
    
    答:
    im 客户连接服务器 （可分布式部署，暂无负载均衡模块) 依赖外部负载均衡模块 比如：lvs
    imr 路由查询服务器（主要解决im分布式部署的问题） 可部署多台
    ims 存储服务器 (主从部署) 可部署多组
     
    
    问题10：im，imr，ims三个程序分别对应什么模块功能？
    
    答:
    im服务器由3个服务组成：
    1. im接入服务器 直接负责和终端的连接 -> im
    2. im存储服务器 存储im的离线消息 -> ims
    3. im路由服务器 在接入服务器之间转发消息 -> imr
    ims存储， imr 转发，imr是不同im实例的桥梁
     
    
    问题11：im，imr，ims的启动
    
    答:
    先配置， 然后启动
    im -logtostderr=true im.cfg 
    ims -logtostderr=true ims.cfg 
    imr -logtostderr=true imr.cfg
     
    
    问题12：问一下，StartHttpServer 这个里面的HTTP接口，暴露出来的话任何前端都人员都可以使用了 
    
    答：监听的是内网ip地址，或者localhost。
    
    
    问题13：
    你们在goBelive和前端用户之间还有处理的逻辑层吗？
    goBelive服务器直接提供服务给前端用户吗？
    前端用户之间提交消息给goBelive服务器，还是中间还有中间件进行处理转发？
    
    答：api是要求你自己实现的，没有中间件。