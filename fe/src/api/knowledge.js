import request from '@/utils/request'

export function getKnowledgeBaseList(params) {
  return request({
    url: '/v1/knowledge-base',
    method: 'get',
    params
  })
}

export function createKnowledgeBase(data) {
  return request({
    url: '/v1/knowledge-base',
    method: 'post',
    data
  })
}

export function updateKnowledgeBase(data) {
  return request({
    url: `/v1/knowledge-base/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteKnowledgeBase(id) {
  return request({
    url: `/v1/knowledge-base/${id}`,
    method: 'delete'
  })
} 