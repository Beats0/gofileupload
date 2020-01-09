## gofileupload

一个用 vue + go 写的网盘服务器

线上体验

[demo](http://47.94.16.206:8081)

![go file upload.png](https://wx1.sinaimg.cn/large/006nOlwNly1gaid15uytij31hc0pe43f.jpg)

*Tip: 双击文件图标可查看文件详情

前端

 - node >= 10
 - vue-cli 3
 - vuetify


后端

 - go >= 1.13.1
 - mysql, Gorm
 - go-gin

TODO:

- [x] 多用户使用
- [x] 单个，多个文件上传
- [x] 大文件分片、断点续传
- [x] 文件MD5检验
- [x] restful api
- [x] token认证
- [x] 文件分类
- [x] 文件搜索
- [x] 文件预览
- [x] 文件下载
- [x] 文件命名
- [x] 视频播放
- [x] 图片缩略图
- [ ] 文件合成下载
- [ ] 优化上传进度控制面板
- [ ] 文件复制、移动
- [ ] 上传文件夹
- [ ] 文件分享
- [ ] 回收站

### 使用

前端

```sh
cd client/gofileupload

### install
yarn install

### dev
yarn run serve

### build
yarn run build
```


默认访问

> http://localhost:8080

后端

```sh
cd

### install
go get -u

### 创建static下assets和upload文件夹
mkdir -p static/assets static/upload

### run debug
go run main.go -mode debug

### build
go build main.go

### run release
go run main.go -mode release
```

默认访问

> http://localhost:8081

### 部署

前端vue打包后

将`client\gofileupload\dist`内所有文件放在 `static\assets`中
 
将`client\gofileupload\dist\index.html`复制到 `templates\index.html` 中

后端添加 `conf\app_release.ini` 配置文件, 内容为自己的服务器部署配置信息

后端打包后，执行`go run main.go -mode release`


配置详见 [config](https://github.com/Beats0/gofileupload/tree/master/conf)
