import request from '@/utils/request'

export function getCrons(params) {
  return request({
    url: '/api/crons',
    method: 'get',
    params
  })
}

export function createCron(data) {
  return request({
    url: '/api/crons',
    method: 'post',
    data
  })
}

export function updateCron(id, data) {
  return request({
    url: `/api/crons/${id}`,
    method: 'put',
    data
  })
}

export function deleteCron(id) {
  return request({
    url: `/api/crons/${id}`,
    method: 'delete'
  })
}
