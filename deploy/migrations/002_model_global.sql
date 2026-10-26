-- 002_model_global.sql
-- 模型改为全局资源：移除 app 维度，uuid 全局唯一
-- 使用：mysql -h <host> -u <user> -p ai_agent < deploy/migrations/002_model_global.sql
--
-- 说明：
--   1. 执行前先确认无跨 app 的重复 uuid，否则唯一索引创建失败：
--      SELECT uuid, COUNT(*) FROM model GROUP BY uuid HAVING COUNT(*) > 1;
--      如有重复需先人工合并。
--   2. 本迁移保留原有模型数据，仅删除 app 列；prompt / call_log 仍按 app 隔离。

ALTER TABLE `model`
  DROP INDEX `idx_model_app_uuid`,
  ADD UNIQUE INDEX `uk_model_uuid` (`uuid`),
  DROP COLUMN `app`;
