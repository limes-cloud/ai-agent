-- 001_app_partition.sql
-- ai-agent 按 app 纬度切分：model / prompt / prompt_classify 增加 app 列；model 增加 status；新增 call_log 表
-- 使用：mysql -h <host> -u <user> -p ai_agent < deploy/migrations/001_app_partition.sql
--
-- 说明：
--   1. 存量数据 app 默认 ''（空），需按业务归属执行 UPDATE（如 UPDATE model SET app='manager' WHERE app='';）
--   2. 若库中已存在多租户时代的残留索引（tenant_id / tenant_id_2），先执行下方“清理残留”段再执行本文件
--   3. 本文件与已执行的线上库保持一致（(app,*) 为唯一索引，替代原 uuid 全局唯一）

-- 清理残留（多租户时代索引，为唯一性的错误载体；不存在时请忽略报错）
-- ALTER TABLE `model` DROP INDEX `tenant_id`, DROP INDEX `tenant_id_2`;
-- ALTER TABLE `prompt` DROP INDEX `tenant_id`;

-- 1. model 表
ALTER TABLE `model`
  ADD COLUMN `app` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '所属应用（app 纬度隔离）' AFTER `id`,
  ADD COLUMN `status` TINYINT UNSIGNED NULL DEFAULT 1 COMMENT '状态：0 禁用 | 1 启用' AFTER `extra`,
  ADD UNIQUE INDEX `idx_model_app_uuid` (`app`, `uuid`);

-- 2. prompt 表
ALTER TABLE `prompt`
  ADD COLUMN `app` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '所属应用（app 纬度隔离）' AFTER `id`,
  ADD UNIQUE INDEX `idx_prompt_app_uuid` (`app`, `uuid`);

-- 3. prompt_classify 表
ALTER TABLE `prompt_classify`
  ADD COLUMN `app` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '所属应用（app 纬度隔离）' AFTER `id`,
  ADD UNIQUE INDEX `idx_prompt_classify_app_name` (`app`, `name`);

-- 4. call_log 表（生成调用日志，按 app 计量）
CREATE TABLE IF NOT EXISTS `call_log` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `app` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '应用标识',
  `model_id` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '模型 ID',
  `model_uuid` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '模型 UUID',
  `supplier` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '供应商',
  `scene` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '生成场景：text|image|prompt_text|prompt_image|embedding',
  `prompt_tokens` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '输入 token 数',
  `completion_tokens` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '输出 token 数',
  `latency_ms` BIGINT NOT NULL DEFAULT 0 COMMENT '生成耗时（毫秒）',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态：0 失败 | 1 成功',
  `error_message` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '失败原因',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（Unix 秒）',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '修改时间（Unix 秒）',
  PRIMARY KEY (`id`),
  KEY `idx_call_log_app_created` (`app`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI 生成调用日志';