## frp-notify

[![Build](https://github.com/arugal/frp-notify/workflows/Build/badge.svg?branch=master)](https://github.com/arugal/frp-notify/actions?query=branch%3Amaster+event%3Apush+workflow%3ABuild)

一个专注于消息通知的 [frp server manager plugin](https://github.com/fatedier/frp/blob/master/doc/server_plugin_zh.md) 实现，让你对进入 `frps` 的连接了如指掌，不再裸奔。

目前支持将 `Login`、`NewProxy`、`NewWorkConn` 和 `NewUserConn` 操作通知到 [gotify-server](https://github.com/gotify/server) 。

## 快速启动

[下载地址](https://github.com/arugal/frp-notify/releases)

### 目录介绍

```bash
* frp-notify
└─── system
|    |          frp-notify.service                  # linux 系统服务配置文件
|
│           frp-notify                              # frp-notify 程序
|           notify-plugin.json                      # 通知插件配置文件
```

### 打印帮助信息

```bash
./frp-notify --help
```

### 命令行启动

```bash
./frp-notify start
```

### docker

```bash
docker run -p 50080:80 -v /etc/frp-notify/notify-plugin.json:/notify-plugin.json arugal/frp-notify:latest start
```

## 配置介绍

### frps

在 `frps.ini` 增加以下配置 

```
[plugin.frp-notify]
addr = 127.0.0.1:80                             // frp-notify 地址
path = /handler                                 // frp-notify url, 固定配置
ops = Login,NewProxy,NewWorkConn,NewUserConn    // 通知的操作
```

### 通知插件

#### gotify

```
{
  "notify_plugins": [
    {
      "name": "gotify",                         // 固定配置
      "config": {
        "server_addr": "127.0.0.1:4000",        // gotify-server 服务地址
        "app_token": "token"                    // gotify-server 配置的 app token
      }
    }
  ]
}
```