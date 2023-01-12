import request from '@/utils/request'

export function getCredentials(params) {
  return request({
    url: '/api/credentials',
    method: 'get',
    params
  })
}

export function createCredential(data) {
  return request({
    url: '/api/credentials',
    method: 'post',
    data
  })
}

export function updateCredential(id, data) {
  return request({
    url: `/api/credentials/${id}`,
    method: 'put',
    data
  })
}

export function deleteCredential(id) {
  return request({
    url: `/api/credentials/${id}`,
    method: 'delete'
  })
}
