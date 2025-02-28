# WP2AI

WP2AI 是一套私有部署的智能搜索程序，它利用先进的 AI 技术，让您的 WordPress 网站拥有更强大的搜索能力。

### 核心功能与优势
* **全站内容深度分析：** WP2AI 会深入分析您 WordPress 网站的每一篇文章，理解其内容含义，而非仅仅依赖关键词匹配。
* **语义搜索，快速找到答案：** 用户可以用更自然的方式提问，WP2AI 能理解他们的真实意图，快速找到相关内容。
* **AI 辅助，智能呈现：** WP2AI 不仅能找到文章，还能结合 AI 技术，以更清晰、更易懂的方式呈现搜索结果。
* **私有部署，数据安全可控：** WP2AI 可以部署在您自己的服务器上，确保您的数据安全和隐私。
* **提升搜索体验，增强用户粘性：** WP2AI 为您的网站带来更流畅、更智能的搜索体验，让用户更容易找到所需信息，从而提升用户满意度和网站互动性。

> 内容越垂直，结果越精准！

### 适用场景

* 需要提升站内搜索效率的博客网站
* 拥有大量专业知识内容的在线平台
* 希望为用户提供更优质搜索服务的企业网站

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