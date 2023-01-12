import request from '@/utils/request'

export function getServers(params) {
  return request({
    url: '/api/servers',
    method: 'get',
    params
  })
}

export function createServer(data) {
  return request({
    url: '/api/servers',
    method: 'post',
    data
  })
}

export function updateServer(id, data) {
  return request({
    url: `/api/servers/${id}`,
    method: 'put',
    data
  })
}

export function deleteServer(id) {
  return request({
    url: `/api/servers/${id}`,
    method: 'delete'
  })
}

export function getRecords(params) {
  return request({
    url: '/api/records',
    method: 'get',
    params
  })
}
