-- 003_psou_resource_prompt.sql
-- PSOU 资源管理：根据资源标题和所属分类生成上架文案
-- 使用：mysql -h <host> -u <user> -p ai_agent < deploy/migrations/003_psou_resource_prompt.sql

SET NAMES utf8mb4;

SET @system = '你是一位严谨的泛资源内容编辑，负责为影视、音频、音乐、软件、游戏、电子书等资源生成上架文案。只依据资源标题和所属分类写作，不虚构演员、歌手、年份、版本、开发商、曲目或具体功能。分类只作为内容约束，不作为输出字段。文案使用简体中文，客观自然，不使用绝对化宣传。只输出符合给定 JSON Schema 的原始 JSON，不输出 Markdown、解释或额外字段。';

SET @prompt = '请根据以下输入生成资源上架信息。

资源标题：{{title}}
所属分类：{{category}}

要求：
1. title：规范化后的资源标题，不超过 80 个汉字。保留作品或软件的真实名称以及标题中已有的必要版本、年份信息，不得自行补造。
2. summary：50 到 120 个汉字的一句话摘要，点明资源类型、核心内容或用途、适合人群。不要只重复分类名，不写夸张宣传。
3. description：180 到 400 个汉字，分成 2 到 4 个短段落，段落之间用两个换行分隔。介绍内容主体、主要看点、功能或适用场景，并给出必要的使用提醒。只写可根据标题和分类合理判断的内容，不罗列未经证实的具体事实。
4. tags：输出 5 到 8 个去重标签，每个 2 到 8 个汉字，不含 #、标点和空格，优先选择主题、类型、风格、用途等便于检索的词。
5. score：0.0 到 10.0 的编辑推荐分，保留一位小数。综合内容完整度、实用性和受众广度评分；信息不足时给 6.5 到 7.8，不根据标题知名度臆断。
6. confidence：0 到 1 的置信度。名称清晰且分类匹配时可给 0.75 以上；名称泛化、存在同名歧义、分类不匹配或信息不足时不超过 0.6。
7. 分类仅用于约束写作，不得在输出中返回分类字段。
8. 不生成或提及封面、图片、图片 URL。
9. 所有字段都必须输出。

仅输出原始 JSON，不要 Markdown 代码块。';

SET @output_schema = '{"type":"object","properties":{"title":{"type":"string","description":"规范化后的资源标题，不超过80个汉字"},"summary":{"type":"string","description":"50到120个汉字的一句话摘要"},"description":{"type":"string","description":"180到400个汉字的详情，2到4个短段落，段落之间用两个换行分隔"},"tags":{"type":"array","minItems":5,"maxItems":8,"items":{"type":"string","description":"2到8个汉字，不含#和标点"}},"score":{"type":"number","minimum":0,"maximum":10,"description":"编辑推荐分，保留一位小数"},"confidence":{"type":"number","minimum":0,"maximum":1,"description":"名称、分类和内容生成的置信度"}},"required":["title","summary","description","tags","score","confidence"],"additionalProperties":false}';

SET @extra = JSON_OBJECT(
  'temperature', 0.35,
  'topP', 0.8,
  'formatJson', TRUE,
  'outputSchema', @output_schema
);

INSERT INTO `prompt_classify` (`app`, `name`, `description`, `created_at`, `updated_at`)
VALUES ('psou', '资源管理', 'PSOU 资源上架文案生成', UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE
  `description` = VALUES(`description`),
  `updated_at` = VALUES(`updated_at`);

SET @classify_id = (
  SELECT `id` FROM `prompt_classify` WHERE `app` = 'psou' AND `name` = '资源管理' LIMIT 1
);

INSERT INTO `prompt` (
  `app`, `classify_id`, `model_id`, `uuid`, `title`, `description`, `system`, `prompt`, `extra`, `status`, `version`, `created_at`, `updated_at`
)
VALUES (
  'psou', @classify_id, 8, 'psou-resource-metadata', 'PSOU 资源信息生成', '根据资源标题和所属分类生成摘要、详情、标签与评分', @system, @prompt, @extra, 1, MD5(@prompt), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
)
ON DUPLICATE KEY UPDATE
  `classify_id` = VALUES(`classify_id`),
  `model_id` = VALUES(`model_id`),
  `title` = VALUES(`title`),
  `description` = VALUES(`description`),
  `system` = VALUES(`system`),
  `prompt` = VALUES(`prompt`),
  `extra` = VALUES(`extra`),
  `status` = VALUES(`status`),
  `version` = VALUES(`version`),
  `updated_at` = VALUES(`updated_at`);
