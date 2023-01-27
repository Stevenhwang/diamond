import request from '@/utils/request'

export function getScripts(params) {
  return request({
    url: '/api/scripts',
    method: 'get',
    params
  })
}

export function createScript(data) {
  return request({
    url: '/api/scripts',
    method: 'post',
    data
  })
}

export function updateScript(id, data) {
  return request({
    url: `/api/scripts/${id}`,
    method: 'put',
    data
  })
}

export function deleteScript(id) {
  return request({
    url: `/api/scripts/${id}`,
    method: 'delete'
  })
}

