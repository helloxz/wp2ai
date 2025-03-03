# WP2AI

WP2AI可以将您的WordPress文章变成智能知识库，并通过AI智能匹配和解读，使其更准确的回答问题。

> 演示地址：[https://ai.xiaoz.top/](https://ai.xiaoz.top/)

### 功能特点

- [x] 扫描WordPress文章
- [x] 向量化WordPress文章数据
- [x] AI问答搜索
- [x] 后台管理
- [x] API接口
- [ ] WordPress插件（用于自动添加、更新、删除文章索引）
- [ ] 后台文章管理支持状态分类和重试
- [ ] 自动翻译文章 
- [ ] 记录用户问答记录


### 适用场景

* 需要提升站内搜索效率的博客网站
* 拥有大量专业知识内容的在线平台
* 希望为用户提供更优质搜索服务的企业网站

### 部分截图

![0771f4d41f551d81.png](https://img.rss.ink/imgs/2025/03/02/0771f4d41f551d81.png)

![ba02206bb4893bce.png](https://img.rss.ink/imgs/2025/03/02/ba02206bb4893bce.png)

![0127483b258c677c.png](https://img.rss.ink/imgs/2025/03/02/0127483b258c677c.png)

![dfd708b02006a65b.png](https://img.rss.ink/imgs/2025/03/02/dfd708b02006a65b.png)

### 安装

目前仅支持Docker安装，请确保您已经安装Docker环境。

**方式一：Docker Compose安装**

`docker-compose.yaml`内容如下：

```yaml
version: '3'
services:
    wp2ai:
        container_name: wp2ai
        volumes:
            - '/opt/wp2ai/data:/opt/wp2ai/data'
        network_mode: "host"
        restart: always
        image: 'helloz/wp2ai'
```

注意上面使用了HOST网络模式，安装完毕后您需要在防火墙或安全组放行`2080`端口。如果WP2AI和WordPress不在同一服务器，也可以使用`bridge`网络模式，只需要注释掉`network_mode`并自定义端口，比如：

```yaml
version: '3'
services:
  wp2ai:
    container_name: wp2ai
    volumes:
      - '/opt/wp2ai/data:/opt/wp2ai/data'
    ports:
      - '2080:2080'
    image: 'helloz/wp2ai'
    restart: always
```

**方式二：Docker命令行安装**

使用HOST网络模式，需要自行放行`2080`端口！

```bash
docker run -d \
  --name wp2ai \
  -v /opt/wp2ai/data:/opt/wp2ai/data \
  --network host \
  --restart always \
  helloz/wp2ai
```

使用`network_mode`模式，可自定义访问端口：

```bash
docker run -d \
  --name wp2ai \
  -v /opt/wp2ai/data:/opt/wp2ai/data \
  -p 2080:2080 \
  --restart always \
  helloz/wp2ai
```



### 使用

1. 安装完毕后访问 `http://IP:2080` 根据提示完成初始化
2. 在后台【参数设置】，填写WordPress数据、向量模型、AI模型等信息
3. 在后台【文章数据 - 批量扫描】，将WrdPress扫描入库，期间系统会自动向量化数据（1分钟大约处理15条）
4. 等等数据处理完毕，在【后台 - AI检索】或者网站首页进行问答测试

### 其他产品

如果您有兴趣，还可以了解我们的其他产品。

* [Zdir](https://www.zdir.pro/zh/) - 一款轻量级、多功能的文件分享程序。
* [OneNav](https://www.onenav.top/) - 高效的浏览器书签管理工具，将您的书签集中式管理。
* [ImgURL](https://www.imgurl.org/) - 2017年上线的免费图床。

### 安装服务

WP2AI目前是免费开源的，但是我们也提供付费安装服务，如果需要我们协助部署，请联系微信：`xiaozme`，收费为`80元/次`

### 联系我们

QQ交流群：`964597848`

![ba55b32730f40f20.jpg](https://img.rss.ink/imgs/2025/02/28/ba55b32730f40f20.jpg)

微信交流群

![0ae4ddce930c5c59.jpg](https://img.rss.ink/imgs/2025/02/28/0ae4ddce930c5c59.jpg)