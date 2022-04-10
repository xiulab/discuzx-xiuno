用户
------

**user(用户) 对应表关系**
```sql
# 用户表
DROP TABLE IF EXISTS `bbs_user`;
CREATE TABLE `bbs_user` (
  uid int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户编号',
  gid smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '用户组编号',	# 如果要屏蔽，调整用户组即可
  email char(40) NOT NULL DEFAULT '' COMMENT '邮箱',
  username char(32) NOT NULL DEFAULT '' COMMENT '用户名',	# 不可以重复
  realname char(16) NOT NULL DEFAULT '' COMMENT '用户名',	# 真实姓名，天朝预留
  idnumber char(19) NOT NULL DEFAULT '' COMMENT '用户名',	# 真实身份证号码，天朝预留
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `password_sms` char(16) NOT NULL DEFAULT '' COMMENT '密码',	# 预留，手机发送的 sms 验证码
  salt char(16) NOT NULL DEFAULT '' COMMENT '密码混杂',
  mobile char(11) NOT NULL DEFAULT '' COMMENT '手机号',		# 预留，供二次开发扩展
  qq char(15) NOT NULL DEFAULT '' COMMENT 'QQ',			# 预留，供二次开发扩展，可以弹出QQ直接聊天
  threads int(11) NOT NULL DEFAULT '0' COMMENT '发帖数',		#
  posts int(11) NOT NULL DEFAULT '0' COMMENT '回帖数',		#
  credits int(11) NOT NULL DEFAULT '0' COMMENT '积分',		# 预留，供二次开发扩展
  golds int(11) NOT NULL DEFAULT '0' COMMENT '金币',		# 预留，虚拟币
  rmbs int(11) NOT NULL DEFAULT '0' COMMENT '人民币',		# 预留，人民币
  create_ip int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时IP',
  create_date int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  login_ip int(11) unsigned NOT NULL DEFAULT '0' COMMENT '登录时IP',
  login_date int(11) unsigned NOT NULL DEFAULT '0' COMMENT '登录时间',
  logins int(11) unsigned NOT NULL DEFAULT '0' COMMENT '登录次数',
  avatar int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户最后更新图像时间',
  PRIMARY KEY (uid),
  UNIQUE KEY username (username),
  UNIQUE KEY email (email),						# 升级的时候可能为空
  KEY gid (gid)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8;
```

### XiunoBBS
```
MySQL [xn]> desc bbs_user;
+--------------+----------------------+------+-----+---------+----------------+
| Field        | Type                 | Null | Key | Default | Extra          |
+--------------+----------------------+------+-----+---------+----------------+
| uid          | int(11) unsigned     | NO   | PRI | NULL    | auto_increment |
| gid          | smallint(6) unsigned | NO   | MUL | 0       |                |
| email        | char(40)             | NO   | UNI |         |                |
| username     | char(32)             | NO   | UNI |         |                |
| realname     | char(16)             | NO   |     |         |                |
| idnumber     | char(19)             | NO   |     |         |                |
| password     | char(32)             | NO   |     |         |                |
| password_sms | char(16)             | NO   |     |         |                |
| salt         | char(16)             | NO   |     |         |                |
| mobile       | char(11)             | NO   |     |         |                |
| qq           | char(15)             | NO   |     |         |                |
| threads      | int(11)              | NO   |     | 0       |                |
| posts        | int(11)              | NO   |     | 0       |                |
| credits      | int(11)              | NO   |     | 0       |                |
| golds        | int(11)              | NO   |     | 0       |                |
| rmbs         | int(11)              | NO   |     | 0       |                |
| create_ip    | int(11) unsigned     | NO   |     | 0       |                |
| create_date  | int(11) unsigned     | NO   |     | 0       |                |
| login_ip     | int(11) unsigned     | NO   |     | 0       |                |
| login_date   | int(11) unsigned     | NO   |     | 0       |                |
| logins       | int(11) unsigned     | NO   |     | 0       |                |
| avatar       | int(11) unsigned     | NO   |     | 0       |                |
+--------------+----------------------+------+-----+---------+----------------+
22 rows in set (0.01 sec)
```

### Discuz
```
MySQL [dx]> desc pre_common_member;
+--------------------+-----------------------+------+-----+---------+----------------+
| Field              | Type                  | Null | Key | Default | Extra          |
+--------------------+-----------------------+------+-----+---------+----------------+
| uid                | mediumint(8) unsigned | NO   | PRI | NULL    | auto_increment |
| email              | char(40)              | NO   | MUL |         |                |
| username           | char(15)              | NO   | UNI |         |                |
| password           | char(32)              | NO   |     |         |                |
| status             | tinyint(1)            | NO   |     | 0       |                |
| emailstatus        | tinyint(1)            | NO   |     | 0       |                |
| avatarstatus       | tinyint(1)            | NO   |     | 0       |                |
| videophotostatus   | tinyint(1)            | NO   |     | 0       |                |
| adminid            | tinyint(1)            | NO   |     | 0       |                |
| groupid            | smallint(6) unsigned  | NO   | MUL | 0       |                |
| groupexpiry        | int(10) unsigned      | NO   |     | 0       |                |
| extgroupids        | char(20)              | NO   |     |         |                |
| regdate            | int(10) unsigned      | NO   | MUL | 0       |                |
| credits            | int(10)               | NO   |     | 0       |                |
| notifysound        | tinyint(1)            | NO   |     | 0       |                |
| timeoffset         | char(4)               | NO   |     |         |                |
| newpm              | smallint(6) unsigned  | NO   |     | 0       |                |
| newprompt          | smallint(6) unsigned  | NO   |     | 0       |                |
| accessmasks        | tinyint(1)            | NO   |     | 0       |                |
| allowadmincp       | tinyint(1)            | NO   |     | 0       |                |
| onlyacceptfriendpm | tinyint(1)            | NO   |     | 0       |                |
| conisbind          | tinyint(1) unsigned   | NO   | MUL | 0       |                |
| freeze             | tinyint(1)            | NO   |     | 0       |                |
+--------------------+-----------------------+------+-----+---------+----------------+
23 rows in set (0.00 sec)

MySQL [dx]> desc pre_ucenter_members;
+---------------+-----------------------+------+-----+---------+----------------+
| Field         | Type                  | Null | Key | Default | Extra          |
+---------------+-----------------------+------+-----+---------+----------------+
| uid           | mediumint(8) unsigned | NO   | PRI | NULL    | auto_increment |
| username      | char(15)              | NO   | UNI |         |                |
| password      | char(32)              | NO   |     |         |                |
| email         | char(32)              | NO   | MUL |         |                |
| myid          | char(30)              | NO   |     |         |                |
| myidkey       | char(16)              | NO   |     |         |                |
| regip         | char(15)              | NO   |     |         |                |
| regdate       | int(10) unsigned      | NO   |     | 0       |                |
| lastloginip   | int(10)               | NO   |     | 0       |                |
| lastlogintime | int(10) unsigned      | NO   |     | 0       |                |
| salt          | char(6)               | NO   |     | NULL    |                |
| secques       | char(8)               | NO   |     |         |                |
+---------------+-----------------------+------+-----+---------+----------------+
12 rows in set (0.00 sec)

MySQL [dx]> desc pre_common_member_status;
+-----------------+-----------------------+------+-----+---------+-------+
| Field           | Type                  | Null | Key | Default | Extra |
+-----------------+-----------------------+------+-----+---------+-------+
| uid             | mediumint(8) unsigned | NO   | PRI | NULL    |       |
| regip           | char(15)              | NO   |     |         |       |
| lastip          | char(15)              | NO   |     |         |       |
| port            | smallint(6) unsigned  | NO   |     | 0       |       |
| lastvisit       | int(10) unsigned      | NO   |     | 0       |       |
| lastactivity    | int(10) unsigned      | NO   | MUL | 0       |       |
| lastpost        | int(10) unsigned      | NO   |     | 0       |       |
| lastsendmail    | int(10) unsigned      | NO   |     | 0       |       |
| invisible       | tinyint(1)            | NO   |     | 0       |       |
| buyercredit     | smallint(6)           | NO   |     | 0       |       |
| sellercredit    | smallint(6)           | NO   |     | 0       |       |
| favtimes        | mediumint(8) unsigned | NO   |     | 0       |       |
| sharetimes      | mediumint(8) unsigned | NO   |     | 0       |       |
| profileprogress | tinyint(2) unsigned   | NO   |     | 0       |       |
+-----------------+-----------------------+------+-----+---------+-------+
14 rows in set (0.02 sec)
```

### 对应关系
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+--------------+----------------------+------+-----+---------+----------------+
| uid          | uid              | 用户ID
| gid          | groupid          | 分组ID
| email        | email            | 邮箱
| username     | username         | 用户名
| realname     | -                | 真实姓名
| idnumber     | -                | 身份证
| password     | ucenter_members.password        |
| password_sms | -                               | 短信登录验证码
| salt         | ucenter_members.salt            | 密码盐值
| mobile       | -                | 手机号
| qq           | -                | QQ
| threads      | -<4>             | 主题数
| posts        | -<4>             | 帖子数
| credits      | credits          | 积分
| golds        | -                | 金币
| rmbs         | -                | 人民币
| create_ip    | pre_common_member_status.regip <3>   | 创建IP
| create_date  | regdate                         | 注册时间
| login_ip     | pre_common_member_status.lastip | 最后登录IP
| login_date   | pre_common_member_status.lastvisit  | 最后登录时间
| logins       | -                               | 登录次数
| avatar       | - avatarstatus=1                | 是否有头像(有头像时，填更新头像时的时间戳)
+--------------+----------------------+------+-----+---------+----------------+
```

## 备注
- 使用到三表: pre_common_member, pre_ucenter_members, pre_common_member_status
- avatar 头像需要重置 avatarstatus=1 时 
> upload/avatar/000/1.png - /uc_server/data/avatar/000/00/00/01_avatar_middle.jpg
- create_ip,login_ip 涉及到IP部分，将正常IP转换为数值形式
- 转换完 threads 和 posts 后再更新统计