-- MySQL dump 10.13  Distrib 9.6.0, for macos26.2 (arm64)
--
-- Host: 127.0.0.1    Database: backup_rule_knowledge
-- ------------------------------------------------------
-- Server version	11.8.3-MariaDB












--
-- Current Database: "backup_rule_knowledge"
--

CREATE SCHEMA IF NOT EXISTS "backup_rule_knowledge";

SET search_path TO "backup_rule_knowledge", public;

--
-- Table structure for table "app_backup_capability"
--

DROP TABLE IF EXISTS "app_backup_capability";


CREATE TABLE "app_backup_capability" (
  "id" varchar(32) NOT NULL,
  "app_object_type" varchar(32) DEFAULT NULL,
  "capability" varchar(32) DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "app_backup_capability"
--

BEGIN;

INSERT INTO "app_backup_capability" VALUES ('ABC_01','mysql','完全备份','MySQL 支持完全备份能力。'),('ABC_02','mysql','增量备份','MySQL 支持增量备份能力。'),('ABC_03','mysql','日志备份','MySQL 支持日志备份能力。'),('ABC_04','oracle','完全备份','Oracle 支持完全备份能力。'),('ABC_05','oracle','增量备份','Oracle 支持增量备份能力。'),('ABC_06','oracle','日志备份','Oracle 支持日志备份能力。'),('ABC_07','file','完全备份','文件类应用支持完全备份能力。'),('ABC_08','file','增量备份','文件类应用支持增量备份能力。');

COMMIT;

--
-- Table structure for table "app_config_rule"
--

DROP TABLE IF EXISTS "app_config_rule";


CREATE TABLE "app_config_rule" (
  "id" varchar(32) NOT NULL,
  "app_object_type" varchar(32) DEFAULT NULL,
  "condition" text DEFAULT NULL,
  "condition_description" text DEFAULT NULL,
  "config_category" varchar(32) DEFAULT NULL,
  "config" text DEFAULT NULL,
  "config_param" varchar(32) DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "app_config_rule"
--

BEGIN;

INSERT INTO "app_config_rule" VALUES ('AC_001_01','mysql','target_rpo_minutes > 60','MySQL标准逻辑备份配置规则','app.mysql.backupmode','备份方式：XtraBackup','backup_type=XtraBackup'),('AC_001_02','mysql','target_rpo_minutes > 60','MySQL标准逻辑备份配置规则','app.mysql.realtimelog','实时日志备份：关闭','is_real_time_log=0'),('AC_001_03','mysql','target_rpo_minutes > 60','MySQL标准逻辑备份配置规则','app.mysql.channel','通道数开关：开启；通道数值：2','data_channel_num=2'),('AC_001_04','mysql','target_rpo_minutes > 60','MySQL标准逻辑备份配置规则','app.mysql.mergebackup','即时合成备份：关闭','is_merge_backup=0'),('AC_002_01','mysql','data_size_gb >= 500G','MySQL大数据量逻辑备份配置规则','app.mysql.backupmode','备份方式：XtraBackup','backup_type=XtraBackup'),('AC_002_02','mysql','data_size_gb >= 500G','MySQL大数据量逻辑备份配置规则','app.mysql.channel','通道数开关：开启；通道数值：4','data_channel_num=4'),('AC_002_03','mysql','data_size_gb >= 500G','MySQL大数据量逻辑备份配置规则','app.mysql.realtimelog','实时日志备份：关闭','is_real_time_log=0'),('AC_002_04','mysql','data_size_gb >= 500G','MySQL大数据量逻辑备份配置规则','app.mysql.mergebackup','即时合成备份：关闭','is_merge_backup=0'),('AC_003_01','mysql','target_rpo_minutes <= 60','MySQL低RPO物理备份配置规则','app.mysql.backupmode','备份方式：XtraBackup','backup_type=XtraBackup'),('AC_003_02','mysql','target_rpo_minutes <= 60','MySQL低RPO物理备份配置规则','app.mysql.realtimelog','实时日志备份：开启','is_real_time_log=1'),('AC_003_03','mysql','target_rpo_minutes <= 60','MySQL低RPO物理备份配置规则','app.mysql.channel','通道数开关：开启；通道数值：2','data_channel_num=2'),('AC_003_04','mysql','target_rpo_minutes <= 60','MySQL低RPO物理备份配置规则','app.mysql.mergebackup','即时合成备份：关闭','is_merge_backup=0');

COMMIT;

--
-- Table structure for table "app_output_rule"
--

DROP TABLE IF EXISTS "app_output_rule";


CREATE TABLE "app_output_rule" (
  "id" varchar(32) NOT NULL,
  "app_object_type" varchar(32) DEFAULT NULL,
  "profile_name" varchar(32) DEFAULT NULL,
  "solution_section_name" varchar(32) DEFAULT NULL,
  "strategy_section_name" varchar(32) DEFAULT NULL,
  "window_section_name" varchar(32) DEFAULT NULL,
  "plan_section_name" varchar(32) DEFAULT NULL,
  "retention_section_name" varchar(32) DEFAULT NULL,
  "basic_config_section_name" varchar(32) DEFAULT NULL,
  "region_config_section_name" varchar(32) DEFAULT NULL,
  "client_section_name" varchar(32) DEFAULT NULL,
  "storage_pool_section_name" varchar(32) DEFAULT NULL,
  "app_config_section_name" varchar(32) DEFAULT NULL,
  "hide_id_fields" smallint DEFAULT NULL,
  "empty_section_policy" varchar(32) DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "app_output_rule"
--

BEGIN;

INSERT INTO "app_output_rule" VALUES ('OUT_DEFAULT','all','通用输出方案','备份方案','备份策略','备份窗口','备份计划','副本保留','基础配置','备份区域配置','备份客户端','备份存储池','应用配置',1,'omit','适用于未定义专属输出规则时的通用展示结构。'),('OUT_FILE','file','文件输出方案','备份方案','备份策略','备份窗口','备份计划','副本保留','基础配置','备份区域配置','备份客户端','备份存储池','应用配置',1,'omit','适用于文件类推荐结果的分层展示结构，并与备份方案配置定义保持一致。'),('OUT_MYSQL','mysql','MySQL输出方案','备份方案','备份策略','备份窗口','备份计划','副本保留','基础配置','备份区域配置','备份客户端','备份存储池','应用配置',1,'omit','适用于 MySQL 推荐结果的分层展示结构，并与备份方案配置定义保持一致。'),('OUT_ORACLE','oracle','Oracle输出方案','备份方案','备份策略','备份窗口','备份计划','副本保留','基础配置','备份区域配置','备份客户端','备份存储池','应用配置',1,'omit','适用于 Oracle 推荐结果的分层展示结构，并与备份方案配置定义保持一致。');

COMMIT;

--
-- Table structure for table "app_rule"
--

DROP TABLE IF EXISTS "app_rule";


CREATE TABLE "app_rule" (
  "id" varchar(32) NOT NULL,
  "app_object_type" varchar(32) DEFAULT NULL,
  "data_type" varchar(32) DEFAULT NULL,
  "legal_requirement" varchar(32) DEFAULT NULL,
  "rule_scope" varchar(32) DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "app_rule"
--

BEGIN;

INSERT INTO "app_rule" VALUES ('AR_01','oracle','实时交易记录','《反洗钱法》','应用约束','适用于oracle的实时交易记录场景，法规要求：《反洗钱法》'),('AR_02','file','生物识别数据','《个人信息保护法》','应用约束','适用于file的生物识别数据场景，法规要求：《个人信息保护法》'),('AR_03','mysql','普通数据','通用要求','应用约束','适用于mysql的普通数据场景，法规要求：通用要求'),('AR_04','oracle','普通数据','通用要求','应用约束','适用于oracle的普通数据场景，法规要求：通用要求'),('AR_05','file','普通数据','通用要求','应用约束','适用于file的普通数据场景，法规要求：通用要求'),('AR_06','mysql','实时交易记录','《反洗钱法》','应用约束','适用于mysql的实时交易记录场景，法规要求：《反洗钱法》'),('AR_07','oracle','财务报告','《会计档案管理办法》','应用约束','适用于oracle的财务报告场景，法规要求：《会计档案管理办法》'),('AR_08','file','普通数据','《政务数据安全管理办法》','应用约束','适用于file的普通数据场景，法规要求：《政务数据安全管理办法》'),('AR_09','mysql','普通数据','《政务数据安全管理办法》','应用约束','适用于mysql的普通数据场景，法规要求：《政务数据安全管理办法》'),('AR_10','mysql','系统操作日志','《能源行业数据安全管理办法》','应用约束','适用于mysql的系统操作日志场景，法规要求：《能源行业数据安全管理办法》'),('AR_11','oracle','普通数据','《能源行业数据安全管理办法》','应用约束','适用于oracle的普通数据场景，法规要求：《能源行业数据安全管理办法》'),('AR_12','file','普通数据','《能源行业数据安全管理办法》','应用约束','适用于file的普通数据场景，法规要求：《能源行业数据安全管理办法》');

COMMIT;

--
-- Table structure for table "backup_plan"
--

DROP TABLE IF EXISTS "backup_plan";


CREATE TABLE "backup_plan" (
  "id" varchar(32) NOT NULL,
  "backup_solution_id" varchar(32) DEFAULT NULL,
  "name" varchar(32) DEFAULT NULL,
  "app_object_type" varchar(32) DEFAULT NULL,
  "plan_type" varchar(32) DEFAULT NULL,
  "period_value" integer DEFAULT NULL,
  "launch_time" varchar(32) DEFAULT NULL,
  "launch_minute" integer DEFAULT NULL,
  "launch_weekday" varchar(32) DEFAULT NULL,
  "launch_basis" varchar(32) DEFAULT NULL,
  "launch_date" integer DEFAULT NULL,
  "launch_week_order" varchar(32) DEFAULT NULL,
  "launch_weekday_in_period" varchar(32) DEFAULT NULL,
  "quarter_scope" varchar(32) DEFAULT NULL,
  "quarter_month" varchar(32) DEFAULT NULL,
  "annual_launch_datetime" varchar(32) DEFAULT NULL,
  "backup_method" varchar(32) DEFAULT NULL,
  "requires_log_backup" smallint DEFAULT NULL,
  "min_data_size_gb" integer DEFAULT NULL,
  "max_data_size_gb" integer DEFAULT NULL,
  "order_index" integer DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "backup_plan"
--

BEGIN;

INSERT INTO "backup_plan" VALUES ('PS_01_FULL','PS_01','MySQL普通数据-逻辑增量策略-完全备份','mysql','周级计划',1,'01:00',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于MySQL普通数据场景的默认推荐策略。'),('PS_01_INCR','PS_01','MySQL普通数据-逻辑增量策略-增量备份','mysql','天级计划',1,'01:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',0,NULL,NULL,2,'适用于MySQL普通数据场景的默认推荐策略。'),('PS_02_FULL','PS_02','MySQL实时交易-物理实时日志策略-完全备份','mysql','周级计划',1,'00:30',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',1,NULL,NULL,1,'适用于要求分钟级RPO的MySQL实时交易业务。'),('PS_02_INCR','PS_02','MySQL实时交易-物理实时日志策略-增量备份','mysql','天级计划',1,'00:30',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',1,NULL,NULL,2,'适用于要求分钟级RPO的MySQL实时交易业务。'),('PS_02_LOG','PS_02','MySQL实时交易-物理实时日志策略-日志备份','mysql','分钟级计划',1,'00:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'日志备份',1,NULL,NULL,3,'适用于要求分钟级RPO的MySQL实时交易业务。'),('PS_03_FULL','PS_03','MySQL重要数据-逻辑增量策略-完全备份','mysql','周级计划',1,'01:30',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于重要级MySQL数据的标准增量备份策略。'),('PS_03_INCR','PS_03','MySQL重要数据-逻辑增量策略-增量备份','mysql','天级计划',1,'01:30',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',0,NULL,NULL,2,'适用于重要级MySQL数据的标准增量备份策略。'),('PS_04_FULL','PS_04','MySQL敏感数据-每日全备策略-完全备份','mysql','天级计划',1,'02:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于敏感MySQL数据的每日全量备份策略。'),('PS_05_FULL','PS_05','Oracle实时交易-物理实时日志策略-完全备份','oracle','周级计划',1,'00:30',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',1,NULL,NULL,1,'适用于Oracle实时交易业务的低RPO策略。'),('PS_05_INCR','PS_05','Oracle实时交易-物理实时日志策略-增量备份','oracle','天级计划',1,'00:30',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',1,NULL,NULL,2,'适用于Oracle实时交易业务的低RPO策略。'),('PS_05_LOG','PS_05','Oracle实时交易-物理实时日志策略-日志备份','oracle','分钟级计划',1,'00:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'日志备份',1,NULL,NULL,3,'适用于Oracle实时交易业务的低RPO策略。'),('PS_06_FULL','PS_06','Oracle财务核心-物理增强策略-完全备份','oracle','周级计划',1,'01:00',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',1,NULL,NULL,1,'适用于Oracle财务和核心数据库的标准增强策略。'),('PS_06_INCR','PS_06','Oracle财务核心-物理增强策略-增量备份','oracle','天级计划',1,'01:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',1,NULL,NULL,2,'适用于Oracle财务和核心数据库的标准增强策略。'),('PS_06_LOG','PS_06','Oracle财务核心-物理增强策略-日志备份','oracle','分钟级计划',30,'00:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'日志备份',1,NULL,NULL,3,'适用于Oracle财务和核心数据库的标准增强策略。'),('PS_07_FULL','PS_07','文件普通数据-标准文件备份策略-完全备份','file','周级计划',1,'23:30',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'permanent_incremental',0,NULL,NULL,1,'适用于文件普通数据的标准文件备份策略。'),('PS_07_INCR','PS_07','文件普通数据-标准文件备份策略-增量备份','file','天级计划',1,'23:30',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'差异备份',0,NULL,NULL,2,'适用于文件普通数据的标准文件备份策略。'),('PS_08_FULL','PS_08','文件敏感数据-加密文件备份策略-完全备份','file','天级计划',1,'23:30',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'permanent_incremental',0,NULL,NULL,1,'适用于敏感文件数据的加密备份策略。'),('PS_09_FULL','PS_09','Oracle普通数据-物理标准策略-完全备份','oracle','周级计划',1,'01:00',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于Oracle普通业务的标准物理备份策略。'),('PS_09_INCR','PS_09','Oracle普通数据-物理标准策略-增量备份','oracle','天级计划',1,'01:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',0,NULL,NULL,2,'适用于Oracle普通业务的标准物理备份策略。'),('PS_10_FULL','PS_10','MySQL普通数据-每日全备稳妥策略-完全备份','mysql','天级计划',1,'02:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于普通MySQL业务的恢复优先型候选策略。'),('PS_11_FULL','PS_11','MySQL普通数据-高频增量增强策略-完全备份','mysql','周级计划',1,'01:00',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于需要更快恢复和更低RPO的MySQL普通业务候选策略。'),('PS_11_INCR','PS_11','MySQL普通数据-高频增量增强策略-增量备份','mysql','天级计划',1,'04:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',0,NULL,NULL,2,'适用于需要更快恢复和更低RPO的MySQL普通业务候选策略。'),('PS_12_FULL','PS_12','Oracle普通数据-高频增量增强策略-完全备份','oracle','周级计划',1,'01:00',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于Oracle普通业务的增强恢复候选策略。'),('PS_12_INCR','PS_12','Oracle普通数据-高频增量增强策略-增量备份','oracle','天级计划',1,'04:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',0,NULL,NULL,2,'适用于Oracle普通业务的增强恢复候选策略。'),('PS_13_FULL','PS_13','文件普通数据-政务归档增强策略-完全备份','file','周级计划',1,'23:00',NULL,'周日',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',0,NULL,NULL,1,'适用于政务归档类文件数据的长期留存候选策略。'),('PS_13_INCR','PS_13','文件普通数据-政务归档增强策略-增量备份','file','天级计划',1,'23:30',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'差异备份',0,NULL,NULL,2,'适用于政务归档类文件数据的长期留存候选策略。'),('PS_14_FULL','PS_14','MySQL系统日志-能源高频增量策略-完全备份','mysql','天级计划',1,'00:30',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'完全备份',1,NULL,NULL,1,'适用于能源行业系统日志和生产明细的高频候选策略。'),('PS_14_INCR','PS_14','MySQL系统日志-能源高频增量策略-增量备份','mysql','天级计划',1,'02:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'增量备份',1,NULL,NULL,2,'适用于能源行业系统日志和生产明细的高频候选策略。'),('PS_14_LOG','PS_14','MySQL系统日志-能源高频增量策略-日志备份','mysql','分钟级计划',15,'00:00',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'日志备份',1,NULL,NULL,3,'适用于能源行业系统日志和生产明细的高频候选策略。');

COMMIT;

--
-- Table structure for table "backup_solution"
--

DROP TABLE IF EXISTS "backup_solution";


CREATE TABLE "backup_solution" (
  "id" varchar(32) NOT NULL,
  "name" varchar(32) DEFAULT NULL,
  "app_object_type" varchar(32) DEFAULT NULL,
  "data_level" varchar(32) DEFAULT NULL,
  "industry_type" varchar(32) DEFAULT NULL,
  "target_rpo_minutes" integer DEFAULT NULL,
  "target_rto_minutes" integer DEFAULT NULL,
  "source_app_rule_id" varchar(32) DEFAULT NULL,
  "source_data_type" varchar(32) DEFAULT NULL,
  "source_legal_rule_id" varchar(32) DEFAULT NULL,
  "source_rpo_rule_id" varchar(32) DEFAULT NULL,
  "description" text DEFAULT NULL,
  "priority" integer DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "backup_solution"
--

BEGIN;

INSERT INTO "backup_solution" VALUES ('PS_01','MySQL普通数据-逻辑增量策略','mysql','普通','未定义',1440,120,'AR_03','普通数据','LR_00','RPO_03','适用于MySQL普通数据场景的默认推荐策略。',50),('PS_02','MySQL实时交易-物理实时日志策略','mysql','核心','金融',5,60,'AR_06','实时交易记录','LR_01','RPO_00','适用于要求分钟级RPO的MySQL实时交易业务。',10),('PS_03','MySQL重要数据-逻辑增量策略','mysql','重要','未定义',60,120,NULL,NULL,NULL,'RPO_02','适用于重要级MySQL数据的标准增量备份策略。',30),('PS_04','MySQL敏感数据-每日全备策略','mysql','敏感','未定义',1440,120,NULL,'个人身份信息','LR_03','RPO_03','适用于敏感MySQL数据的每日全量备份策略。',40),('PS_05','Oracle实时交易-物理实时日志策略','oracle','核心','金融',5,60,'AR_01','实时交易记录','LR_01','RPO_00','适用于Oracle实时交易业务的低RPO策略。',5),('PS_06','Oracle财务核心-物理增强策略','oracle','核心','金融',60,120,'AR_07','财务报告','LR_02','RPO_02','适用于Oracle财务和核心数据库的标准增强策略。',20),('PS_07','文件普通数据-标准文件备份策略','file','普通','未定义',1440,240,'AR_05','普通数据','LR_00','RPO_03','适用于文件普通数据的标准文件备份策略。',60),('PS_08','文件敏感数据-加密文件备份策略','file','敏感','未定义',1440,120,'AR_02','生物识别数据','LR_03','RPO_03','适用于敏感文件数据的加密备份策略。',35),('PS_09','Oracle普通数据-物理标准策略','oracle','普通','未定义',1440,120,'AR_04','普通数据','LR_00','RPO_03','适用于Oracle普通业务的标准物理备份策略。',45),('PS_10','MySQL普通数据-每日全备稳妥策略','mysql','普通','未定义',1440,90,'AR_03','普通数据','LR_00','RPO_03','适用于普通MySQL业务的恢复优先型候选策略。',52),('PS_11','MySQL普通数据-高频增量增强策略','mysql','普通','未定义',240,60,'AR_03','普通数据','LR_00','RPO_02','适用于需要更快恢复和更低RPO的MySQL普通业务候选策略。',28),('PS_12','Oracle普通数据-高频增量增强策略','oracle','普通','未定义',240,60,'AR_04','普通数据','LR_00','RPO_02','适用于Oracle普通业务的增强恢复候选策略。',26),('PS_13','文件普通数据-政务归档增强策略','file','普通','政府',1440,120,'AR_08','普通数据','LR_04','RPO_03','适用于政务归档类文件数据的长期留存候选策略。',34),('PS_14','MySQL系统日志-能源高频增量策略','mysql','重要','能源',60,90,'AR_10','系统操作日志','LR_05','RPO_02','适用于能源行业系统日志和生产明细的高频候选策略。',22);

COMMIT;

--
-- Table structure for table "backup_window"
--

DROP TABLE IF EXISTS "backup_window";


CREATE TABLE "backup_window" (
  "id" varchar(32) NOT NULL,
  "backup_solution_id" varchar(32) DEFAULT NULL,
  "name" varchar(32) DEFAULT NULL,
  "window_type" varchar(32) DEFAULT NULL,
  "start_weekday" varchar(32) DEFAULT NULL,
  "end_weekday" varchar(32) DEFAULT NULL,
  "start_time" varchar(32) DEFAULT NULL,
  "end_time" varchar(32) DEFAULT NULL,
  "max_speed_mib" integer DEFAULT NULL,
  "priority" integer DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "backup_window"
--

BEGIN;

INSERT INTO "backup_window" VALUES ('PS_01_friday_allowed','PS_01','PS_01-friday-allowed','限速备份','周五','周五','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_friday_forbidden_after','PS_01','PS_01-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_monday_allowed','PS_01','PS_01-monday-allowed','限速备份','周一','周一','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_monday_forbidden_after','PS_01','PS_01-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_saturday_allowed','PS_01','PS_01-saturday-allowed','限速备份','周六','周六','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_saturday_forbidden_after','PS_01','PS_01-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_sunday_allowed','PS_01','PS_01-sunday-allowed','限速备份','周日','周日','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_sunday_forbidden_after','PS_01','PS_01-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_thursday_allowed','PS_01','PS_01-thursday-allowed','限速备份','周四','周四','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_thursday_forbidden_after','PS_01','PS_01-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_tuesday_allowed','PS_01','PS_01-tuesday-allowed','限速备份','周二','周二','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_tuesday_forbidden_after','PS_01','PS_01-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_wednesday_allowed','PS_01','PS_01-wednesday-allowed','限速备份','周三','周三','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_01_wednesday_forbidden_after','PS_01','PS_01-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_02_friday_allowed','PS_02','PS_02-friday-allowed','限速备份','周五','周五','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_friday_forbidden_after','PS_02','PS_02-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_monday_allowed','PS_02','PS_02-monday-allowed','限速备份','周一','周一','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_monday_forbidden_after','PS_02','PS_02-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_saturday_allowed','PS_02','PS_02-saturday-allowed','限速备份','周六','周六','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_saturday_forbidden_after','PS_02','PS_02-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_sunday_allowed','PS_02','PS_02-sunday-allowed','限速备份','周日','周日','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_sunday_forbidden_after','PS_02','PS_02-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_thursday_allowed','PS_02','PS_02-thursday-allowed','限速备份','周四','周四','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_thursday_forbidden_after','PS_02','PS_02-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_tuesday_allowed','PS_02','PS_02-tuesday-allowed','限速备份','周二','周二','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_tuesday_forbidden_after','PS_02','PS_02-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_wednesday_allowed','PS_02','PS_02-wednesday-allowed','限速备份','周三','周三','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_02_wednesday_forbidden_after','PS_02','PS_02-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_03_friday_allowed','PS_03','PS_03-friday-allowed','限速备份','周五','周五','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_friday_forbidden_after','PS_03','PS_03-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_monday_allowed','PS_03','PS_03-monday-allowed','限速备份','周一','周一','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_monday_forbidden_after','PS_03','PS_03-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_saturday_allowed','PS_03','PS_03-saturday-allowed','限速备份','周六','周六','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_saturday_forbidden_after','PS_03','PS_03-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_sunday_allowed','PS_03','PS_03-sunday-allowed','限速备份','周日','周日','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_sunday_forbidden_after','PS_03','PS_03-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_thursday_allowed','PS_03','PS_03-thursday-allowed','限速备份','周四','周四','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_thursday_forbidden_after','PS_03','PS_03-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_tuesday_allowed','PS_03','PS_03-tuesday-allowed','限速备份','周二','周二','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_tuesday_forbidden_after','PS_03','PS_03-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_wednesday_allowed','PS_03','PS_03-wednesday-allowed','限速备份','周三','周三','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_03_wednesday_forbidden_after','PS_03','PS_03-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_04_friday_allowed','PS_04','PS_04-friday-allowed','限速备份','周五','周五','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_friday_forbidden_after','PS_04','PS_04-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_monday_allowed','PS_04','PS_04-monday-allowed','限速备份','周一','周一','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_monday_forbidden_after','PS_04','PS_04-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_saturday_allowed','PS_04','PS_04-saturday-allowed','限速备份','周六','周六','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_saturday_forbidden_after','PS_04','PS_04-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_sunday_allowed','PS_04','PS_04-sunday-allowed','限速备份','周日','周日','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_sunday_forbidden_after','PS_04','PS_04-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_thursday_allowed','PS_04','PS_04-thursday-allowed','限速备份','周四','周四','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_thursday_forbidden_after','PS_04','PS_04-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_tuesday_allowed','PS_04','PS_04-tuesday-allowed','限速备份','周二','周二','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_tuesday_forbidden_after','PS_04','PS_04-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_wednesday_allowed','PS_04','PS_04-wednesday-allowed','限速备份','周三','周三','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_04_wednesday_forbidden_after','PS_04','PS_04-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_04'),('PS_05_friday_allowed','PS_05','PS_05-friday-allowed','限速备份','周五','周五','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_friday_forbidden_after','PS_05','PS_05-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_monday_allowed','PS_05','PS_05-monday-allowed','限速备份','周一','周一','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_monday_forbidden_after','PS_05','PS_05-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_saturday_allowed','PS_05','PS_05-saturday-allowed','限速备份','周六','周六','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_saturday_forbidden_after','PS_05','PS_05-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_sunday_allowed','PS_05','PS_05-sunday-allowed','限速备份','周日','周日','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_sunday_forbidden_after','PS_05','PS_05-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_thursday_allowed','PS_05','PS_05-thursday-allowed','限速备份','周四','周四','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_thursday_forbidden_after','PS_05','PS_05-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_tuesday_allowed','PS_05','PS_05-tuesday-allowed','限速备份','周二','周二','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_tuesday_forbidden_after','PS_05','PS_05-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_wednesday_allowed','PS_05','PS_05-wednesday-allowed','限速备份','周三','周三','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_05_wednesday_forbidden_after','PS_05','PS_05-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_06_friday_allowed','PS_06','PS_06-friday-allowed','限速备份','周五','周五','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_friday_forbidden_after','PS_06','PS_06-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_monday_allowed','PS_06','PS_06-monday-allowed','限速备份','周一','周一','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_monday_forbidden_after','PS_06','PS_06-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_saturday_allowed','PS_06','PS_06-saturday-allowed','限速备份','周六','周六','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_saturday_forbidden_after','PS_06','PS_06-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_sunday_allowed','PS_06','PS_06-sunday-allowed','限速备份','周日','周日','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_sunday_forbidden_after','PS_06','PS_06-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_thursday_allowed','PS_06','PS_06-thursday-allowed','限速备份','周四','周四','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_thursday_forbidden_after','PS_06','PS_06-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_tuesday_allowed','PS_06','PS_06-tuesday-allowed','限速备份','周二','周二','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_tuesday_forbidden_after','PS_06','PS_06-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_wednesday_allowed','PS_06','PS_06-wednesday-allowed','限速备份','周三','周三','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_06_wednesday_forbidden_after','PS_06','PS_06-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_01'),('PS_07_friday_allowed_early','PS_07','PS_07-friday-allowed_early','限速备份','周五','周五','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_friday_allowed_late','PS_07','PS_07-friday-allowed_late','限速备份','周五','周五','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_friday_forbidden_mid','PS_07','PS_07-friday-forbidden_mid','禁止备份','周五','周五','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_monday_allowed_early','PS_07','PS_07-monday-allowed_early','限速备份','周一','周一','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_monday_allowed_late','PS_07','PS_07-monday-allowed_late','限速备份','周一','周一','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_monday_forbidden_mid','PS_07','PS_07-monday-forbidden_mid','禁止备份','周一','周一','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_saturday_allowed_early','PS_07','PS_07-saturday-allowed_early','限速备份','周六','周六','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_saturday_allowed_late','PS_07','PS_07-saturday-allowed_late','限速备份','周六','周六','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_saturday_forbidden_mid','PS_07','PS_07-saturday-forbidden_mid','禁止备份','周六','周六','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_sunday_allowed_early','PS_07','PS_07-sunday-allowed_early','限速备份','周日','周日','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_sunday_allowed_late','PS_07','PS_07-sunday-allowed_late','限速备份','周日','周日','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_sunday_forbidden_mid','PS_07','PS_07-sunday-forbidden_mid','禁止备份','周日','周日','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_thursday_allowed_early','PS_07','PS_07-thursday-allowed_early','限速备份','周四','周四','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_thursday_allowed_late','PS_07','PS_07-thursday-allowed_late','限速备份','周四','周四','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_thursday_forbidden_mid','PS_07','PS_07-thursday-forbidden_mid','禁止备份','周四','周四','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_tuesday_allowed_early','PS_07','PS_07-tuesday-allowed_early','限速备份','周二','周二','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_tuesday_allowed_late','PS_07','PS_07-tuesday-allowed_late','限速备份','周二','周二','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_tuesday_forbidden_mid','PS_07','PS_07-tuesday-forbidden_mid','禁止备份','周二','周二','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_wednesday_allowed_early','PS_07','PS_07-wednesday-allowed_early','限速备份','周三','周三','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_wednesday_allowed_late','PS_07','PS_07-wednesday-allowed_late','限速备份','周三','周三','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_07_wednesday_forbidden_mid','PS_07','PS_07-wednesday-forbidden_mid','禁止备份','周三','周三','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_08_friday_allowed_early','PS_08','PS_08-friday-allowed_early','限速备份','周五','周五','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_friday_allowed_late','PS_08','PS_08-friday-allowed_late','限速备份','周五','周五','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_friday_forbidden_mid','PS_08','PS_08-friday-forbidden_mid','禁止备份','周五','周五','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_monday_allowed_early','PS_08','PS_08-monday-allowed_early','限速备份','周一','周一','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_monday_allowed_late','PS_08','PS_08-monday-allowed_late','限速备份','周一','周一','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_monday_forbidden_mid','PS_08','PS_08-monday-forbidden_mid','禁止备份','周一','周一','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_saturday_allowed_early','PS_08','PS_08-saturday-allowed_early','限速备份','周六','周六','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_saturday_allowed_late','PS_08','PS_08-saturday-allowed_late','限速备份','周六','周六','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_saturday_forbidden_mid','PS_08','PS_08-saturday-forbidden_mid','禁止备份','周六','周六','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_sunday_allowed_early','PS_08','PS_08-sunday-allowed_early','限速备份','周日','周日','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_sunday_allowed_late','PS_08','PS_08-sunday-allowed_late','限速备份','周日','周日','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_sunday_forbidden_mid','PS_08','PS_08-sunday-forbidden_mid','禁止备份','周日','周日','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_thursday_allowed_early','PS_08','PS_08-thursday-allowed_early','限速备份','周四','周四','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_thursday_allowed_late','PS_08','PS_08-thursday-allowed_late','限速备份','周四','周四','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_thursday_forbidden_mid','PS_08','PS_08-thursday-forbidden_mid','禁止备份','周四','周四','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_tuesday_allowed_early','PS_08','PS_08-tuesday-allowed_early','限速备份','周二','周二','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_tuesday_allowed_late','PS_08','PS_08-tuesday-allowed_late','限速备份','周二','周二','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_tuesday_forbidden_mid','PS_08','PS_08-tuesday-forbidden_mid','禁止备份','周二','周二','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_wednesday_allowed_early','PS_08','PS_08-wednesday-allowed_early','限速备份','周三','周三','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_wednesday_allowed_late','PS_08','PS_08-wednesday-allowed_late','限速备份','周三','周三','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_08_wednesday_forbidden_mid','PS_08','PS_08-wednesday-forbidden_mid','禁止备份','周三','周三','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_03'),('PS_09_friday_allowed','PS_09','PS_09-friday-allowed','限速备份','周五','周五','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_friday_forbidden_after','PS_09','PS_09-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_monday_allowed','PS_09','PS_09-monday-allowed','限速备份','周一','周一','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_monday_forbidden_after','PS_09','PS_09-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_saturday_allowed','PS_09','PS_09-saturday-allowed','限速备份','周六','周六','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_saturday_forbidden_after','PS_09','PS_09-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_sunday_allowed','PS_09','PS_09-sunday-allowed','限速备份','周日','周日','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_sunday_forbidden_after','PS_09','PS_09-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_thursday_allowed','PS_09','PS_09-thursday-allowed','限速备份','周四','周四','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_thursday_forbidden_after','PS_09','PS_09-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_tuesday_allowed','PS_09','PS_09-tuesday-allowed','限速备份','周二','周二','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_tuesday_forbidden_after','PS_09','PS_09-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_wednesday_allowed','PS_09','PS_09-wednesday-allowed','限速备份','周三','周三','00:00','06:00',150,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_09_wednesday_forbidden_after','PS_09','PS_09-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_05'),('PS_10_friday_allowed','PS_10','PS_10-friday-allowed','限速备份','周五','周五','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_friday_forbidden_after','PS_10','PS_10-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_monday_allowed','PS_10','PS_10-monday-allowed','限速备份','周一','周一','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_monday_forbidden_after','PS_10','PS_10-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_saturday_allowed','PS_10','PS_10-saturday-allowed','限速备份','周六','周六','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_saturday_forbidden_after','PS_10','PS_10-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_sunday_allowed','PS_10','PS_10-sunday-allowed','限速备份','周日','周日','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_sunday_forbidden_after','PS_10','PS_10-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_thursday_allowed','PS_10','PS_10-thursday-allowed','限速备份','周四','周四','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_thursday_forbidden_after','PS_10','PS_10-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_tuesday_allowed','PS_10','PS_10-tuesday-allowed','限速备份','周二','周二','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_tuesday_forbidden_after','PS_10','PS_10-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_wednesday_allowed','PS_10','PS_10-wednesday-allowed','限速备份','周三','周三','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_10_wednesday_forbidden_after','PS_10','PS_10-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_friday_allowed','PS_11','PS_11-friday-allowed','限速备份','周五','周五','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_friday_forbidden_after','PS_11','PS_11-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_monday_allowed','PS_11','PS_11-monday-allowed','限速备份','周一','周一','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_monday_forbidden_after','PS_11','PS_11-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_saturday_allowed','PS_11','PS_11-saturday-allowed','限速备份','周六','周六','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_saturday_forbidden_after','PS_11','PS_11-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_sunday_allowed','PS_11','PS_11-sunday-allowed','限速备份','周日','周日','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_sunday_forbidden_after','PS_11','PS_11-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_thursday_allowed','PS_11','PS_11-thursday-allowed','限速备份','周四','周四','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_thursday_forbidden_after','PS_11','PS_11-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_tuesday_allowed','PS_11','PS_11-tuesday-allowed','限速备份','周二','周二','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_tuesday_forbidden_after','PS_11','PS_11-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_wednesday_allowed','PS_11','PS_11-wednesday-allowed','限速备份','周三','周三','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_11_wednesday_forbidden_after','PS_11','PS_11-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_friday_allowed','PS_12','PS_12-friday-allowed','限速备份','周五','周五','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_friday_forbidden_after','PS_12','PS_12-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_monday_allowed','PS_12','PS_12-monday-allowed','限速备份','周一','周一','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_monday_forbidden_after','PS_12','PS_12-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_saturday_allowed','PS_12','PS_12-saturday-allowed','限速备份','周六','周六','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_saturday_forbidden_after','PS_12','PS_12-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_sunday_allowed','PS_12','PS_12-sunday-allowed','限速备份','周日','周日','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_sunday_forbidden_after','PS_12','PS_12-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_thursday_allowed','PS_12','PS_12-thursday-allowed','限速备份','周四','周四','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_thursday_forbidden_after','PS_12','PS_12-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_tuesday_allowed','PS_12','PS_12-tuesday-allowed','限速备份','周二','周二','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_tuesday_forbidden_after','PS_12','PS_12-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_wednesday_allowed','PS_12','PS_12-wednesday-allowed','限速备份','周三','周三','00:00','06:00',300,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_12_wednesday_forbidden_after','PS_12','PS_12-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_06'),('PS_13_friday_allowed_early','PS_13','PS_13-friday-allowed_early','限速备份','周五','周五','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_friday_allowed_late','PS_13','PS_13-friday-allowed_late','限速备份','周五','周五','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_friday_forbidden_mid','PS_13','PS_13-friday-forbidden_mid','禁止备份','周五','周五','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_monday_allowed_early','PS_13','PS_13-monday-allowed_early','限速备份','周一','周一','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_monday_allowed_late','PS_13','PS_13-monday-allowed_late','限速备份','周一','周一','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_monday_forbidden_mid','PS_13','PS_13-monday-forbidden_mid','禁止备份','周一','周一','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_saturday_allowed_early','PS_13','PS_13-saturday-allowed_early','限速备份','周六','周六','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_saturday_allowed_late','PS_13','PS_13-saturday-allowed_late','限速备份','周六','周六','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_saturday_forbidden_mid','PS_13','PS_13-saturday-forbidden_mid','禁止备份','周六','周六','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_sunday_allowed_early','PS_13','PS_13-sunday-allowed_early','限速备份','周日','周日','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_sunday_allowed_late','PS_13','PS_13-sunday-allowed_late','限速备份','周日','周日','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_sunday_forbidden_mid','PS_13','PS_13-sunday-forbidden_mid','禁止备份','周日','周日','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_thursday_allowed_early','PS_13','PS_13-thursday-allowed_early','限速备份','周四','周四','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_thursday_allowed_late','PS_13','PS_13-thursday-allowed_late','限速备份','周四','周四','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_thursday_forbidden_mid','PS_13','PS_13-thursday-forbidden_mid','禁止备份','周四','周四','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_tuesday_allowed_early','PS_13','PS_13-tuesday-allowed_early','限速备份','周二','周二','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_tuesday_allowed_late','PS_13','PS_13-tuesday-allowed_late','限速备份','周二','周二','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_tuesday_forbidden_mid','PS_13','PS_13-tuesday-forbidden_mid','禁止备份','周二','周二','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_wednesday_allowed_early','PS_13','PS_13-wednesday-allowed_early','限速备份','周三','周三','00:00','06:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_wednesday_allowed_late','PS_13','PS_13-wednesday-allowed_late','限速备份','周三','周三','23:00','24:00',500,2,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_13_wednesday_forbidden_mid','PS_13','PS_13-wednesday-forbidden_mid','禁止备份','周三','周三','06:00','23:00',NULL,1,'由旧窗口 23:00-06:00 转换；来源基础配置 CC_07'),('PS_14_friday_allowed','PS_14','PS_14-friday-allowed','限速备份','周五','周五','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_friday_forbidden_after','PS_14','PS_14-friday-forbidden_after','禁止备份','周五','周五','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_monday_allowed','PS_14','PS_14-monday-allowed','限速备份','周一','周一','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_monday_forbidden_after','PS_14','PS_14-monday-forbidden_after','禁止备份','周一','周一','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_saturday_allowed','PS_14','PS_14-saturday-allowed','限速备份','周六','周六','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_saturday_forbidden_after','PS_14','PS_14-saturday-forbidden_after','禁止备份','周六','周六','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_sunday_allowed','PS_14','PS_14-sunday-allowed','限速备份','周日','周日','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_sunday_forbidden_after','PS_14','PS_14-sunday-forbidden_after','禁止备份','周日','周日','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_thursday_allowed','PS_14','PS_14-thursday-allowed','限速备份','周四','周四','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_thursday_forbidden_after','PS_14','PS_14-thursday-forbidden_after','禁止备份','周四','周四','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_tuesday_allowed','PS_14','PS_14-tuesday-allowed','限速备份','周二','周二','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_tuesday_forbidden_after','PS_14','PS_14-tuesday-forbidden_after','禁止备份','周二','周二','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_wednesday_allowed','PS_14','PS_14-wednesday-allowed','限速备份','周三','周三','00:00','06:00',200,2,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02'),('PS_14_wednesday_forbidden_after','PS_14','PS_14-wednesday-forbidden_after','禁止备份','周三','周三','06:00','24:00',NULL,1,'由旧窗口 00:00-06:00 转换；来源基础配置 CC_02');

COMMIT;

--
-- Table structure for table "basic_config_rule"
--

DROP TABLE IF EXISTS "basic_config_rule";


CREATE TABLE "basic_config_rule" (
  "id" varchar(32) NOT NULL,
  "app_object_type" varchar(32) DEFAULT NULL,
  "condition" text DEFAULT NULL,
  "condition_description" text DEFAULT NULL,
  "config_category" varchar(32) DEFAULT NULL,
  "config" text DEFAULT NULL,
  "config_param" varchar(128) DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "basic_config_rule"
--

BEGIN;

INSERT INTO "basic_config_rule" VALUES ('BC_001_01','all','default','默认基础配置规则','basic.backuptarget','备份存储：备份存储服务自动选择；备份存储池自动选择','storage_service_id=auto_select;storage_pool_id=auto_select'),('BC_001_03','all','default','默认基础配置规则','basic.transportencryption','传输加密：开启','encryption_trans=1'),('BC_002_01','mysql,oracle','default','数据库基础配置默认规则','basic.retry','备份失败重试：开启；失败重试次数 2；失败重试等待时间 5 分钟','failure_retry=1;failure_retry_count=2;failure_retry_interval=5'),('BC_003_01','mysql,oracle,file','target_rpo_minutes <= 60','低RPO场景基础配置增强规则','basic.retry','备份失败重试：开启；失败重试次数 3；失败重试等待时间 3 分钟','failure_retry=1;failure_retry_count=3;failure_retry_interval=3'),('BC_003_02','mysql,oracle,file','target_rpo_minutes <= 60','低RPO场景基础配置增强规则','basic.forcedretention','强制数据保留：开启；强制数据保留天数：1 天','forced_retention_switch=1;forced_retention_cycle=1'),('BC_004_01','file','default','文件系统基础配置默认规则','basic.compress','数据压缩：开启；数据压缩工作线程数 4；数据压缩算法：标准；数据压缩位置：源端','compress=1;compress_tread_num=4;compress_algorithm=2;compress_location=1'),('BC_004_02','file','default','文件系统基础配置默认规则','basic.deduplication','重复数据删除：开启；重复数据删除位置：目标端；重复数据删除工作线程数：4；指纹库模式：自动配置指纹库；指纹库：自动配置','deduplication=1;deduplication_location=2;deduplication_tread_num=4;finger_library_id=auto_config'),('BC_005_01','mysql,oracle,file','data_size_gb >= 500G','生产资源数据大于等于500G','basic.compress','数据压缩：开启；数据压缩工作线程数 4；数据压缩算法：标准；数据压缩位置：源端','compress=1;compress_tread_num=4;compress_algorithm=2;compress_location=1'),('BC_005_02','mysql,oracle,file','data_size_gb >= 500G','生产资源数据大于等于500G','basic.deduplication','重复数据删除：开启；重复数据删除位置：源端；重复数据删除工作线程数：4；指纹库模式：自动配置指纹库；指纹库：自动配置','deduplication=1;deduplication_location=1;deduplication_tread_num=4;finger_library_id=auto_config'),('BC_006_01','mysql,oracle','data_level = 普通','普通数据库轻量基础配置规则','basic.compress','数据压缩：关闭','compress=0'),('BC_006_02','mysql,oracle','data_level = 普通','普通数据库轻量基础配置规则','basic.deduplication','重复数据删除：关闭','deduplication=0'),('BC_006_03','mysql,oracle','data_level = 普通','普通数据库轻量基础配置规则','basic.transportencryption','传输加密：开启','encryption_trans=1'),('BC_007_01','mysql,oracle,file','data_level in [敏感,核心] or industry_type in [金融,政府,能源]','敏感数据或高合规场景基础配置规则','basic.transportencryption','传输加密：开启','encryption_trans=1'),('BC_007_03','mysql,oracle,file','data_level in [敏感,核心] or industry_type in [金融,政府,能源]','敏感数据或高合规场景基础配置规则','basic.retry','备份失败重试：开启；失败重试次数：2；失败重试等待时间：1 天','failure_retry=1;failure_retry_count=2;failure_retry_interval=1440'),('BC_008_01','mysql,oracle,file','backup_window in [23:00-06:00,00:00-06:00]','夜间备份窗口基础配置规则','basic.backuptarget','备份存储：备份存储服务自动选择；备份存储池自动选择','storage_service_id=auto_select;storage_pool_id=auto_select'),('BC_008_02','mysql,oracle,file','backup_window in [23:00-06:00,00:00-06:00]','夜间备份窗口基础配置规则','basic.retry','备份失败重试：开启；失败重试次数：2；失败重试等待时间：5 分钟','failure_retry=1;failure_retry_count=2;failure_retry_interval=5');

COMMIT;

--
-- Table structure for table "data_rule"
--

DROP TABLE IF EXISTS "data_rule";


CREATE TABLE "data_rule" (
  "data_type" varchar(32) DEFAULT NULL,
  "data_level" varchar(32) DEFAULT NULL,
  "data_content" text DEFAULT NULL,
  "rule_scope" varchar(32) DEFAULT NULL,
  "description" text DEFAULT NULL
) ;


--
-- Dumping data for table "data_rule"
--

BEGIN;

INSERT INTO "data_rule" VALUES ('个人身份信息','敏感','姓名、身份证号','敏感数据增强保护','典型内容：姓名、身份证号'),('实时交易记录','核心','大额跨境转账记录、支付流水','核心数据保护','典型内容：大额跨境转账记录、支付流水'),('普通数据','普通','业务普通数据、一般配置、非敏感附件','普通数据保护','典型内容：业务普通数据、一般配置、非敏感附件'),('生物识别数据','核心/敏感','指纹、人脸','核心数据保护','典型内容：指纹、人脸'),('系统操作日志','核心/重要','敏感级及以上数据的操作日志','核心数据保护','典型内容：敏感级及以上数据的操作日志'),('财务报告','核心','资产负债表、利润表','核心数据保护','典型内容：资产负债表、利润表');

COMMIT;

--
-- Table structure for table "legal_rule"
--

DROP TABLE IF EXISTS "legal_rule";


CREATE TABLE "legal_rule" (
  "id" varchar(32) NOT NULL,
  "name" varchar(32) DEFAULT NULL,
  "industry_type" varchar(32) DEFAULT NULL,
  "legal_requirement" varchar(32) DEFAULT NULL,
  "data_type" varchar(32) DEFAULT NULL,
  "data_level" varchar(32) DEFAULT NULL,
  "max_rpo_minutes" integer DEFAULT NULL,
  "max_rto_minutes" integer DEFAULT NULL,
  "min_retention_days" integer DEFAULT NULL,
  "compliance_scope" varchar(32) DEFAULT NULL,
  "severity" integer DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "legal_rule"
--

BEGIN;

INSERT INTO "legal_rule" VALUES ('LR_00','通用基线约束','未定义','通用要求','普通数据','普通',1440,120,14,'行业合规',40,'未指定行业时使用的通用备份合规基线。'),('LR_01','反洗钱法约束','金融','《反洗钱法》','实时交易记录','核心',5,60,2190,'行业合规',100,'金融实时交易数据要求分钟级RPO和长期离线保留。'),('LR_02','会计法约束','金融','《会计档案管理办法》','财务报告','核心',60,120,2190,'行业合规',90,'财务档案要求长期保留，并优先选择可恢复性更强的数据库备份策略。'),('LR_03','个保法约束','未定义','《个人信息保护法》','生物识别数据','敏感',1440,120,1825,'行业合规',85,'个人敏感信息要求加密与长期保留，优先保障完整性。'),('LR_04','政务数据约束','政府','《政务数据安全管理办法》','普通数据','普通',240,120,3650,'行业合规',95,'政务数据优先考虑长期归档和异地留存。'),('LR_05','能源生产日志约束','能源','《能源行业数据安全管理办法》','系统操作日志','重要',60,90,1095,'行业合规',90,'能源生产日志要求较低RPO和中长期留存。'),('LR_06','能源业务数据约束','能源','《能源行业数据安全管理办法》','普通数据','重要',240,120,365,'行业合规',75,'能源业务普通数据强调可恢复性和至少一年保留。');

COMMIT;

--
-- Table structure for table "recommendation_policy"
--

DROP TABLE IF EXISTS "recommendation_policy";


CREATE TABLE "recommendation_policy" (
  "id" varchar(32) NOT NULL,
  "name" varchar(32) DEFAULT NULL,
  "industry_type" varchar(32) DEFAULT NULL,
  "default_data_type" varchar(32) DEFAULT NULL,
  "default_rpo_minutes" integer DEFAULT NULL,
  "default_rto_minutes" integer DEFAULT NULL,
  "default_backup_window" varchar(32) DEFAULT NULL,
  "max_candidate_count" integer DEFAULT NULL,
  "matched_legal_rule_ids" varchar(32) DEFAULT NULL,
  "compliance_weight" integer DEFAULT NULL,
  "data_type_weight" integer DEFAULT NULL,
  "rpo_weight" integer DEFAULT NULL,
  "rto_weight" integer DEFAULT NULL,
  "window_weight" integer DEFAULT NULL,
  "priority_weight" integer DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "recommendation_policy"
--

BEGIN;

INSERT INTO "recommendation_policy" VALUES ('REC_00','通用候选推荐策略','未定义','普通数据',1440,120,'00:00-06:00',3,'LR_00,LR_03',30,15,20,10,10,15,'未定义行业的默认候选评估策略。'),('REC_ENERGY','能源候选推荐策略','能源','普通数据',1440,120,'00:00-06:00',3,'LR_00,LR_05,LR_06',35,10,20,15,5,15,'能源行业优先平衡合规、恢复时间与连续保护能力。'),('REC_FIN','金融候选推荐策略','金融','普通数据',1440,120,'00:00-06:00',3,'LR_00,LR_01,LR_02',40,10,20,10,5,15,'金融行业优先满足法规与低RPO要求。'),('REC_GOV','政务候选推荐策略','政府','普通数据',1440,120,'00:00-06:00',3,'LR_00,LR_04',45,10,10,10,10,15,'政务行业优先满足长期保留与归档要求。');

COMMIT;

--
-- Table structure for table "retention_plan"
--

DROP TABLE IF EXISTS "retention_plan";


CREATE TABLE "retention_plan" (
  "id" varchar(32) NOT NULL,
  "backup_solution_id" varchar(32) DEFAULT NULL,
  "name" varchar(32) DEFAULT NULL,
  "data_level" varchar(32) DEFAULT NULL,
  "storage_media" varchar(32) DEFAULT NULL,
  "data_retention_enabled" smallint DEFAULT NULL,
  "data_retention_mode" varchar(32) DEFAULT NULL,
  "retention_days" integer DEFAULT NULL,
  "data_retention_duration_value" integer DEFAULT NULL,
  "data_retention_duration_unit" varchar(32) DEFAULT NULL,
  "data_retention_count" integer DEFAULT NULL,
  "periodic_retention_enabled" smallint DEFAULT NULL,
  "daily_copy_count" integer DEFAULT NULL,
  "weekly_copy_count" integer DEFAULT NULL,
  "monthly_copy_count" integer DEFAULT NULL,
  "quarterly_copy_count" integer DEFAULT NULL,
  "yearly_copy_count" integer DEFAULT NULL,
  "log_retention_enabled" smallint DEFAULT NULL,
  "log_retention_duration_value" integer DEFAULT NULL,
  "log_retention_duration_unit" varchar(32) DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "retention_plan"
--

BEGIN;

INSERT INTO "retention_plan" VALUES ('PS_01_RET','PS_01','MySQL普通数据-逻辑增量策略-保留计划','普通','本地/同城/异地',1,'按时长保留',90,90,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:普通级基线保留-本地(RR_04_A)；同城:普通级基线保留-同城(RR_04_C)；异地:普通级基线保留-异地(RR_04_B)'),('PS_02_RET','PS_02','MySQL实时交易-物理实时日志策略-保留计划','核心','本地/同城/异地',1,'按时长保留',1825,5,'年',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:核心级合规保留-本地(RR_01_A)；同城:核心级交易保留-同城(RR_01_D)；异地:核心级交易保留-异地(RR_01_E)'),('PS_03_RET','PS_03','MySQL重要数据-逻辑增量策略-保留计划','重要','本地/同城/异地',1,'按时长保留',2190,6,'年',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:重要级合规保留-本地(RR_02_A)；同城:重要级合规保留-同城(RR_02_B)；异地:重要级合规保留-异地(RR_02_C)'),('PS_04_RET','PS_04','MySQL敏感数据-每日全备策略-保留计划','敏感','本地/异地',1,'按时长保留',1825,5,'年',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:敏感级合规保留-本地(RR_03_A)；异地:敏感级合规保留-异地(RR_03_B)'),('PS_05_RET','PS_05','Oracle实时交易-物理实时日志策略-保留计划','核心','本地/同城/异地',1,'按时长保留',365000,365000,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:核心级合规保留-本地(RR_01_A)；同城:核心级合规保留-同城(RR_01_B)；异地:核心级合规保留-异地(RR_01_C)'),('PS_06_RET','PS_06','Oracle财务核心-物理增强策略-保留计划','核心','本地/同城/异地',1,'按时长保留',365000,365000,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:核心级合规保留-本地(RR_01_A)；同城:核心级合规保留-同城(RR_01_B)；异地:核心级合规保留-异地(RR_01_C)'),('PS_07_RET','PS_07','文件普通数据-标准文件备份策略-保留计划','普通','本地/同城/异地',1,'按时长保留',90,90,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:普通级基线保留-本地(RR_04_A)；同城:普通级基线保留-同城(RR_04_C)；异地:普通级基线保留-异地(RR_04_B)'),('PS_08_RET','PS_08','文件敏感数据-加密文件备份策略-保留计划','敏感','本地/异地',1,'按时长保留',1825,5,'年',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:敏感级合规保留-本地(RR_03_A)；异地:敏感级合规保留-异地(RR_03_B)'),('PS_09_RET','PS_09','Oracle普通数据-物理标准策略-保留计划','普通','本地/同城/异地',1,'按时长保留',90,90,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:普通级基线保留-本地(RR_04_A)；同城:普通级基线保留-同城(RR_04_C)；异地:普通级基线保留-异地(RR_04_B)'),('PS_10_RET','PS_10','MySQL普通数据-每日全备稳妥策略-保留计划','普通','本地/同城/异地',1,'按时长保留',90,90,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:普通级基线保留-本地(RR_04_A)；同城:普通级基线保留-同城(RR_04_C)；异地:普通级基线保留-异地(RR_04_B)'),('PS_11_RET','PS_11','MySQL普通数据-高频增量增强策略-保留计划','普通','本地/同城/异地',1,'按时长保留',90,90,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:普通级基线保留-本地(RR_04_A)；同城:普通级基线保留-同城(RR_04_C)；异地:普通级基线保留-异地(RR_04_B)'),('PS_12_RET','PS_12','Oracle普通数据-高频增量增强策略-保留计划','普通','本地/同城/异地',1,'按时长保留',90,90,'天',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:普通级基线保留-本地(RR_04_A)；同城:普通级基线保留-同城(RR_04_C)；异地:普通级基线保留-异地(RR_04_B)'),('PS_13_RET','PS_13','文件普通数据-政务归档增强策略-保留计划','普通','本地/异地',1,'按时长保留',3650,10,'年',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:政务数据归档保留-本地(RR_05_A)；异地:政务数据归档保留-异地(RR_05_B)'),('PS_14_RET','PS_14','MySQL系统日志-能源高频增量策略-保留计划','重要','本地/同城/异地',1,'按时长保留',1095,3,'年',NULL,0,NULL,NULL,NULL,NULL,NULL,0,NULL,NULL,'本地:能源日志保留-本地(RR_06_A)；同城:能源日志保留-同城(RR_06_C)；异地:能源日志保留-异地(RR_06_B)');

COMMIT;

--
-- Table structure for table "rpo_rule"
--

DROP TABLE IF EXISTS "rpo_rule";


CREATE TABLE "rpo_rule" (
  "id" varchar(32) NOT NULL,
  "name" varchar(32) DEFAULT NULL,
  "min_time" integer DEFAULT NULL,
  "max_time" integer DEFAULT NULL,
  "app_type_req" varchar(32) DEFAULT NULL,
  "rule_scope" varchar(32) DEFAULT NULL,
  "description" text DEFAULT NULL,
  PRIMARY KEY ("id")
) ;


--
-- Dumping data for table "rpo_rule"
--

BEGIN;

INSERT INTO "rpo_rule" VALUES ('RPO_00','极致要求(支持日志的数据库)',0,4,'db_with_log_support','RPO约束','适用条件：db_with_log_support'),('RPO_01','极高要求',5,15,NULL,'RPO约束','极高要求'),('RPO_02','高要求',16,60,NULL,'RPO约束','高要求'),('RPO_03','一般要求',61,1440,NULL,'RPO约束','一般要求');

COMMIT;










-- Dump completed on 2026-04-26  9:36:59
