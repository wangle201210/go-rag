import request from '@/utils/request'

// 获取文档列表
export function getDocumentList(params) {
  return request({
    url: '/api/v1/document/list',
    method: 'get',
    params
  })
}

// 获取文档详情
export function getDocumentDetail(id) {
  return request({
    url: `/api/v1/document/${id}`,
    method: 'get'
  })
}

// 删除文档
export function deleteDocument(id) {
  return request({
    url: `/api/v1/document/${id}`,
    method: 'delete'
  })
}

// 更新文档
export function updateDocument(id, data) {
  return request({
    url: `/api/v1/document/${id}`,
    method: 'put',
    data
  })
}

// 上传文档
export function uploadDocument(data) {
  return request({
    url: '/api/v1/document/upload',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
} 