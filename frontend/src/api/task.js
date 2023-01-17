import request from '@/utils/request'

export function getTasks(params) {
  return request({
    url: '/api/tasks',
    method: 'get',
    params
  })
}

export function createTask(data) {
  return request({
    url: '/api/tasks',
    method: 'post',
    data
  })
}

export function updateTask(id, data) {
  return request({
    url: `/api/tasks/${id}`,
    method: 'put',
    data
  })
}

export function deleteTask(id) {
  return request({
    url: `/api/tasks/${id}`,
    method: 'delete'
  })
}

export function invokeTask(id) {
  return request({
    url: `/api/tasks/${id}`,
    method: 'post'
  })
}

export function getTaskHist(params) {
  return request({
    url: '/api/taskhist',
    method: 'get',
    params
  })
}
