用户组
------

**group(用户组) 对应表关系**
```sql
# 用户组表
DROP TABLE IF EXISTS `bbs_group`;
CREATE TABLE `bbs_group` (
  gid smallint(6) unsigned NOT NULL,			#	
  name char(20) NOT NULL default '',			# 用户组名称
  creditsfrom int(11) NOT NULL default '0',		# 积分从
  creditsto int(11) NOT NULL default '0',		# 积分到
  allowread int(11) NOT NULL default '0',		# 允许访问
  allowthread int(11) NOT NULL default '0',		# 允许发主题
  allowpost int(11) NOT NULL default '0',		# 允许回帖
  allowattach int(11) NOT NULL default '0',		# 允许上传文件
  allowdown int(11) NOT NULL default '0',		# 允许下载文件
  allowtop int(11) NOT NULL default '0',		# 允许置顶
  allowupdate int(11) NOT NULL default '0',		# 允许编辑
  allowdelete int(11) NOT NULL default '0',		# 允许删除
  allowmove int(11) NOT NULL default '0',		# 允许移动
  allowbanuser int(11) NOT NULL default '0',		# 允许禁止用户
  allowdeleteuser int(11) NOT NULL default '0',		# 允许删除用户
  allowviewip int(11) unsigned NOT NULL default '0',	# 允许查看用户敏感信息
  PRIMARY KEY (gid)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

### XiunoBBS
```
MySQL [xn]> desc bbs_group;
+-----------------+----------------------+------+-----+---------+-------+
| Field           | Type                 | Null | Key | Default | Extra |
+-----------------+----------------------+------+-----+---------+-------+
| gid             | smallint(6) unsigned | NO   | PRI | NULL    |       |
| name            | char(20)             | NO   |     |         |       |
| creditsfrom     | int(11)              | NO   |     | 0       |       |
| creditsto       | int(11)              | NO   |     | 0       |       |
| allowread       | int(11)              | NO   |     | 0       |       |
| allowthread     | int(11)              | NO   |     | 0       |       |
| allowpost       | int(11)              | NO   |     | 0       |       |
| allowattach     | int(11)              | NO   |     | 0       |       |
| allowdown       | int(11)              | NO   |     | 0       |       |
| allowtop        | int(11)              | NO   |     | 0       |       |
| allowupdate     | int(11)              | NO   |     | 0       |       |
| allowdelete     | int(11)              | NO   |     | 0       |       |
| allowmove       | int(11)              | NO   |     | 0       |       |
| allowbanuser    | int(11)              | NO   |     | 0       |       |
| allowdeleteuser | int(11)              | NO   |     | 0       |       |
| allowviewip     | int(11) unsigned     | NO   |     | 0       |       |
+-----------------+----------------------+------+-----+---------+-------+
16 rows in set (0.00 sec)
```

### Discuz
```
MySQL [dx]> desc pre_common_usergroup;
+-----------------+-----------------------------------+------+-----+---------+----------------+
| Field           | Type                              | Null | Key | Default | Extra          |
+-----------------+-----------------------------------+------+-----+---------+----------------+
| groupid         | smallint(6) unsigned              | NO   | PRI | NULL    | auto_increment |
| radminid        | tinyint(3)                        | NO   |     | 0       |                |
| type            | enum('system','special','member') | NO   |     | member  |                |
| system          | varchar(255)                      | NO   |     | private |                |
| grouptitle      | varchar(255)                      | NO   |     |         |                |
| creditshigher   | int(10)                           | NO   | MUL | 0       |                |
| creditslower    | int(10)                           | NO   |     | 0       |                |
| stars           | tinyint(3)                        | NO   |     | 0       |                |
| color           | varchar(255)                      | NO   |     |         |                |
| icon            | varchar(255)                      | NO   |     |         |                |
| allowvisit      | tinyint(1)                        | NO   |     | 0       |                |
| allowsendpm     | tinyint(1)                        | NO   |     | 1       |                |
| allowinvite     | tinyint(1)                        | NO   |     | 0       |                |
| allowmailinvite | tinyint(1)                        | NO   |     | 0       |                |
| maxinvitenum    | tinyint(3) unsigned               | NO   |     | 0       |                |
| inviteprice     | smallint(6) unsigned              | NO   |     | 0       |                |
| maxinviteday    | smallint(6) unsigned              | NO   |     | 0       |                |
+-----------------+-----------------------------------+------+-----+---------+----------------+
17 rows in set (0.00 sec)

MySQL [dx]> desc pre_common_usergroup_field;
+------------------------+-----------------------+------+-----+---------+-------+
| Field                  | Type                  | Null | Key | Default | Extra |
+------------------------+-----------------------+------+-----+---------+-------+
| groupid                | smallint(6) unsigned  | NO   | PRI | NULL    |       |
| readaccess             | tinyint(3) unsigned   | NO   |     | 0       |       |
| allowpost              | tinyint(1)            | NO   |     | 0       |       |
| allowreply             | tinyint(1)            | NO   |     | 0       |       |
| allowpostpoll          | tinyint(1)            | NO   |     | 0       |       |
| allowpostreward        | tinyint(1)            | NO   |     | 0       |       |
| allowposttrade         | tinyint(1)            | NO   |     | 0       |       |
| allowpostactivity      | tinyint(1)            | NO   |     | 0       |       |
| allowdirectpost        | tinyint(1)            | NO   |     | 0       |       |
| allowgetattach         | tinyint(1)            | NO   |     | 0       |       |
| allowgetimage          | tinyint(1)            | NO   |     | 0       |       |
| allowpostattach        | tinyint(1)            | NO   |     | 0       |       |
| allowpostimage         | tinyint(1)            | NO   |     | 0       |       |
| allowvote              | tinyint(1)            | NO   |     | 0       |       |
| allowsearch            | tinyint(1)            | NO   |     | 0       |       |
| allowcstatus           | tinyint(1)            | NO   |     | 0       |       |
| allowinvisible         | tinyint(1)            | NO   |     | 0       |       |
| allowtransfer          | tinyint(1)            | NO   |     | 0       |       |
| allowsetreadperm       | tinyint(1)            | NO   |     | 0       |       |
| allowsetattachperm     | tinyint(1)            | NO   |     | 0       |       |
| allowposttag           | tinyint(1)            | NO   |     | 0       |       |
| allowhidecode          | tinyint(1)            | NO   |     | 0       |       |
| allowhtml              | tinyint(1)            | NO   |     | 0       |       |
| allowanonymous         | tinyint(1)            | NO   |     | 0       |       |
| allowsigbbcode         | tinyint(1)            | NO   |     | 0       |       |
| allowsigimgcode        | tinyint(1)            | NO   |     | 0       |       |
| allowmagics            | tinyint(1) unsigned   | NO   |     | NULL    |       |
| disableperiodctrl      | tinyint(1)            | NO   |     | 0       |       |
| reasonpm               | tinyint(1)            | NO   |     | 0       |       |
| maxprice               | smallint(6) unsigned  | NO   |     | 0       |       |
| maxsigsize             | smallint(6) unsigned  | NO   |     | 0       |       |
| maxattachsize          | int(10) unsigned      | NO   |     | 0       |       |
| maxsizeperday          | int(10) unsigned      | NO   |     | 0       |       |
| maxthreadsperhour      | tinyint(3) unsigned   | NO   |     | 0       |       |
| maxpostsperhour        | tinyint(3) unsigned   | NO   |     | 0       |       |
| attachextensions       | char(100)             | NO   |     |         |       |
| raterange              | char(150)             | NO   |     |         |       |
| loginreward            | char(150)             | NO   |     |         |       |
| mintradeprice          | smallint(6) unsigned  | NO   |     | 1       |       |
| maxtradeprice          | smallint(6) unsigned  | NO   |     | 0       |       |
| minrewardprice         | smallint(6) unsigned  | NO   |     | 1       |       |
| maxrewardprice         | smallint(6) unsigned  | NO   |     | 0       |       |
| magicsdiscount         | tinyint(1)            | NO   |     | NULL    |       |
| maxmagicsweight        | smallint(6) unsigned  | NO   |     | NULL    |       |
| allowpostdebate        | tinyint(1)            | NO   |     | 0       |       |
| tradestick             | tinyint(1) unsigned   | NO   |     | NULL    |       |
| exempt                 | tinyint(1) unsigned   | NO   |     | NULL    |       |
| maxattachnum           | smallint(6)           | NO   |     | 0       |       |
| allowposturl           | tinyint(1)            | NO   |     | 3       |       |
| allowrecommend         | tinyint(1) unsigned   | NO   |     | 1       |       |
| allowpostrushreply     | tinyint(1)            | NO   |     | 0       |       |
| maxfriendnum           | smallint(6) unsigned  | NO   |     | 0       |       |
| maxspacesize           | int(10) unsigned      | NO   |     | 0       |       |
| allowcomment           | tinyint(1)            | NO   |     | 0       |       |
| allowcommentarticle    | smallint(6)           | NO   |     | 0       |       |
| searchinterval         | smallint(6) unsigned  | NO   |     | 0       |       |
| searchignore           | tinyint(1)            | NO   |     | 0       |       |
| allowblog              | tinyint(1)            | NO   |     | 0       |       |
| allowdoing             | tinyint(1)            | NO   |     | 0       |       |
| allowupload            | tinyint(1)            | NO   |     | 0       |       |
| allowshare             | tinyint(1)            | NO   |     | 0       |       |
| allowblogmod           | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowdoingmod          | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowuploadmod         | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowsharemod          | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowcss               | tinyint(1)            | NO   |     | 0       |       |
| allowpoke              | tinyint(1)            | NO   |     | 0       |       |
| allowfriend            | tinyint(1)            | NO   |     | 0       |       |
| allowclick             | tinyint(1)            | NO   |     | 0       |       |
| allowmagic             | tinyint(1)            | NO   |     | 0       |       |
| allowstat              | tinyint(1)            | NO   |     | 0       |       |
| allowstatdata          | tinyint(1)            | NO   |     | 0       |       |
| videophotoignore       | tinyint(1)            | NO   |     | 0       |       |
| allowviewvideophoto    | tinyint(1)            | NO   |     | 0       |       |
| allowmyop              | tinyint(1)            | NO   |     | 0       |       |
| magicdiscount          | tinyint(1)            | NO   |     | 0       |       |
| domainlength           | smallint(6) unsigned  | NO   |     | 0       |       |
| seccode                | tinyint(1)            | NO   |     | 1       |       |
| disablepostctrl        | tinyint(1)            | NO   |     | 0       |       |
| allowbuildgroup        | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowgroupdirectpost   | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowgroupposturl      | tinyint(1) unsigned   | NO   |     | 0       |       |
| edittimelimit          | smallint(6) unsigned  | NO   |     | 0       |       |
| allowpostarticle       | tinyint(1)            | NO   |     | 0       |       |
| allowdownlocalimg      | tinyint(1)            | NO   |     | 0       |       |
| allowdownremoteimg     | tinyint(1)            | NO   |     | 0       |       |
| allowpostarticlemod    | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowspacediyhtml      | tinyint(1)            | NO   |     | 0       |       |
| allowspacediybbcode    | tinyint(1)            | NO   |     | 0       |       |
| allowspacediyimgcode   | tinyint(1)            | NO   |     | 0       |       |
| allowcommentpost       | tinyint(1)            | NO   |     | 2       |       |
| allowcommentitem       | tinyint(1)            | NO   |     | 0       |       |
| allowcommentreply      | tinyint(1)            | NO   |     | 0       |       |
| allowreplycredit       | tinyint(1)            | NO   |     | 0       |       |
| ignorecensor           | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowsendallpm         | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowsendpmmaxnum      | smallint(6) unsigned  | NO   |     | 0       |       |
| maximagesize           | mediumint(8) unsigned | NO   |     | 0       |       |
| allowmediacode         | tinyint(1)            | NO   |     | 0       |       |
| allowbegincode         | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowat                | smallint(6) unsigned  | NO   |     | 0       |       |
| allowsetpublishdate    | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowfollowcollection  | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowcommentcollection | tinyint(1) unsigned   | NO   |     | 0       |       |
| allowcreatecollection  | smallint(6) unsigned  | NO   |     | 0       |       |
| forcesecques           | tinyint(1) unsigned   | NO   |     | 0       |       |
| forcelogin             | tinyint(1) unsigned   | NO   |     | 0       |       |
| closead                | tinyint(1) unsigned   | NO   |     | 0       |       |
| buildgroupcredits      | smallint(6) unsigned  | NO   |     | 0       |       |
| allowimgcontent        | tinyint(1) unsigned   | NO   |     | 0       |       |
+------------------------+-----------------------+------+-----+---------+-------+
110 rows in set (0.00 sec)

MySQL [dx]> desc pre_common_admingroup;
+-----------------------+----------------------+------+-----+---------+-------+
| Field                 | Type                 | Null | Key | Default | Extra |
+-----------------------+----------------------+------+-----+---------+-------+
| admingid              | smallint(6) unsigned | NO   | PRI | 0       |       |
| alloweditpost         | tinyint(1)           | NO   |     | 0       |       |
| alloweditpoll         | tinyint(1)           | NO   |     | 0       |       |
| allowstickthread      | tinyint(1)           | NO   |     | 0       |       |
| allowmodpost          | tinyint(1)           | NO   |     | 0       |       |
| allowdelpost          | tinyint(1)           | NO   |     | 0       |       |
| allowmassprune        | tinyint(1)           | NO   |     | 0       |       |
| allowrefund           | tinyint(1)           | NO   |     | 0       |       |
| allowcensorword       | tinyint(1)           | NO   |     | 0       |       |
| allowviewip           | tinyint(1)           | NO   |     | 0       |       |
| allowbanip            | tinyint(1)           | NO   |     | 0       |       |
| allowedituser         | tinyint(1)           | NO   |     | 0       |       |
| allowmoduser          | tinyint(1)           | NO   |     | 0       |       |
| allowbanuser          | tinyint(1)           | NO   |     | 0       |       |
| allowbanvisituser     | tinyint(1)           | NO   |     | 0       |       |
| allowpostannounce     | tinyint(1)           | NO   |     | 0       |       |
| allowviewlog          | tinyint(1)           | NO   |     | 0       |       |
| allowbanpost          | tinyint(1)           | NO   |     | 0       |       |
| supe_allowpushthread  | tinyint(1)           | NO   |     | 0       |       |
| allowhighlightthread  | tinyint(1)           | NO   |     | 0       |       |
| allowlivethread       | tinyint(1)           | NO   |     | 0       |       |
| allowdigestthread     | tinyint(1)           | NO   |     | 0       |       |
| allowrecommendthread  | tinyint(1)           | NO   |     | 0       |       |
| allowbumpthread       | tinyint(1)           | NO   |     | 0       |       |
| allowclosethread      | tinyint(1)           | NO   |     | 0       |       |
| allowmovethread       | tinyint(1)           | NO   |     | 0       |       |
| allowedittypethread   | tinyint(1)           | NO   |     | 0       |       |
| allowstampthread      | tinyint(1)           | NO   |     | 0       |       |
| allowstamplist        | tinyint(1)           | NO   |     | 0       |       |
| allowcopythread       | tinyint(1)           | NO   |     | 0       |       |
| allowmergethread      | tinyint(1)           | NO   |     | 0       |       |
| allowsplitthread      | tinyint(1)           | NO   |     | 0       |       |
| allowrepairthread     | tinyint(1)           | NO   |     | 0       |       |
| allowwarnpost         | tinyint(1)           | NO   |     | 0       |       |
| allowviewreport       | tinyint(1)           | NO   |     | 0       |       |
| alloweditforum        | tinyint(1)           | NO   |     | 0       |       |
| allowremovereward     | tinyint(1)           | NO   |     | 0       |       |
| allowedittrade        | tinyint(1)           | NO   |     | 0       |       |
| alloweditactivity     | tinyint(1)           | NO   |     | 0       |       |
| allowstickreply       | tinyint(1)           | NO   |     | 0       |       |
| allowmanagearticle    | tinyint(1)           | NO   |     | 0       |       |
| allowaddtopic         | tinyint(1)           | NO   |     | 0       |       |
| allowmanagetopic      | tinyint(1)           | NO   |     | 0       |       |
| allowdiy              | tinyint(1)           | NO   |     | 0       |       |
| allowclearrecycle     | tinyint(1)           | NO   |     | 0       |       |
| allowmanagetag        | tinyint(1)           | NO   |     | 0       |       |
| alloweditusertag      | tinyint(1)           | NO   |     | 0       |       |
| managefeed            | tinyint(1)           | NO   |     | 0       |       |
| managedoing           | tinyint(1)           | NO   |     | 0       |       |
| manageshare           | tinyint(1)           | NO   |     | 0       |       |
| manageblog            | tinyint(1)           | NO   |     | 0       |       |
| managealbum           | tinyint(1)           | NO   |     | 0       |       |
| managecomment         | tinyint(1)           | NO   |     | 0       |       |
| managemagiclog        | tinyint(1)           | NO   |     | 0       |       |
| managereport          | tinyint(1)           | NO   |     | 0       |       |
| managehotuser         | tinyint(1)           | NO   |     | 0       |       |
| managedefaultuser     | tinyint(1)           | NO   |     | 0       |       |
| managevideophoto      | tinyint(1)           | NO   |     | 0       |       |
| managemagic           | tinyint(1)           | NO   |     | 0       |       |
| manageclick           | tinyint(1)           | NO   |     | 0       |       |
| allowmanagecollection | tinyint(1)           | NO   |     | 0       |       |
| allowmakehtml         | tinyint(1)           | NO   |     | 0       |       |
+-----------------------+----------------------+------+-----+---------+-------+
62 rows in set (0.03 sec)
```

### 对应关系
```
+-----------+---------------------+------+-----+---------+----------------+
| XiunoBBS  | Discuz              |   描述
+--------------+----------------------+------+-----+---------+----------------+
| gid             | groupid                         | 用户组ID
| name            | grouptitle                      | 用户组名称
| creditsfrom     | creditslower                    | 积分从
| creditsto       | creditshigher                   | 积分到
| allowread       | allowvisit                      | 允许访问
| allowthread     | _field.allowpost                | 允许发主题
| allowpost       | _field.allowreply               | 允许回帖
| allowattach     | _field.allowpostattach          | 允许上传文件
| allowdown       | _field.allowgetattach           | 允许下载文件
| allowtop        | _admingroup.allowstickthread    | 允许置顶
| allowupdate     | _admingroup.alloweditpost       | 允许编辑
| allowdelete     | _admingroup.allowdelpost        | 允许删除
| allowmove       | _admingroup.allowmovethread     | 允许移动
| allowbanuser    | _admingroup.allowbanvisituser   | 允许禁止用户
| allowdeleteuser | -<4>                            | 允许删除用户
| allowviewip     | _admingroup.allowviewip         | 允许查看用户IP
+--------------+----------------------+------+-----+---------+----------------+
```

## 备注
- 使用到三表: pre_common_usergroup, pre_common_usergroup_field, pre_common_admingroup
- 当 type 值有三个 enum('system','special','member')
- 当 type == member 时, 填充管理组的值(allowtop,allowupdate,allowdelete,allowmove,allowbanuser,allowviewip)
- 删除用户功能手动修改数据库 (或者后面在扩展功能添加此功能)
- 扩展功能：因XiunoBBS 新用户注册固定用户组为 **101**，现决定将creditsfrom为0，creditsto不为0的组ID改为101