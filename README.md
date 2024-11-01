# 北京无线电考试预约工具

一键预约北京市业余无线电考试（注意本工具只会预约可预约列表的第一个考试，一般来说已足够）

## 用法

首先需要在[北京无线电协会业余无线电服务平台](https://xt.bjwxdxh.org.cn/static/member/#/static/member/user/login) 注册一个账号，并通过审核。

```bash
./main <username> <password>

```

可以通过下面方式查看调试信息：

```bash
DEBUG=1 ./main <username> <password>

```
