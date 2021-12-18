-- all values are handled by testing, we just need strucutre
CREATE TABLE `alias` (
  `address` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  `accesspolicy` varchar(30) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`address`),
  KEY `domain` (`domain`),
  KEY `expired` (`expired`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `alias_domain` (
  `alias_domain` varchar(255) NOT NULL,
  `target_domain` varchar(255) NOT NULL,
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`alias_domain`),
  KEY `target_domain` (`target_domain`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `domain` (
  `domain` varchar(255) NOT NULL DEFAULT '',
  `description` text DEFAULT NULL,
  `disclaimer` text DEFAULT NULL,
  `aliases` int(10) NOT NULL DEFAULT 0,
  `mailboxes` int(10) NOT NULL DEFAULT 0,
  `maillists` int(10) NOT NULL DEFAULT 0,
  `maxquota` bigint(20) NOT NULL DEFAULT 0,
  `quota` bigint(20) NOT NULL DEFAULT 0,
  `transport` varchar(255) NOT NULL DEFAULT 'dovecot',
  `backupmx` tinyint(1) NOT NULL DEFAULT 0,
  `settings` text DEFAULT NULL,
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`domain`),
  KEY `backupmx` (`backupmx`),
  KEY `expired` (`expired`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `forwardings` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(255) NOT NULL DEFAULT '',
  `forwarding` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `dest_domain` varchar(255) NOT NULL DEFAULT '',
  `is_maillist` tinyint(1) NOT NULL DEFAULT 0,
  `is_list` tinyint(1) NOT NULL DEFAULT 0,
  `is_forwarding` tinyint(1) NOT NULL DEFAULT 0,
  `is_alias` tinyint(1) NOT NULL DEFAULT 0,
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`,`forwarding`),
  KEY `domain` (`domain`),
  KEY `dest_domain` (`dest_domain`),
  KEY `is_maillist` (`is_maillist`),
  KEY `is_list` (`is_list`),
  KEY `is_alias` (`is_alias`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8;

--
-- Table structure for table `last_login`
--

DROP TABLE IF EXISTS `last_login`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `last_login` (
  `username` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `imap` int(11) DEFAULT NULL,
  `pop3` int(11) DEFAULT NULL,
  `lda` int(11) DEFAULT NULL,
  PRIMARY KEY (`username`),
  KEY `domain` (`domain`),
  KEY `imap` (`imap`),
  KEY `pop3` (`pop3`),
  KEY `lda` (`lda`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mailbox`
--

DROP TABLE IF EXISTS `mailbox`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mailbox` (
  `username` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  `language` varchar(5) NOT NULL DEFAULT '',
  `mailboxformat` varchar(50) NOT NULL DEFAULT 'maildir',
  `mailboxfolder` varchar(50) NOT NULL DEFAULT 'Maildir',
  `storagebasedirectory` varchar(255) NOT NULL DEFAULT '',
  `storagenode` varchar(255) NOT NULL DEFAULT '',
  `maildir` varchar(255) NOT NULL DEFAULT '',
  `quota` bigint(20) NOT NULL DEFAULT 0,
  `domain` varchar(255) NOT NULL DEFAULT '',
  `transport` varchar(255) NOT NULL DEFAULT '',
  `department` varchar(255) NOT NULL DEFAULT '',
  `rank` varchar(255) NOT NULL DEFAULT 'normal',
  `employeeid` varchar(255) DEFAULT '',
  `isadmin` tinyint(1) NOT NULL DEFAULT 0,
  `isglobaladmin` tinyint(1) NOT NULL DEFAULT 0,
  `enablesmtp` tinyint(1) NOT NULL DEFAULT 1,
  `enablesmtpsecured` tinyint(1) NOT NULL DEFAULT 1,
  `enablepop3` tinyint(1) NOT NULL DEFAULT 1,
  `enablepop3secured` tinyint(1) NOT NULL DEFAULT 1,
  `enablepop3tls` tinyint(1) NOT NULL DEFAULT 1,
  `enableimap` tinyint(1) NOT NULL DEFAULT 1,
  `enableimapsecured` tinyint(1) NOT NULL DEFAULT 1,
  `enableimaptls` tinyint(1) NOT NULL DEFAULT 1,
  `enabledeliver` tinyint(1) NOT NULL DEFAULT 1,
  `enablelda` tinyint(1) NOT NULL DEFAULT 1,
  `enablemanagesieve` tinyint(1) NOT NULL DEFAULT 1,
  `enablemanagesievesecured` tinyint(1) NOT NULL DEFAULT 1,
  `enablesieve` tinyint(1) NOT NULL DEFAULT 1,
  `enablesievesecured` tinyint(1) NOT NULL DEFAULT 1,
  `enablesievetls` tinyint(1) NOT NULL DEFAULT 1,
  `enableinternal` tinyint(1) NOT NULL DEFAULT 1,
  `enabledoveadm` tinyint(1) NOT NULL DEFAULT 1,
  `enablelib-storage` tinyint(1) NOT NULL DEFAULT 1,
  `enablequota-status` tinyint(1) NOT NULL DEFAULT 1,
  `enableindexer-worker` tinyint(1) NOT NULL DEFAULT 1,
  `enablelmtp` tinyint(1) NOT NULL DEFAULT 1,
  `enabledsync` tinyint(1) NOT NULL DEFAULT 1,
  `enablesogo` tinyint(1) NOT NULL DEFAULT 1,
  `enablesogowebmail` varchar(1) NOT NULL DEFAULT 'y',
  `enablesogocalendar` varchar(1) NOT NULL DEFAULT 'y',
  `enablesogoactivesync` varchar(1) NOT NULL DEFAULT 'y',
  `allow_nets` text DEFAULT NULL,
  `disclaimer` text DEFAULT NULL,
  `settings` text DEFAULT NULL,
  `passwordlastchange` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`username`),
  KEY `domain` (`domain`),
  KEY `department` (`department`),
  KEY `employeeid` (`employeeid`),
  KEY `isadmin` (`isadmin`),
  KEY `isglobaladmin` (`isglobaladmin`),
  KEY `enablesmtp` (`enablesmtp`),
  KEY `enablesmtpsecured` (`enablesmtpsecured`),
  KEY `enablepop3` (`enablepop3`),
  KEY `enablepop3secured` (`enablepop3secured`),
  KEY `enablepop3tls` (`enablepop3tls`),
  KEY `enableimap` (`enableimap`),
  KEY `enableimapsecured` (`enableimapsecured`),
  KEY `enableimaptls` (`enableimaptls`),
  KEY `enabledeliver` (`enabledeliver`),
  KEY `enablelda` (`enablelda`),
  KEY `enablemanagesieve` (`enablemanagesieve`),
  KEY `enablemanagesievesecured` (`enablemanagesievesecured`),
  KEY `enablesieve` (`enablesieve`),
  KEY `enablesievesecured` (`enablesievesecured`),
  KEY `enablesievetls` (`enablesievetls`),
  KEY `enablelmtp` (`enablelmtp`),
  KEY `enableinternal` (`enableinternal`),
  KEY `enabledoveadm` (`enabledoveadm`),
  KEY `enablelib-storage` (`enablelib-storage`),
  KEY `enablequota-status` (`enablequota-status`),
  KEY `enableindexer-worker` (`enableindexer-worker`),
  KEY `enabledsync` (`enabledsync`),
  KEY `enablesogo` (`enablesogo`),
  KEY `passwordlastchange` (`passwordlastchange`),
  KEY `expired` (`expired`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `maillist_owners`
--

DROP TABLE IF EXISTS `maillist_owners`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `maillist_owners` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(255) NOT NULL DEFAULT '',
  `owner` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `dest_domain` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`,`owner`),
  KEY `owner` (`owner`),
  KEY `domain` (`domain`),
  KEY `dest_domain` (`dest_domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `maillists`
--

DROP TABLE IF EXISTS `maillists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `maillists` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `transport` varchar(255) NOT NULL DEFAULT '',
  `accesspolicy` varchar(30) NOT NULL DEFAULT '',
  `maxmsgsize` bigint(20) NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL DEFAULT '',
  `description` text DEFAULT NULL,
  `mlid` varchar(36) NOT NULL DEFAULT '',
  `is_newsletter` tinyint(1) NOT NULL DEFAULT 0,
  `settings` text DEFAULT NULL,
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`),
  UNIQUE KEY `mlid` (`mlid`),
  KEY `is_newsletter` (`is_newsletter`),
  KEY `domain` (`domain`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `moderators`
--

DROP TABLE IF EXISTS `moderators`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `moderators` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(255) NOT NULL DEFAULT '',
  `moderator` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `dest_domain` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`,`moderator`),
  KEY `domain` (`domain`),
  KEY `dest_domain` (`dest_domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `recipient_bcc_domain`
--

DROP TABLE IF EXISTS `recipient_bcc_domain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `recipient_bcc_domain` (
  `domain` varchar(255) NOT NULL DEFAULT '',
  `bcc_address` varchar(255) NOT NULL DEFAULT '',
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`domain`),
  KEY `bcc_address` (`bcc_address`),
  KEY `expired` (`expired`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `recipient_bcc_user`
--

DROP TABLE IF EXISTS `recipient_bcc_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `recipient_bcc_user` (
  `username` varchar(255) NOT NULL DEFAULT '',
  `bcc_address` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`username`),
  KEY `bcc_address` (`bcc_address`),
  KEY `expired` (`expired`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sender_bcc_domain`
--

DROP TABLE IF EXISTS `sender_bcc_domain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sender_bcc_domain` (
  `domain` varchar(255) NOT NULL DEFAULT '',
  `bcc_address` varchar(255) NOT NULL DEFAULT '',
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`domain`),
  KEY `bcc_address` (`bcc_address`),
  KEY `expired` (`expired`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sender_bcc_user`
--

DROP TABLE IF EXISTS `sender_bcc_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sender_bcc_user` (
  `username` varchar(255) NOT NULL DEFAULT '',
  `bcc_address` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `created` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `modified` datetime NOT NULL DEFAULT '1970-01-01 01:01:01',
  `expired` datetime NOT NULL DEFAULT '9999-12-31 00:00:00',
  `active` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`username`),
  KEY `bcc_address` (`bcc_address`),
  KEY `domain` (`domain`),
  KEY `expired` (`expired`),
  KEY `active` (`active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sender_relayhost`
--

DROP TABLE IF EXISTS `sender_relayhost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sender_relayhost` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(255) NOT NULL DEFAULT '',
  `relayhost` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `share_folder`
--

DROP TABLE IF EXISTS `share_folder`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `share_folder` (
  `from_user` varchar(255) CHARACTER SET ascii NOT NULL,
  `to_user` varchar(255) CHARACTER SET ascii NOT NULL,
  `dummy` char(1) DEFAULT NULL,
  PRIMARY KEY (`from_user`,`to_user`),
  KEY `from_user` (`from_user`),
  KEY `to_user` (`to_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `used_quota`
--

DROP TABLE IF EXISTS `used_quota`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `used_quota` (
  `username` varchar(255) NOT NULL,
  `bytes` bigint(20) NOT NULL DEFAULT 0,
  `messages` bigint(20) NOT NULL DEFAULT 0,
  `domain` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`username`),
  KEY `domain` (`domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;