附件
------

**attach(附件) 对应表关系**

```sql
#论坛附件表  只能按照从上往下的方式查找和删除！ 此表如果大，可以考虑通过 aid 分区。
DROP TABLE IF EXISTS bbs_attach;
CREATE TABLE bbs_attach (
  aid int(11) unsigned NOT NULL auto_increment ,	# 附件id
  tid int(11) NOT NULL default '0',			# 主题id
  pid int(11) NOT NULL default '0',			# 帖子id
  uid int(11) NOT NULL default '0',			# 用户id
  filesize int(8) unsigned NOT NULL default '0',	# 文件尺寸，单位字节
  width mediumint(8) unsigned NOT NULL default '0',	# width > 0 则为图片
  height mediumint(8) unsigned NOT NULL default '0',	# height
  filename char(120) NOT NULL default '',		# 文件名称，会过滤，并且截断，保存后的文件名，不包含URL前缀 upload_url
  orgfilename char(120) NOT NULL default '',		# 上传的原文件名
  filetype char(7) NOT NULL default '',			# 文件类型: image/txt/zip，小图标显示 <i class="icon filetype image"></i>
  create_date int(11) unsigned NOT NULL default '0',	# 文件上传时间 UNIX 时间戳
  comment char(100) NOT NULL default '',		# 文件注释 方便于搜索
  downloads int(11) NOT NULL default '0',		# 下载次数，预留
  credits int(11) NOT NULL default '0',			# 需要的积分，预留
  golds int(11) NOT NULL default '0',			# 需要的金币，预留
  rmbs int(11) NOT NULL default '0',			# 需要的人民币，预留
  isimage tinyint(11) NOT NULL default '0',		# 是否为图片
  PRIMARY KEY (aid),					# aid
  KEY pid (pid),					# 每个帖子下多个附件
  KEY uid (uid)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

### XiunoBBS
```
MySQL [xn]> desc bbs_attach;
+-------------+-----------------------+------+-----+---------+----------------+
| Field       | Type                  | Null | Key | Default | Extra          |
+-------------+-----------------------+------+-----+---------+----------------+
| aid         | int(11) unsigned      | NO   | PRI | NULL    | auto_increment |
| tid         | int(11)               | NO   |     | 0       |                |
| pid         | int(11)               | NO   | MUL | 0       |                |
| uid         | int(11)               | NO   | MUL | 0       |                |
| filesize    | int(8) unsigned       | NO   |     | 0       |                |
| width       | mediumint(8) unsigned | NO   |     | 0       |                |
| height      | mediumint(8) unsigned | NO   |     | 0       |                |
| filename    | char(120)             | NO   |     |         |                |
| orgfilename | char(120)             | NO   |     |         |                |
| filetype    | char(7)               | NO   |     |         |                |
| create_date | int(11) unsigned      | NO   |     | 0       |                |
| comment     | char(100)             | NO   |     |         |                |
| downloads   | int(11)               | NO   |     | 0       |                |
| credits     | int(11)               | NO   |     | 0       |                |
| golds       | int(11)               | NO   |     | 0       |                |
| rmbs        | int(11)               | NO   |     | 0       |                |
| isimage     | tinyint(1)            | NO   |     | 0       |                |
+-------------+-----------------------+------+-----+---------+----------------+
17 rows in set (0.01 sec)
```

### Discuz
```
MySQL [dx]> desc pre_forum_attachment;
+-----------+-----------------------+------+-----+---------+----------------+
| Field     | Type                  | Null | Key | Default | Extra          |
+-----------+-----------------------+------+-----+---------+----------------+
| aid       | mediumint(8) unsigned | NO   | PRI | NULL    | auto_increment |
| tid       | mediumint(8) unsigned | NO   | MUL | 0       |                |
| pid       | int(10) unsigned      | NO   | MUL | 0       |                |
| uid       | mediumint(8) unsigned | NO   | MUL | 0       |                |
| tableid   | tinyint(1) unsigned   | 对应在附件表ID
| downloads | mediumint(8)          | 下载次数
+-----------+-----------------------+------+-----+---------+----------------+
6 rows in set (0.00 sec)

MySQL [dx]> desc pre_forum_attachment_0;
+-------------+-----------------------+------+-----+---------+-------+
| Field       | Type                  | Null | Key | Default | Extra |
+-------------+-----------------------+------+-----+---------+-------+
| aid         | mediumint(8) unsigned | NO   | PRI | NULL    |       |
| tid         | mediumint(8) unsigned | NO   | MUL | 0       |       |
| pid         | int(10) unsigned      | NO   | MUL | 0       |       |
| uid         | mediumint(8) unsigned | NO   | MUL | 0       |       |
| dateline    | int(10) unsigned      | NO   |     | 0       |       |
| filename    | varchar(255)          | NO   |     |         |       |
| filesize    | int(10) unsigned      | NO   |     | 0       |       |
| attachment  | varchar(255)          | NO   |     |         |       |
| remote      | tinyint(1) unsigned   | NO   |     | 0       |       |
| description | varchar(255)          | NO   |     | NULL    |       |
| readperm    | tinyint(3) unsigned   | NO   |     | 0       |       |
| price       | smallint(6) unsigned  | NO   |     | 0       |       |
| isimage     | tinyint(1)            | NO   |     | 0       |       |
| width       | smallint(6) unsigned  | NO   |     | 0       |       |
| thumb       | tinyint(1) unsigned   | NO   |     | 0       |       |
| picid       | mediumint(8)          | NO   |     | 0       |       |
+-------------+-----------------------+------+-----+---------+-------+
16 rows in set (0.00 sec)
```

### 对应关系 - attach 帖子
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+-----------+---------------------+------+-----+---------+----------------+
| aid         | aid <2>
| tid         | tid <2>
| pid         | pid <2>
| uid         | uid <2>
| filesize    | *.filesize
| width       | *.width
| height      | -
| filename    | *.attachment
| orgfilename | *.filename
| filetype    | - <5>
| create_date | *.dateline
| comment     | *.description
| downloads   | downloads <2>
| credits     | -
| golds       | -
| rmbs        | -
| isimage     | *.isimage
+-----------+---------------------+------+-----+---------+----------------+
```

## 备注
- 使用到 discuz 的两(分)表: pre_forum_attachment, pre_forum_attachment_x
- 直接提取附件分表即可
- upload/attach/201804/1_J8JZUCK2446ZVG3.jpg - upload/attach/
- data/attachment/forum/201804/11/094749cyy1evvv7051uf50.jpg - data/attachment/forum/
- filetype 分为 image,text,zip,other (建议从后缀名处提取)
- 重置 post 和 thread 的 images 和 files