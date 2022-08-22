# dataserver 
# RESTful HTTP 接口
## objects
- 文件下载
``` 
    GET http://${host}/objects/<file_hash>
```
## temp
- 文件元信息上传，建立临时文件存储
``` 
    POST http://${host}/temp/<file_name>
    Request Header.Size: file_size
    Response Body:${uuid}
```
- 文件信息上传，将请求体（即上传文件）进行保存
``` 
    PATCH http://${host}/temp/<uuid>
    Request Body: binary file
```
- 确认完成文件上传,临时文件转为正式文件(移动文件位置,并删除文件元信息)
``` 
    PUT http://${host}/temp/<uuid>
```
- 对象数据校验未通过,删除缓存区的临时文件
``` 
    DELETE http://${host}/temp/<uuid>
```
# 消息队列接口
## 心跳信息
每隔5s向apiServers交换机发送一次自身的网络地址，告知本服务存在。
- **Exchange**: apiServers
- **Message Body**: ${ListenAddr}
- **frequent**:5s
## 定位信息
监听dataServers交换机的消息，从消息体中得到定位文件的name,响应该文件是否存在。   
若该文件存在本节点中，通过消息队列回复本机地址。
- 监听消息
  - 交换机: dataServer
  - 消息体: ${file_hash}
- 响应消息(若文件存在)
  - 交换机: 直接回复Relay
  - 消息体: ${ListenAddr} 