### Discuz!X To XiunoBBS

---

基于 `Go` 语言的 `Discuz!X 3.x To XiunoBBS 4.x` 转换工具

### 开发进度

- 基础构架 ✔
- 数据转换 ✔
- 附件迁移 ✔
- 数据优化 ✔

### 编译指南

- 拉取主项目 `git clone https://github.com/xiulab/discuzx-xiuno.git`
- 进入项目目录, 执行 `go get`
- 编译程序 `go build`
- 完成，文件在 `当前目录(go build)` 下
- 修改配置文件信息 `config.toml`

**温馨提示:**

> 如果已配置好`GOBIN`或者将 `$GOPATH/bin` 环境变量，即可以在任何目录下执行 **`discuzx-xiuno`** 启动本程序。  
> 程序必须有**可执行权限**。

### 工具使用教程

- **操作前, 请仔细查阅 config.toml 配置文件**
- 先建一个 `XiunoBBS 4.x` 论坛。
- 下载源码编译的二进制程序（或自行按上述步骤编译）。非 Windows 平台需要可执行权限。
- 配置`confit.toml`, 执行本程序 `./discuzx-xiuno`（Windows 平台, 建议使用 `PowerShell`控制台, 执行`discuzx-xiuno.exe`）
- 登录后台，**记得**更新缓存统计。

### 配置文件说明

> 请认真查阅配置文件的注释，主要修改 database (数据库)、tables.xiuno.user (用户表) 及 extension.file (附件迁移) 这几部分。

<details>
<summary>配置文件内容</summary>

```toml
[setting]

# 日志配置
[log]
    # 日志等级 all.全部日志, prod.一般日志(warning,error), alert.提示日志(warning,error,info), warning.警告日志, info.正常日志, error.错误日志
    level = "alert"
    # 日志保存目录,会在此目录下生成一个当天的日志文件
    path = "logs"
    # 是否输出错误位置,转换出错时建议打开,以便配合作者修复 BUG
    trace = false

# 数据库配置
[database]
    # XiunoBBS
    [[database.xiuno]]
        type = "mysql"      # 数据库类型(不可修改)
        host = "127.0.0.1"  # IP
        port = "3306"       # 端口
        user = "root"       # 数据库用户名
        pass = "123456"     # 密码
        name = "xiuno"      # 数据库名
        prefix = "bbs_"     # 表前缀
        charset = "utf8"    # 字符集
        debug = false     # 日志调试,建议关闭

    # Discuz!X
    [[database.discuz]]
        type = "mysql"
        host = "127.0.0.1"
        port = "3306"
        user = "root"
        pass = "123456"
        name = "discuzx"
        prefix = "pre_"
        charset = "utf8" # 不可改动
        debug = false    # 日志调试,建议关闭

    # UCenter
    [[database.uc]]
        type = "mysql"
        host = "127.0.0.1"
        port = "3306"
        user = "root"
        pass = "123456"
        name = "discuzx"
        prefix = "pre_ucenter_"
        charset = "utf8" # 不可改动
        debug = false    # 日志调试,建议关闭

# 需要转换的表配置
[tables]
    [tables.xiuno]
        # 用户表
        [tables.xiuno.user]
            # 表名
            name = "user"
            # 是否转换
            convert = true
            # 每次更新条数(留空或 < 2, 则默认为 1 条), 当 ucenter 与 discuz!X 不同一个库中 或 multiple_email 值为 2 时, batch 则默认为 1 条, 不作批量导入
            batch = 100
            # 去除 email 的唯一索引(Discuz!X 遗留问题, 若存在多用户用同一个 email 时, 则需要去除索引 或 修改重复的 email)
            # 建议先默认 0, 用工具进去 MySQL 执行 SELECT count(*) c,uid,email FROM `pre_common_member` GROUP BY email ORDER BY `c` DESC
            # 若 c > 1 的数据很多, 则可以设置为 1; 否则, 可以手动将重复的 email 修改掉, 默认 0 即可
            # 0. 正常流程, 1. 去除索引方式, 2. 在重复的 email 前添加 UID_(若 UID 为 555 的用户 email: abc@qq.com 重复, 将变更为 555_abc@qq.com)
            multiple_email = 2

        # 用户组表
        [tables.xiuno.group]
            # 表名
            name = "group"
            # 是否转换
            convert = true
            # 是否使用 xiunobbs 官方用户组
            official = true

        # 版块表
        [tables.xiuno.forum]
            # 表名
            name = "forum"
            # 是否转换
            convert = true

        # 附件表
        [tables.xiuno.attach]
            # 表名
            name = "attach"
            # 是否转换
            convert = true
            # 每次更新条数(留空或 < 2, 则默认为 1 条), 单条导入时, 错误不会导致程序退出
            batch = 1

        # 主题表
        [tables.xiuno.thread]
            # 表名
            name = "thread"
            # 是否转换
            convert = true
            # 每次更新条数(留空或 < 2, 则默认为 1 条; 数据过大时, 建议设置为 1, 否则可能会导致 Killed)
            batch = 100
            # 取 >= TID 的数据。当上次转换出错时, 记录此 TID, 方便再次导入
            last_tid = 0

        # 帖子表
        [tables.xiuno.post]
            # 表名
            name = "post"
            # 是否转换
            convert = true
            # 每次更新条数(留空或 < 2, 则默认为 1 条; 数据过大时, 建议设置为 1, 否则可能会导致 Killed)
            batch = 100
            # 取 >= PID 的数据。当上次转换出错时, 记录此 PID, 方便再次导入
            last_pid = 0

        # 置顶帖子表
        [tables.xiuno.thread_top]
            # 表名
            name = "thread_top"
            # 是否转换
            convert = true

        # 我的主题表
        [tables.xiuno.mythread]
            # 表名
            name = "mythread"
            # 是否转换
            convert = true

        # 我的帖子表
        [tables.xiuno.mypost]
            # 表名
            name = "mypost"
            # 是否转换
            convert = true

# 扩展功能
[extension]
    [extension.forum]
        # 是否导入论坛版主 (不建议使用)
        moderators = false

    [extension.file]
        # 是否启用转移附件文件功能
        enable = false

        # Windows 平台的目录请使用 \\ 或 / 代替 \, 比如 C:\\dist\\abc 或 C:/dist/abc
        # XiunoBBS 论坛绝对路径
        # 若不配置, 则附件、头像及版块 icon 等资源将会复制到当前目录的 files 目录下, 迁移完成后，复制 files 下的 upload 到 XiunoBBS 根目录覆盖即可
        xiuno_path = ""
        # Discuz!X 论坛绝对路径
        discuzx_path = ""

        # 附件转移
        attach = true
        # 头像转移
        avatar = true
        # 版块 ICON 转移
        icon = true

    [extension.group]
        # 管理员 UID
        admin_id = 1

    [extension.user]
        # 是否修正主题 (post.first 全部为0,第一条变更为主题)
        fix_thread = true
        # 是否修正用户主题数和帖子数(帖子数=主题+回复), 非常耗时
        total = true
        # 修正最低等级积分的用户 gid 为 101 的用户组
        normal_user = true

    [extension.thread_post]
        # 是否修正主题的 lastpid 和 lastuid, 比较耗时
        fix_last = true
        # 是否修正帖子内附件统计数量
        post_attach_total = true
        # 是否修正主题内附件统计数量
        thread_attach_total = true

```

</details>

### 使用到的开源项目

- **XiunoBBS**
- **Discuz!X**
- https://github.com/gogf/gf (基础框架)
- https://github.com/frustra/bbcode (内容 BBCODE 转 HTML)
