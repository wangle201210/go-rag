CREATE TABLE IF NOT EXISTS `document_mapping` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `knowledge_base_id` bigint(20) NOT NULL COMMENT '知识库ID',
  `document_id` varchar(64) NOT NULL COMMENT '文档ID',
  `document_name` varchar(255) NOT NULL COMMENT '文档名称',
  `document_type` varchar(20) NOT NULL COMMENT '文档类型：md/pdf/html',
  `document_path` varchar(500) NOT NULL COMMENT '文档路径',
  `document_size` bigint(20) NOT NULL COMMENT '文档大小(字节)',
  `chunk_count` int(11) NOT NULL DEFAULT '0' COMMENT '分块数量',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_doc_id` (`document_id`),
  KEY `idx_kb_id` (`knowledge_base_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文档映射关系表'; 