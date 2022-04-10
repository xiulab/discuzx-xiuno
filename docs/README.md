### Discuz!X 3.x 转换 XiunoBBS 4.x
------
XiunoBBS 主程序含 **17**表

------
### 需要处理的表
- attach - 附件表，存放图片与附件 - [attach.md](attach.md)
- forum - 版块表 - [forum.md](forum.md)
- group - 用户组表 - [group.md](group.md)
- post - 帖子表 - [post.md](post.md)
- mypost - 我的帖子表 - [post.md](post.md#%E5%AF%B9%E5%BA%94%E5%85%B3%E7%B3%BB---mypost-%E6%88%91%E7%9A%84%E5%B8%96%E5%AD%90)
- thread - 主题表 - [thread.md](thread.md)
- mythread - 我的主题表 - [thread.md](thread.md)
- thread_top - 置顶主题表 - [thread.md](thread.md)
- user - 用户表 - [user.md](user.md)

- extension - 最终修正附件、头像、版块 icon、管理员等功能 - [extension.md](extension.md)

### 不作处理的表
- cache - 缓存表，用来保存临时数据
- forum_access 版块访问规则
- kv - 持久的 key value 数据存储
- modlog -  版主操作日志
- queue - 临时队列，用来保存临时数据
- session - session 表
- session_data 
- table_day - 系统表

**Xiuno 表解析参考:** https://www.kancloud.cn/xiuno/bbs4/214347   
**Discuz 表解析参考:** http://discuzt.cr180.com/discuzcode-db.html   

------
### 作者
Author: Skiychan   
Email : dev@skiy.net   
Link  : https://www.skiy.net      