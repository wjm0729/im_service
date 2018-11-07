对照README文档

1，安装go，百度之。

2，克隆之。 https://github.com/GoBelieveIO/im_service.git

3，编译proto文件，虽然还不知道是啥，不过，编呗。

FQ？不用的，应该说翻了也没用，FQ之后还需要设置代理什么的，貌似ss中的socks5代理不能代理命令行。

只有一个go get google.golang.org/grpc需要FQ去下，在github上找到替代方案，如下：

进入$GOPATH/src/google.golang.org目录下，没有的话新建之，

git clone https://github.com/grpc/grpc-go.git grpc 【注意改名，虽然不知道不改会怎么样，不想试了，心累】

之后跟着文档走：

go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

export PATH=$PATH:$GOPATH/bin

protoc -Irpc/ rpc/rpc.proto --go_out=plugins=grpc:rpc

此时，提示protoc没有该命令，到$GOPATH/bin/下查看，果然没有，google之。

参考链接：http://lihaoquan.me/2017/6/29/how-to-use-protobuf.html

到https://github.com/google/protobuf/releases下面下载对应版本，我下的linux x86_64.zip，

解压到$GOPATH/bin/，将解压后的bin下的protoc拷贝到$GOPATH/bin/，其他的就可以过河拆桥了。

【小插曲：wget该链接居然无法解析主机，将wifi换成外网就好了，猜测是又有人改动公司的内网了，心塞】

执行protoc -Irpc/ rpc/rpc.proto --go_out=plugins=grpc:rpc，没有报错即是成功。

python -m grpc.tools.protoc -Irpc --python_out=rpc/ --grpc_python_out=rpc/ rpc/rpc.proto

报错：/usr/bin/python: No module named grpc.tools

没有gtpc.tools模块，安装之，

sudo pip install grpcio.tools【注意1：sudo，注意2：grpcio.tools，io别少】

执行python -m grpc.tools.protoc -Irpc --python_out=rpc/ --grpc_python_out=rpc/ rpc/rpc.proto，成功。

4，编译

cd im_service

mkdir bin

go get github.com/bitly/go-simplejson

go get github.com/golang/glog

go get github.com/go-sql-driver/mysql

go get github.com/garyburd/redigo/redis

go get github.com/googollee/go-engine.io

go get github.com/richmonkey/cfg

go get github.com/valyala/gorpc

一切顺利

make install

在bin下生成im ims imr三个可执行文件，【ok】

之后的过程，需要安装mysql以及redis，此处略。