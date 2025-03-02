# WP2AI

WP2AI可以将您的WordPress文章变成智能知识库，并通过AI智能匹配和解读，使其更准确的回答问题。

### 功能特点

- [x] 扫描WordPress文章
- [x] 向量化WordPress文章数据
- [x] AI问答搜索
- [x] 后台管理
- [x] API接口
- [] WordPress插件（用于自动添加、更新、删除文章索引）
- [] 后台文章管理支持状态分类和重试
- [] 自动翻译文章 


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

目前我们仅支持Docker安装，请确保您已经安装Docker环境。（建议安装在WordPress同一服务器）

使用`docker-compose.yaml`，内容如下：

```yaml
version: '3'
services:
    wp2ai:
        container_name: wp2ai
        volumes:
            - '/opt/wp2ai/data:/opt/wp2ai/data'
        network_mode: "host"
        # ports:
        #     - '2080:2080'
        restart: always
        image: 'pub.tcp.mk/helloz/wp2ai:dev-2025022818'
        restart: always
```

注意上面使用了HOST网络模式，安装完毕后您需要在防火墙或安全组放行`2080`端口。

### 使用

1. 安装完毕后访问 `http://IP:2080` 根据提示完成初始化
2. 在后台【参数设置 - WordPress】，填写WordPress数据库地址等信息
3. 在后台【参数设置 - 向量模型】，填写向量模型API信息
4. 在后台【参数设置 - AI模型】，填写AI模型API信息
5. 在后台【文章数据 - 批量扫描】，将WrdPress扫描入库，期间系统会自动向量化数据（1分钟处理10条）
6. 等等数据处理完毕，在后台【AI检索】进行问答测试

### 其他产品

如果您有兴趣，还可以了解我们的其他产品。

* [Zdir](https://www.zdir.pro/zh/) - 一款轻量级、多功能的文件分享程序。
* [OneNav](https://www.onenav.top/) - 高效的浏览器书签管理工具，将您的书签集中式管理。
* [ImgURL](https://www.imgurl.org/) - 2017年上线的免费图床。

### 联系我们

QQ交流群：`964597848`

![ba55b32730f40f20.jpg](https://img.rss.ink/imgs/2025/02/28/ba55b32730f40f20.jpg)

微信交流群

![0ae4ddce930c5c59.jpg](https://img.rss.ink/imgs/2025/02/28/0ae4ddce930c5c59.jpg)