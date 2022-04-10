版块
------

**forum(版块) 及 bbs_forum_access(版块规则) 对应表关系**

```sql
# 板块表，一级, runtime 中存放 forumlist 格式化以后的数据。
DROP TABLE IF EXISTS bbs_forum;
CREATE TABLE bbs_forum (				
  fid int(11) unsigned NOT NULL auto_increment,		# fid
 # fup int(11) unsigned NOT NULL auto_increment,	# 上一级版块，二级版块作为插件
  name char(16) NOT NULL default '',			# 版块名称
  rank tinyint(3) unsigned NOT NULL default '0',	# 显示，倒序，数字越大越靠前
  threads mediumint(8) unsigned NOT NULL default '0',	# 主题数
  todayposts mediumint(8) unsigned NOT NULL default '0',# 今日发帖，计划任务每日凌晨０点清空为０，
  todaythreads mediumint(8) unsigned NOT NULL default '0',# 今日发主题，计划任务每日凌晨０点清空为０
  brief text NOT NULL,					# 版块简介 允许HTML
  announcement text NOT NULL,				# 版块公告 允许HTML
  accesson int(11) unsigned NOT NULL default '0',	# 是否开启权限控制
  orderby tinyint(11) NOT NULL default '0',		# 默认列表排序，0: 顶贴时间 last_date， 1: 发帖时间 tid
  create_date int(11) unsigned NOT NULL default '0',	# 板块创建时间
  icon int(11) unsigned NOT NULL default '0',		# 板块是否有 icon，存放最后更新时间
  moduids char(120) NOT NULL default '',		# 每个版块有多个版主，最多10个： 10*12 = 120，删除用户的时候，如果是版主，则调整后再删除。逗号分隔
  seo_title char(64) NOT NULL default '',		# SEO 标题，如果设置会代替版块名称
  seo_keywords char(64) NOT NULL default '',		# SEO keyword
  PRIMARY KEY (fid)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
INSERT INTO bbs_forum SET fid='1', name='默认版块', brief='默认版块介绍';
#  cache_date int(11) NOT NULL default '0',		# 最后 threadlist 缓存的时间，6种排序前10页结果缓存。如果是前10页，先读缓存，并依据此字段过期。更新条件：发贴
```

```sql
# 版块访问规则, forum.accesson 开启时生效, 记录行数： fid * gid
DROP TABLE IF EXISTS bbs_forum_access;
CREATE TABLE bbs_forum_access (				# 字段中文名
  fid int(11) unsigned NOT NULL default '0',		# fid
  gid int(11) unsigned NOT NULL default '0',		# fid
  allowread tinyint(1) unsigned NOT NULL default '0',	# 允许查看
  allowthread tinyint(1) unsigned NOT NULL default '0',	# 允许发主题
  allowpost tinyint(1) unsigned NOT NULL default '0',	# 允许回复
  allowattach tinyint(1) unsigned NOT NULL default '0',	# 允许上传附件
  allowdown tinyint(1) unsigned NOT NULL default '0',	# 允许下载附件
  PRIMARY KEY (fid, gid)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

### XiunoBBS
```
MySQL [xn]> desc bbs_forum;
+--------------+-----------------------+------+-----+---------+----------------+
| Field        | Type                  | Null | Key | Default | Extra          |
+--------------+-----------------------+------+-----+---------+----------------+
| fid          | int(11) unsigned      | NO   | PRI | NULL    | auto_increment |
| name         | char(16)              | NO   |     |         |                |
| rank         | tinyint(3) unsigned   | NO   |     | 0       |                |
| threads      | mediumint(8) unsigned | NO   |     | 0       |                |
| todayposts   | mediumint(8) unsigned | NO   |     | 0       |                |
| todaythreads | mediumint(8) unsigned | NO   |     | 0       |                |
| brief        | text                  | NO   |     | NULL    |                |
| announcement | text                  | NO   |     | NULL    |                |
| accesson     | int(11) unsigned      | NO   |     | 0       |                |
| orderby      | tinyint(11)           | NO   |     | 0       |                |
| create_date  | int(11) unsigned      | NO   |     | 0       |                |
| icon         | int(11) unsigned      | NO   |     | 0       |                |
| moduids      | char(120)             | NO   |     |         |                |
| seo_title    | char(64)              | NO   |     |         |                |
| seo_keywords | char(64)              | NO   |     |         |                |
+--------------+-----------------------+------+-----+---------+----------------+
15 rows in set (0.00 sec)

MySQL [xn]> desc bbs_forum_access;
+-------------+---------------------+------+-----+---------+-------+
| Field       | Type                | Null | Key | Default | Extra |
+-------------+---------------------+------+-----+---------+-------+
| fid         | int(11) unsigned    | NO   | PRI | 0       |       |
| gid         | int(11) unsigned    | NO   | PRI | 0       |       |
| allowread   | tinyint(1) unsigned | NO   |     | 0       |       |
| allowthread | tinyint(1) unsigned | NO   |     | 0       |       |
| allowpost   | tinyint(1) unsigned | NO   |     | 0       |       |
| allowattach | tinyint(1) unsigned | NO   |     | 0       |       |
| allowdown   | tinyint(1) unsigned | NO   |     | 0       |       |
+-------------+---------------------+------+-----+---------+-------+
7 rows in set (0.01 sec)
```

### Discuz
```
MySQL [dx]> desc pre_forum_forum;
+------------------+-----------------------------+------+-----+---------+----------------+
| Field            | Type                        | Null | Key | Default | Extra          |
+------------------+-----------------------------+------+-----+---------+----------------+
| fid              | mediumint(8) unsigned       | NO   | 版块ID                          |
| fup              | mediumint(8) unsigned       | NO   | 上级版块ID                       |
| type             | enum('group','forum','sub') | NO   | 版块类型                         |
| name             | char(50)                    | NO   | 名称                            |
| status           | tinyint(1)                  | NO   | 显示状态                         |
| displayorder     | smallint(6)                 | NO   | 显示顺序                         |
| styleid          | smallint(6) unsigned        | NO   | 风格ID                          |
| threads          | mediumint(8) unsigned       | NO   | 主题数量                         |
| posts            | mediumint(8) unsigned       | NO   | 帖子数量                         |
| todayposts       | mediumint(8) unsigned       | NO   | 今日发帖数量                     |
| yesterdayposts   | mediumint(8) unsigned       | NO   | 昨日发帖数量                     |
| rank             | smallint(6) unsigned        | NO   | 版块主题排序(按最新帖/最后回复排序) |
| oldrank          | smallint(6) unsigned        | NO   |     | 0       |                |
| lastpost         | char(110)                   | NO   | 最后发帖 (主题ID 标题 时间戳 用户名 |
| domain           | char(15)                    | NO   | 绑定的二级域名                   |
| allowsmilies     | tinyint(1)                  | NO   |     | 0       |                |
| allowhtml        | tinyint(1)                  | NO   |     | 0       |                |
| allowbbcode      | tinyint(1)                  | NO   |     | 0       |                |
| allowimgcode     | tinyint(1)                  | NO   |     | 0       |                |
| allowmediacode   | tinyint(1)                  | NO   |     | 0       |                |
| allowanonymous   | tinyint(1)                  | NO   |     | 0       |                |
| allowpostspecial | smallint(6) unsigned        | NO   |     | 0       |                |
| allowspecialonly | tinyint(1) unsigned         | NO   |     | 0       |                |
| allowappend      | tinyint(1) unsigned         | NO   |     | 0       |                |
| alloweditrules   | tinyint(1)                  | NO   |     | 0       |                |
| allowfeed        | tinyint(1)                  | NO   |     | 1       |                |
| allowside        | tinyint(1)                  | NO   |     | 0       |                |
| recyclebin       | tinyint(1)                  | NO   |     | 0       |                |
| modnewposts      | tinyint(1)                  | NO   |     | 0       |                |
| jammer           | tinyint(1)                  | NO   |     | 0       |                |
| disablewatermark | tinyint(1)                  | NO   |     | 0       |                |
| inheritedmod     | tinyint(1)                  | NO   |     | 0       |                |
| autoclose        | smallint(6)                 | NO   |     | 0       |                |
| forumcolumns     | tinyint(3) unsigned         | NO   |     | 0       |                |
| catforumcolumns  | tinyint(3) unsigned         | NO   |     | 0       |                |
| threadcaches     | tinyint(1)                  | NO   |     | 0       |                |
| alloweditpost    | tinyint(1) unsigned         | NO   |     | 1       |                |
| simple           | tinyint(1) unsigned         | NO   |     | 0       |                |
| modworks         | tinyint(1) unsigned         | NO   |     | 0       |                |
| allowglobalstick | tinyint(1)                  | NO   |     | 1       |                |
| level            | smallint(6)                 | NO   |     | 0       |                |
| commoncredits    | int(10) unsigned            | NO   |     | 0       |                |
| archive          | tinyint(1)                  | NO   |     | 0       |                |
| recommend        | smallint(6) unsigned        | NO   |     | 0       |                |
| favtimes         | mediumint(8) unsigned       | NO   |     | 0       |                |
| sharetimes       | mediumint(8) unsigned       | NO   |     | 0       |                |
| disablethumb     | tinyint(1)                  | NO   |     | 0       |                |
| disablecollect   | tinyint(1)                  | NO   |     | 0       |                |
+------------------+-----------------------------+------+-----+---------+----------------+
48 rows in set (0.00 sec)

MySQL [dx]> desc pre_forum_forumfield;
+------------------+-----------------------+------+-----+---------+-------+
| Field            | Type                  | Null | Key | Default | Extra |
+------------------+-----------------------+------+-----+---------+-------+
| fid              | mediumint(8) unsigned | NO   | PRI | 0       |版块ID  |
| description      | text                  | NO   |     | NULL    |       |
| password         | varchar(12)           | NO   |     |         |       |
| icon             | varchar(255)          | NO   |     |         |       |
| redirect         | varchar(255)          | NO   |     |         |       |
| attachextensions | varchar(255)          | NO   |     |         |       |
| creditspolicy    | mediumtext            | NO   |     | NULL    |       |
| formulaperm      | text                  | NO   |     | NULL    |       |
| moderators       | text                  | NO   |     | NULL    |       |
| rules            | text                  | NO   |     | NULL    |       |
| threadtypes      | text                  | NO   |     | NULL    |       |
| threadsorts      | text                  | NO   |     | NULL    |       |
| viewperm         | text                  | NO   |     | NULL    |       |
| postperm         | text                  | NO   |     | NULL    |       |
| replyperm        | text                  | NO   |     | NULL    |       |
| getattachperm    | text                  | NO   |     | NULL    |       |
| postattachperm   | text                  | NO   |     | NULL    |       |
| postimageperm    | text                  | NO   |     | NULL    |       |
| spviewperm       | text                  | NO   |     | NULL    |       |
| seotitle         | text                  | NO   |     | NULL    |       |
| keywords         | text                  | NO   |     | NULL    |       |
| seodescription   | text                  | NO   |     | NULL    |       |
| supe_pushsetting | text                  | NO   |     | NULL    |       |
| modrecommend     | text                  | NO   |     | NULL    |       |
| threadplugin     | text                  | NO   |     | NULL    |       |
| replybg          | text                  | NO   |     | NULL    |       |
| extra            | text                  | NO   |     | NULL    |       |
| jointype         | tinyint(1)            | NO   |     | 0       |       |
| gviewperm        | tinyint(1) unsigned   | NO   |     | 0       |       |
| membernum        | smallint(6) unsigned  | NO   | MUL | 0       |       |
| dateline         | int(10) unsigned      | NO   | MUL | 0       |       |
| lastupdate       | int(10) unsigned      | NO   | MUL | 0       |       |
| activity         | int(10) unsigned      | NO   | MUL | 0       |       |
| founderuid       | mediumint(8) unsigned | NO   |     | 0       |       |
| foundername      | varchar(255)          | NO   |     |         |       |
| banner           | varchar(255)          | NO   |     |         |       |
| groupnum         | smallint(6) unsigned  | NO   |     | 0       |       |
| commentitem      | text                  | NO   |     | NULL    |       |
| relatedgroup     | text                  | NO   |     | NULL    |       |
| picstyle         | tinyint(1)            | NO   |     | 0       |       |
| widthauto        | tinyint(1)            | NO   |     | 0       |       |
| noantitheft      | tinyint(1)            | NO   |     | 0       |       |
| noforumhidewater | tinyint(1)            | NO   |     | 0       |       |
| noforumrecommend | tinyint(1)            | NO   |     | 0       |       |
| livetid          | mediumint(8) unsigned | NO   |     | 0       |       |
| price            | mediumint(8) unsigned | NO   |     | 0       |       |
+------------------+-----------------------+------+-----+---------+-------+
46 rows in set (0.00 sec)
```

### 对应关系 - forum 版块
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| fid       | fid                 |  版块ID
| name      | name                |  版块名称
| rank      | rank                |  版块排序
| threads   | threads             |  主题数
| todayposts| todayposts          |  今日发帖数
| todaythreads | -                          |  今日主题数
| brief        | - 对应forumfield.description|  版块介绍
| announcement | - 对应forumfield.rules      |  版块公告
| access       | -                          |  是否开启权限控制
| orderby      | -                          |  默认列表排序
| create_date  | -                          |  版块创建时间
| icon         | - 对应forumfield.icon<3>    |  图标最后存放时间(0为无图标)
| moduids      | -<6>                       |  版主id(1,2,3)
| seo_title    | - 对应forumfield.seotitle   |  SEO标题
| seo_keywords | - 对应forumfield.keywords   |  SEO关键词
+-----------+---------------------+------+-----+---------+----------------+
```

### 对应关系 - bbs_forum_access 版块规则(当forum.access开启时启用)
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| fid       | fid                 |  版块ID
+-----------+---------------------+------+-----+---------+----------------+
```

## 备注
- 使用到两表: pre_forum_forum, pre_forum_forumfield
- pre_forum_forum.type = 'forum', status == 1 时值才对应
- icon - 如果存在 icon 则值为时间戳
- data/attachment/common/a5/common_{$fid}_icon.png - upload/forum/{$fid}.png 
- data/attachment/common/a5/common_{$fid}_icon.png - data/attachment/common/
- moduids 版主清空 - 对应pre_forum_forumfield.moderators (用户名1\n用户名2) (扩展处修复)