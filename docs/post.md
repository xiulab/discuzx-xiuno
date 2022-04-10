帖子
------

**post(帖子)、mypost(我的帖子) 对应表关系**

```sql
# 论坛帖子数据
DROP TABLE IF EXISTS bbs_post;
CREATE TABLE bbs_post (
  tid int(11) unsigned NOT NULL default '0',		# 主题id
  pid int(11) unsigned NOT NULL auto_increment,		# 帖子id
  uid int(11) unsigned NOT NULL default '0',		# 用户id
  isfirst int(11) unsigned NOT NULL default '0',	# 是否为首帖，与 thread.firstpid 呼应
  create_date int(11) unsigned NOT NULL default '0',	# 发贴时间
  userip int(11) unsigned NOT NULL default '0',		# 发帖时用户ip ip2long()
  images smallint(6) NOT NULL default '0',		# 附件中包含的图片数
  files smallint(6) NOT NULL default '0',		# 附件中包含的文件数
  doctype tinyint(3) NOT NULL default '0',		# 类型，0: html, 1: txt; 2: markdown; 3: ubb
  quotepid int(11) NOT NULL default '0',		# 引用哪个 pid，可能不存在
  message longtext NOT NULL,				# 内容，用户提示的原始数据
  message_fmt longtext NOT NULL,			# 内容，存放的过滤后的html内容，可以定期清理，减肥。
  PRIMARY KEY (pid),
  KEY (tid, pid)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

```sql
# 我的帖子
CREATE TABLE IF NOT EXISTS `bbs_mypost` (
  `uid` int(11) UNSIGNED NOT NULL DEFAULT '0',
  `tid` int(11) UNSIGNED NOT NULL DEFAULT '0',
  `pid` int(11) UNSIGNED NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`,`pid`),
  KEY `tid` (`tid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
```

### XiunoBBS
```
MySQL [xn]> desc bbs_post;
+-------------+------------------+------+-----+---------+----------------+
| Field       | Type             | Null | Key | Default | Extra          |
+-------------+------------------+------+-----+---------+----------------+
| tid         | int(11) unsigned | NO   | MUL | 0       |                |
| pid         | int(11) unsigned | NO   | PRI | NULL    | auto_increment |
| uid         | int(11) unsigned | NO   | MUL | 0       |                |
| isfirst     | int(11) unsigned | NO   |     | 0       |                |
| create_date | int(11) unsigned | NO   |     | 0       |                |
| userip      | int(11) unsigned | NO   |     | 0       |                |
| images      | smallint(6)      | NO   |     | 0       |                |
| files       | smallint(6)      | NO   |     | 0       |                |
| doctype     | tinyint(3)       | NO   |     | 0       |                |
| quotepid    | int(11)          | NO   |     | 0       |                |
| message     | longtext         | NO   |     | NULL    |                |
| message_fmt | longtext         | NO   |     | NULL    |                |
+-------------+------------------+------+-----+---------+----------------+
12 rows in set (0.00 sec)

MySQL [xn]> desc bbs_mypost;
+-------+------------------+------+-----+---------+-------+
| Field | Type             | Null | Key | Default | Extra |
+-------+------------------+------+-----+---------+-------+
| uid   | int(11) unsigned | NO   | PRI | 0       |       |
| tid   | int(11) unsigned | NO   | MUL | 0       |       |
| pid   | int(11) unsigned | NO   | PRI | 0       |       |
+-------+------------------+------+-----+---------+-------+
3 rows in set (0.00 sec)
```

### Discuz
```
MySQL [dx]> desc pre_forum_post;
+-------------+-----------------------+------+-----+---------+----------------+
| Field       | Type                  | Null | Key | Default | Extra          |
+-------------+-----------------------+------+-----+---------+----------------+
| pid         | int(10) unsigned      | NO   | UNI | NULL    |                |
| fid         | mediumint(8) unsigned | NO   | MUL | 0       |                |
| tid         | mediumint(8) unsigned | NO   | PRI | 0       |                |
| first       | tinyint(1)            | NO   |     | 0       |                |
| author      | varchar(15)           | NO   |     |         |                |
| authorid    | mediumint(8) unsigned | NO   | MUL | 0       |                |
| subject     | varchar(80)           | NO   |     |         |                |
| dateline    | int(10) unsigned      | NO   | MUL | 0       |                |
| message     | mediumtext            | NO   |     | NULL    |                |
| useip       | varchar(15)           | NO   |     |         |                |
| port        | smallint(6) unsigned  | NO   |     | 0       |                |
| invisible   | tinyint(1)            | NO   | MUL | 0       |                |
| anonymous   | tinyint(1)            | NO   |     | 0       |                |
| usesig      | tinyint(1)            | NO   |     | 0       |                |
| htmlon      | tinyint(1)            | NO   |     | 0       |                |
| bbcodeoff   | tinyint(1)            | NO   |     | 0       |                |
| smileyoff   | tinyint(1)            | NO   |     | 0       |                |
| parseurloff | tinyint(1)            | NO   |     | 0       |                |
| attachment  | tinyint(1)            | NO   |     | 0       |                |
| rate        | smallint(6)           | NO   |     | 0       |                |
| ratetimes   | tinyint(3) unsigned   | NO   |     | 0       |                |
| status      | int(10)               | NO   |     | 0       |                |
| tags        | varchar(255)          | NO   |     | 0       |                |
| comment     | tinyint(1)            | NO   |     | 0       |                |
| replycredit | int(10)               | NO   |     | 0       |                |
| position    | int(8) unsigned       | NO   | PRI | NULL    | auto_increment |
+-------------+-----------------------+------+-----+---------+----------------+
26 rows in set (0.02 sec)
```

### 对应关系 - post 帖子
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| tid         | tid               | 主题ID
| pid         | pid               | 帖子ID
| uid         | authorid          | 用户ID
| isfirst     | first             | 是否为首帖(主题)
| create_date | dateline          | 创建时间
| userip      | useip             | 创建IP
| images      | -<4>              | 图片数
| files       | -<4>              | 附件数
| doctype     | -(默认0.html)         | message 的类型，0: html, 1: txt; 2: markdown; 3: ubb
| quotepid    | -                 | 引用pid
| message     | message           | 内容原数据 (直接转 html 后的数据)
| message_fmt | message (将ubb转换后的)  | ubb转html后的message
+-----------+---------------------+------+-----+---------+----------------+
```

### 对应关系 - mypost 我的帖子
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| uid   | authorid                | 用户ID
| tid   | tid                     | 主题ID
| pid   | pid                     | 帖子ID
+-----------+---------------------+------+-----+---------+----------------+
```

## 备注
- 使用到 xiuno 的三表: bbs_post, bbs_mypost
- message_fmt 由 message内容并由ubb转html所得
- ✔align、table、font 的ubb标签暂时无法解析
- 图片数及附件数从 attach 表中提取
