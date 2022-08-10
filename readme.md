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
# 拉取镜像
docker pull  elasticsearch:7.17.5
# 创建容器
docker run -d --name es -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.17.5
# 
```
> -d：后台启动  
> --name：容器名称  
> -p：端口映射  
> -e：设置环境变量  
> discovery.type=single-node：单机运行 


``` shell 
# 拉取镜像
docker pull mobz/elasticsearch-head:5
# 创建容器
docker run -d --name elasticsearch-head -p 9100:9100  mobz/elasticsearch-head:5
```
管理web 的端口：http://ip:9100 

跨域问题：
容器内没有vi/vim命令,直接使用echo追加即可
``` shell 
echo http.cors.enabled: true  >> config/elasticsearch.yml
echo http.cors.allow-origin: \"*\"  >> config/elasticsearch.yml
```
参考：  
- https://zhuanlan.zhihu.com/p/257867352