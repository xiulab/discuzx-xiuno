主题
------

**thread(主题)、thread_top(置顶帖)、mythread(我的主题) 对应表关系**

```sql
# 论坛主题
DROP TABLE IF EXISTS bbs_thread;
CREATE TABLE bbs_thread (
  fid smallint(6) NOT NULL default '0',			# 版块 id
  tid int(11) unsigned NOT NULL auto_increment,		# 主题id
  top tinyint(1) NOT NULL default '0',			# 置顶级别: 0: 普通主题, 1-3 置顶的顺序
  uid int(11) unsigned NOT NULL default '0',		# 用户id
  userip int(11) unsigned NOT NULL default '0',		# 发帖时用户ip ip2long()，主要用来清理
  subject char(128) NOT NULL default '',		# 主题
  create_date int(11) unsigned NOT NULL default '0',	# 发帖时间
  last_date int(11) unsigned NOT NULL default '0',	# 最后回复时间
  views int(11) unsigned NOT NULL default '0',		# 查看次数, 剥离出去，单独的服务，避免 cache 失效
  posts int(11) unsigned NOT NULL default '0',		# 回帖数
  images tinyint(6) NOT NULL default '0',		# 附件中包含的图片数
  files tinyint(6) NOT NULL default '0',		# 附件中包含的文件数
  mods tinyint(6) NOT NULL default '0',			# 预留：版主操作次数，如果 > 0, 则查询 modlog，显示斑竹的评分
  closed tinyint(1) unsigned NOT NULL default '0',	# 预留：是否关闭，关闭以后不能再回帖、编辑。
  firstpid int(11) unsigned NOT NULL default '0',	# 首贴 pid
  lastuid int(11) unsigned NOT NULL default '0',	# 最近参与的 uid
  lastpid int(11) unsigned NOT NULL default '0',	# 最后回复的 pid
  PRIMARY KEY (tid),					# 主键
  KEY (lastpid),					# 最后回复排序
  KEY (fid, tid),					# 发帖时间排序，正序。数据量大时可以考虑建立小表，对小表进行分区优化，只有数据量达到千万级以上时才需要。
  KEY (fid, lastpid)					# 顶贴时间排序，倒序
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

```sql
# 置顶主题
DROP TABLE IF EXISTS bbs_thread_top;
CREATE TABLE bbs_thread_top (
  fid smallint(6) NOT NULL default '0',			# 查找板块置顶
  tid int(11) unsigned NOT NULL default '0',		# tid
  top int(11) unsigned NOT NULL default '0',		# top: 0 是普通最新贴，> 0 置顶贴。
  PRIMARY KEY (tid),					#
  KEY (top, tid),					# 最新贴：top=0 order by tid desc / 全局置顶： top=3
  KEY (fid, top)					# 版块置顶的贴 fid=1 and top=1
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

```sql
# 我的主题，每个主题不管回复多少次，只记录一次。大表，需要分区。
DROP TABLE IF EXISTS bbs_mythread;
CREATE TABLE bbs_mythread (
  uid int(11) unsigned NOT NULL default '0',		# uid
  tid int(11) unsigned NOT NULL default '0',		# 用来清理，删除板块的时候需要
  PRIMARY KEY (uid, tid)				# 每一个帖子只能插入一次 unique
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

### XiunoBBS
```
mysql> desc bbs_thread;
+-------------+---------------------+------+-----+---------+----------------+
| Field       | Type                | Null | Key | Default | Extra          |
+-------------+---------------------+------+-----+---------+----------------+
| fid         | smallint(6)         | NO   | MUL | 0       |                |
| tid         | int(11) unsigned    | NO   | PRI | NULL    | auto_increment |
| top         | tinyint(1)          | NO   |     | 0       |                |
| uid         | int(11) unsigned    | NO   |     | 0       |                |
| userip      | int(11) unsigned    | NO   |     | 0       |                |
| subject     | char(128)           | NO   |     |         |                |
| create_date | int(11) unsigned    | NO   |     | 0       |                |
| last_date   | int(11) unsigned    | NO   |     | 0       |                |
| views       | int(11) unsigned    | NO   |     | 0       |                |
| posts       | int(11) unsigned    | NO   |     | 0       |                |
| images      | tinyint(6)          | NO   |     | 0       |                |
| files       | tinyint(6)          | NO   |     | 0       |                |
| mods        | tinyint(6)          | NO   |     | 0       |                |
| closed      | tinyint(1) unsigned | NO   |     | 0       |                |
| firstpid    | int(11) unsigned    | NO   |     | 0       |                |
| lastuid     | int(11) unsigned    | NO   |     | 0       |                |
| lastpid     | int(11) unsigned    | NO   | MUL | 0       |                |
+-------------+---------------------+------+-----+---------+----------------+
17 rows in set (0.00 sec)

mysql> desc bbs_thread_top;
+-------+------------------+------+-----+---------+-------+
| Field | Type             | Null | Key | Default | Extra |
+-------+------------------+------+-----+---------+-------+
| fid   | smallint(6)      | NO   | MUL | 0       |       |
| tid   | int(11) unsigned | NO   | PRI | 0       |       |
| top   | int(11) unsigned | NO   | MUL | 0       |       |
+-------+------------------+------+-----+---------+-------+
3 rows in set (0.00 sec)

mysql> desc bbs_mythread;
+-------+------------------+------+-----+---------+-------+
| Field | Type             | Null | Key | Default | Extra |
+-------+------------------+------+-----+---------+-------+
| uid   | int(11) unsigned | NO   | PRI | 0       |       |
| tid   | int(11) unsigned | NO   | PRI | 0       |       |
+-------+------------------+------+-----+---------+-------+
2 rows in set (0.00 sec)
```

### Discuz
```
mysql> desc pre_forum_thread;
+---------------+-----------------------+------+-----+---------+----------------+
| Field         | Type                  | Null | Key | Default | Extra          |
+---------------+-----------------------+------+-----+---------+----------------+
| tid           | mediumint(8) unsigned | NO   | PRI | NULL    | auto_increment |
| fid           | mediumint(8) unsigned | NO   | MUL | 0       |                |
| posttableid   | smallint(6) unsigned  | NO   |     | 0       |                |
| typeid        | smallint(6) unsigned  | NO   |     | 0       |                |
| sortid        | smallint(6) unsigned  | NO   | MUL | 0       |                |
| readperm      | tinyint(3) unsigned   | NO   |     | 0       |                |
| price         | smallint(6)           | NO   |     | 0       |                |
| author        | char(15)              | NO   |     |         |                |
| authorid      | mediumint(8) unsigned | NO   | MUL | 0       |                |
| subject       | char(80)              | NO   |     |         |                |
| dateline      | int(10) unsigned      | NO   |     | 0       |                |
| lastpost      | int(10) unsigned      | NO   |     | 0       |                |
| lastposter    | char(15)              | NO   |     |         |                |
| views         | int(10) unsigned      | NO   |     | 0       |                |
| replies       | mediumint(8) unsigned | NO   |     | 0       |                |
| displayorder  | tinyint(1)            | NO   | 置顶  1本版,2分类,3全局          |
| highlight     | tinyint(1)            | NO   |     | 0       |                |
| digest        | tinyint(1)            | NO   | MUL | 0       |                |
| rate          | tinyint(1)            | NO   |     | 0       |                |
| special       | tinyint(1)            | NO   | MUL | 0       |                |
| attachment    | tinyint(1)            | NO   |     | 0       |                |
| moderated     | tinyint(1)            | NO   |     | 0       |                |
| closed        | mediumint(8) unsigned | NO   |     | 0       |                |
| stickreply    | tinyint(1) unsigned   | NO   |     | 0       |                |
| recommends    | smallint(6)           | NO   | MUL | 0       |                |
| recommend_add | smallint(6)           | NO   |     | 0       |                |
| recommend_sub | smallint(6)           | NO   |     | 0       |                |
| heats         | int(10) unsigned      | NO   | MUL | 0       |                |
| status        | smallint(6) unsigned  | NO   |     | 0       |                |
| isgroup       | tinyint(1)            | NO   | MUL | 0       |                |
| favtimes      | mediumint(8)          | NO   |     | 0       |                |
| sharetimes    | mediumint(8)          | NO   |     | 0       |                |
| stamp         | tinyint(3)            | NO   |     | -1      |                |
| icon          | tinyint(3)            | NO   |     | -1      |                |
| pushedaid     | mediumint(8)          | NO   |     | 0       |                |
| cover         | smallint(6)           | NO   |     | 0       |                |
| replycredit   | smallint(6)           | NO   |     | 0       |                |
| relatebytag   | char(255)             | NO   |     | 0       |                |
| maxposition   | int(8) unsigned       | NO   |     | 0       |                |
| bgcolor       | char(8)               | NO   |     |         |                |
| comments      | int(10) unsigned      | NO   |     | 0       |                |
| hidden        | smallint(6) unsigned  | NO   |     | 0       |                |
+---------------+-----------------------+------+-----+---------+----------------+
42 rows in set (0.00 sec)
```

### 对应关系 - thread 主题
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| fid         | fid                     | 归属版块
| tid         | tid                     | 主题 ID
| top         | displayorder            | 置顶
| uid         | authorid                | 作者 ID
| userip      | - forum_post.useip      | 发帖 IP
| subject     | subject                 | 标题
| create_date | dateline                | 创建时间
| last_date   | lastpost                | 最后回复时间
| views       | views                   | 浏览次数
| posts       | replies                 | 回复次数
| images      | -<4>                    | 图片数
| files       | -<4>                    | 文件数
| mods        | -                       | 版主修改次数
| closed      | closed                  | 是否已关闭
| firstpid    | - forum_post.pid > first| 主帖 ID
| lastuid     | - lastposter <3>        | 最后回复者 ID
| lastpid     | - <3>                   | 最后回复者 IP
+-----------+---------------------+------+-----+---------+----------------+
```

### 对应关系 - thread_top 置顶帖
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| fid   | fid                     | 版块 ID
| tid   | tid                     | 主题 ID
| top   | displayorder (>0)       | 置顶状态 1-3
+-----------+---------------------+------+-----+---------+----------------+
```

### 对应关系 - mythread 我的主题
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| uid   | authorid                | 用户 ID
| tid   | tid                     | 主题 ID
+-----------+---------------------+------+-----+---------+----------------+
```

## 备注
- 使用到 xiuno 的三表: bbs_thread, bbs_thread_top, bbs_mythread
- 使用到 dx 的表 pre_forum_post
- 最后发帖者及最后帖子 PID 最终扩展处再处理
- 图片数及附件数从 attach 表中提取