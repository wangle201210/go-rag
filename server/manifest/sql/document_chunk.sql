CREATE TABLE IF NOT EXISTS `document_chunk` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `document_id` varchar(64) NOT NULL COMMENT '文档ID',
    `chunk_id` varchar(64) NOT NULL COMMENT '块ID',
    `content` text NOT NULL COMMENT '块内容',
    `metadata` text NOT NULL COMMENT '块元数据',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_chunk_id` (`chunk_id`),
    KEY `idx_document_id` (`document_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档块表'; 