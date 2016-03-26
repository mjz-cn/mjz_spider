# 从公开的proxy网站上采集爬虫代理

**采用go语言编写，用来练手**

## 工程结构
1. handlers  
   在爬取一个新的网站内容时，只需要在相应的handlers模块下添加针对这个网站的handler即可
   以添加一个针对[incloak][incloak]网站的handler为例：
   1. 从[incloak][incloak]网站上抓取代理地址
   2. 调用models提供的方法存入数据库
   3. 将此handler注册到全局的Handlers中
2. models  
   此模块专门负责数据的存储，使用[beego/orm](http://beego.me/docs/mvc/model/overview.md)
3. config  
   project配置信息，使用yaml格式
4. utils  
   提供一些帮助性的方法 

## 爬取的网站有

* [incloak][incloak] （需翻墙）
* [hidemyass][hidemyass]


**To-Do**

1.   Review代码，美化代码
	*   工程结构和变量的重命名
2.   重新梳理工程结构
3.   高效爬取和存储


[incloak]: http://incloak.com/
[hidemyass]: http://proxylist.hidemyass.com
