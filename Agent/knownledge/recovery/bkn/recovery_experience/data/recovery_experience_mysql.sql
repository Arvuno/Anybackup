-- MySQL dump for Recovery Experience Knowledge Network
--
-- Host: 127.0.0.1    Database: recovery_experience
-- ------------------------------------------------------
-- Server version	MySQL 8.0+

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

--
-- Table structure for table `availability_checkpoint_template`
--

DROP TABLE IF EXISTS `availability_checkpoint_template`;

CREATE TABLE `availability_checkpoint_template` (
  `checkpointTemplateId` varchar(128) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `appType` varchar(32) DEFAULT NULL,
  `checkpointType` varchar(32) DEFAULT NULL,
  `targetScope` varchar(32) DEFAULT NULL,
  `verificationMethod` text DEFAULT NULL,
  `successCriteria` text DEFAULT NULL,
  `faultPatternId` varchar(128) DEFAULT NULL,
  `strategyTemplateId` varchar(128) DEFAULT NULL,
  `enabled` boolean DEFAULT NULL,
  PRIMARY KEY (`checkpointTemplateId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `availability_checkpoint_template`
--

INSERT INTO `availability_checkpoint_template` VALUES 
('ckpt_mysql_basic_sql_verification','MySQL基础SQL检索验证模板','MySQL','sql_query','instance','连接 MySQL 实例并执行 select User from user，验证 mysql.user 表是否可正常检索。','SQL 执行成功并返回 result=true.',NULL,NULL,true);

--
-- Table structure for table `fault_pattern`
--

DROP TABLE IF EXISTS `fault_pattern`;

CREATE TABLE `fault_pattern` (
  `patternId` varchar(128) NOT NULL,
  `name`' varchar(255) DEFAULT NULL,
  `appType` varchar(32) DEFAULT NULL,
  `faultCategory` varchar(32) DEFAULT NULL,
  `affectedGranularity` varchar(64) DEFAULT NULL,
  `symptomKeywords` text DEFAULT NULL,
  `intentKeywords` text DEFAULT NULL,
  `requiredClarification` text DEFAULT NULL,
  `disposalHint` text DEFAULT NULL,
  `severityBaseline` varchar(16) DEFAULT NULL,
  `enabled` boolean DEFAULT NULL,
  PRIMARY KEY (`patternId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `fault_pattern`
--

--

INSERT INTO `fault_pattern` VALUES 
('mysql_service_unavailable','业务系统不可用','MySQL','service_outage','instance_or_database','无法连接,服务停止,宕机,访问失败,连接异常,业务不可用','恢复业务,修复数据库服务,恢复系统可用性','需明确提供业务的是实例还是数据库；需明确是全部业务不可用还是部分业务不可用；需明确故障发生时间。','先明确受影响粒度；若实例整体不可用，优先考虑实例级恢复；若仅部分数据库受影响，优先考虑库级恢复。','high',true),
('mysql_data_loss','数据丢失','MySQL','data_loss','instance_or_database_or_table_or_log','误删除,DROP,TRUNCATE,数据不见了,数据缺失,记录丢失','恢复数据,找回误删数据,恢复丢失数据','需明确数据丢失时间；需明确影响范围；需明确是实例、库、表还是日志时间段；需明确是否要求恢复到指定时间点。','优先选择故障前最近可用时间点；根据确认后的范围决定实例级、库级、表级间接或日志级恢复。','high',true),
('mysql_data_corruption','数据损坏导致对象访问异常','MySQL','data_corruption','instance_or_database_or_table','数据页损坏,索引损坏,表无法访问,查询报错,表损坏,逻辑损坏,物理损坏','修复数据损坏,恢复损坏对象,恢复异常读写,修复表损坏','需明确故障时间；需明确影响范围；需明确是单表、多个表、库还是实例受损；需明确当前是否还能部分读写。','优先评估受损对象范围；表级问题优先采用间接恢复；范围扩大时再考虑库级或实例级恢复。','high',true),
('mysql_ransomware','勒索或数据被加密','MySQL','ransomware','instance_or_database_or_table','被加密,勒索,文件被加密,数据被锁定,异常加密','恢复被加密数据,隔离恢复,恢复业务可用','需明确影响范围；需明确受影响时间范围；需明确是否要求恢复回原环境；需明确隔离环境要求。','优先在隔离环境恢复；避免原环境二次污染；默认不建议直接恢复回原环境。','critical',true),
('mysql_configuration_error','配置错误或配置异常','MySQL','configuration_error','instance_or_database','配置文件损坏,参数错误,无法启动,配置异常,启动失败','修复配置,恢复启动,恢复服务可用',''需明确是纯配置问题、实例问题，还是已涉及数据库对象损坏；需明确是否还能连接实例。','若未涉及数据损坏，默认不进入恢复流程；优先进行配置修正与启动验证。','medium',true);

--
-- Table structure for table `recovery_capability`
--

DROP TABLE IF EXISTS `recovery_capability`;

CREATE TABLE `recovery_capability` (
  `capabilityId` varchar(128) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `vendor` varchar(32) DEFAULT NULL,
  `appType` varchar(32) DEFAULT NULL,
  `supportedGranularity` varchar(255) DEFAULT NULL,
  `supportedTechnique` varchar(255) DEFAULT NULL,
  `supportedMode` varchar(255) DEFAULT NULL,
  `supportsOriginalRestore` boolean DEFAULT NULL,
  `supportsRemoteRestore` boolean DEFAULT NULL,
  `supportsPointInTimeRestore` boolean DEFAULT NULL,
  `supportsLogRestore` boolean DEFAULT NULL,
  `tableRecoveryMode` varchar(32) DEFAULT NULL,
  `capabilitySummary` text DEFAULT NULL,
  `enabled` boolean DEFAULT NULL,
  PRIMARY KEY (`capabilityId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `recovery_capability`
--

INSERT INTO `recovery_capability` VALUES 
('mysql_recovery_capability_v1','MySQL恢复能力集','AISHU','MySQL','instance,database,log_file,table_indirect','restore,mount','shortest_time,latest_state,specified_time',true,true,true,true,'indirect_only','支持实例级、库级、日志文件级恢复；表级仅支持间接恢复；支持数据恢复与挂载恢复(暂不支持，占坑用)；支持原机和异机；支持最短时间恢复、恢复到最新状态和恢复到指定时间。',true);

--
-- Table structure for table `recovery_strategy_template`
--

DROP TABLE IF EXISTS `recovery_strategy_template`;

CREATE TABLE `recovery_strategy_template` (
  `strategyTemplateId` varchar(128) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `vendor` varchar(32) DEFAULT NULL,
  `appType` varchar(32) DEFAULT NULL,
  `faultPatternId` varchar(128) DEFAULT NULL,
  `recoveryGranularity` varchar(32) DEFAULT NULL,
  `destinationType` varchar(32) DEFAULT NULL,
  `recoveryMethod` varchar(32) DEFAULT NULL,
  `requiresRecovery` boolean DEFAULT NULL,
  `strategySummary` text DEFAULT NULL,
  `riskBaseline` varchar(16) DEFAULT NULL,
  `approvalRequired` boolean DEFAULT NULL,
  `enabled`' boolean DEFAULT NULL,
  PRIMARY KEY (`strategyTemplateId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `recovery_strategy_template`
--

INSERT INTO `recovery_strategy_template` VALUES 
('mysql_instance_latest_restore','MySQL实例级最新点恢复模板','AISHU','MySQL','mysql_service_unavailable','instance','original','latest_state',true,'实例故障默认优先原机恢复到最新状态；该操作可能覆盖生产实例，执行前必须明确告知高风险并获得用户显式确认。','high',true,true),
('mysql_database_latest_restore','MySQL库级最新点恢复模板','AISHU','MySQL','mysql_service_unavailable','database','original','latest_state',true,'当确认仅某个或某些数据库不可用时，默认优先原机库级恢复到最新状态；该操作可能覆盖生产数据库，执行前必须明确告知高风险并获得用户显式确认。','high',true,true),
('mysql_instance_before_fault_restore','MySQL实例级定点恢复模板','AISHU','MySQL','mysql_data_loss','instance','original','specified_time',true,'当确认数据丢失范围为实例级时，需先明确故障发生时间，再选择故障前最近可用时间点默认恢复到原机实例；若原机 client 离线或不可用，再转异机恢复。','high',true,true),
('mysql_database_before_fault_restore','MySQL库级定点恢复模板','AISHU','MySQL','mysql_data_loss','database','original','specified_time',true,'当确认数据丢失范围为数据库级时，需先明确故障发生时间，再选择故障前最近可用时间点默认恢复到原机生产库；若原机 client 离线或不可用，再转异机恢复。','high',true,true),
('mysql_table_indirect_specified_time_restore','MySQL表级间接定点恢复模板','AISHU','MySQL','mysql_data_loss','database','remote','specified_time',true,'当确认问题范围为表级且需要间接恢复时，需先明确故障发生时间，再选择故障前最近可用时间点进行异机库级恢复，随后导出目标表或受影响对象并回灌。','high',true,true),
('mysql_log_file'_time_range_restore','MySQL日志文件时间范围恢复模板','AISHU','MySQL','mysql_data_loss','log_file','original','specified_time',true,'当确认数据丢失范围为日志时间段或需要恢复 binlog 文件时，按指定时间范围默认恢复归档日志文件到原机目标路径；若原机 client 离线线或不可用，再转异机恢复。','high',true,true),
('mysql_ransomware_isolated_restore','MySQL勒索场景隔离恢复模板','AISHU','MySQL','mysql_ransomware','instance','remote','specified_time',true,'勒索场景必须恢复到网络隔离环境，并优先恢复到受影响前的安全时间点。','high',true,true),
('mysql_configuration_no_restore','MySQL配置错误免恢复模板','AISHU','MySQL','mysql_configuration_error','none','none','none',false,'配置错误默认不恢复，先做配置修正。','low',false,true);

--
-- Table structure for table `risk_rule`
--

DROP TABLE IF EXISTS `risk_rule`;

CREATE TABLE `risk_rule` (
  `riskRuleId` varchar(128) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `appType` varchar(32) DEFAULT NULL,
  `strategyTemplateId` varchar(128) DEFAULT NULL,
  `triggerCondition` text DEFAULT NULL,
  `riskLevel` varchar(16) DEFAULT NULL,
  `approvalRequired` boolean DEFAULT NULL,
  `mitigationAdvice` text DEFAULT NULL,
  `enabled` boolean DEFAULT NULL,
  PRIMARY KEY (`riskRuleId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `risk_rule`
--

INSERT INTO `risk_rule` VALUES 
('risk_instance_original_restore','原机实例恢复风险控制','MySQL','mysql_instance_latest_restore','原机恢复且可能覆盖生产实例，属于高风险覆盖操作，未获得用户显式确认前不得自动执行','high',true,'先明确告知覆盖生产实例风险，再确认业务切换与回退方案，并取得用户显式确认后方可执行。',true),
('risk_database_original_restore','原机库级恢复风险控制','MySQL','mysql_database_latest_restore','原机库级恢复且可能覆盖生产数据库，属于高风险覆盖操作，未获得用户显式确认前不得自动执行','high',true,'先明确告知覆盖生产数据库风险，再确认业务影响范围与回退方案，并取得用户显式确认后方可执行。',true),
('risk_instance_restore_backfill_prod','实例级定点恢复风险控制','MySQL','mysql_instance_before_fault_restore','实例级定点恢复默认恢复到原机，可能覆盖生产实例；若原机 client 不可用转异机恢复，切换生产仍有较高风险。','high',true,'先确认故障时间、业务影响范围和回退方案；原机恢复前确认覆盖影响；若转异机恢复，需先验证实例可用性后再人工审批切换。',true),
('risk_database_restore_backfill_prod','库级定点恢复风险控制','MySQL','mysql_database_before_fault_restore','库级定点恢复默认恢复到原机生产库，可能覆盖现网数据库；若原机 client 不可用转异机恢复，切换或回灌生产仍有较高风险。','high',true,'先确认故障时间、受影响数据库范围和回退方案；原机恢复前确认覆盖影响；若转异机恢复，需先隔离验证后再人工审批切换或回灌。',true),
('risk_table_indirect_backfill_prod','表级间接恢复回灌生产风险控制','MySQL','mysql_table_indirect_specified_time_restore','需要将异机恢复数据库中的目标表或受影响对象回灌生产','high',true,'先隔离校验，再人工审批回灌；禁止直接覆盖生产对象。',true),
('risk_log_file_restore_prod','日志文件恢复风险控制','MySQL','mysql_log_file_time_range_restore','日志文件恢复默认恢复到原机目标路径，可能覆盖已有归档日志或影响后续日志链使用；若原机 client 不可用转异机恢复，也需确认目标路径和后续恢复链影响。','high',true,'先确认时间范围、目标路径和日志链连续性；避免覆盖已有日志；必要时先在隔离环境验证恢复日志文件可用性。',true),
('risk_ransomware_original_env','勒索场景回原环境风险控制','MySQL','mysql_ransomware_isolated_restore','用户要求恢复回原环境','high',true,'必须先隔离验证并确认环境安全。',true),
('risk_configuration_no_restore','配置错误非恢复流程风险控制','MySQL','mysql_configuration_no_restore','实际不涉及数据恢复','low',false,'Agent直接退出恢复流程，转配置修复建议。',true),
('risk_unverified_inferred_strategy','未验证策略风险控制','MySQL','','未命中正式 recovery_strategy_template 且策略来源为 inferred 或 unverified','high',true,'先人工确认；必要时先异机验证；不得自动执行。',true);

SET FOREIGN_KEY_CHECKS = 1;

-- Dump completed