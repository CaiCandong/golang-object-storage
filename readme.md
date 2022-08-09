# golang-object-storage 
## rabbitqm docker 启动
``` shell 
docker run -d --hostname my-rabbit --name rabbit -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=center -e RABBITMQ_DEFAULT_PASS=123qwe -e RABBITMQ_DEFAULT_VHOST=center  rabbitmq:management
```
> - 账户：center
> - 密码：123qwe
> - v-host: center
> - 15672：控制台端口号
> - 5672：应用访问端口号
> - 管理web 的端口：http://ip:15672

## elasticsearch docker 启动
``` shell 
docker pull  elasticsearch:7.17.5
docker run -d --name es -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.17.5
# 
```
> -d：后台启动  
> --name：容器名称  
> -p：端口映射  
> -e：设置环境变量  
> discovery.type=single-node：单机运行 
> 管理web 的端口：http://ip:9200 

跨域问题：
https://zhuanlan.zhihu.com/p/257867352
https://blog.csdn.net/cecurio/article/details/105578136