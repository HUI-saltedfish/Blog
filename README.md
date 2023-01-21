## [原项目地址](https://github.com/wangsongyan/wblog)
## [博客地址](http://www.huilearn.work/)

## 基于原项目所做的更改
### 1. go vender包管理机制替换成了 go mod包管理机制
### 2. gorm版本由1.19升级到了1.24，并摆脱了对于原生sql语句的依赖
### 3. 将smms图床替换成了阿里oss对象存储服务
### 4. 将sqlite数据库替换成了docker中的mysql数据库