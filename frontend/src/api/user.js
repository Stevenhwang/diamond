import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/api/login',
    method: 'post',
    data
  })
}

export function getInfo() {
  return request({
    url: '/api/user_info',
    method: 'get'
  })
}

export function resetPass(data) {
  return request({
    url: '/api/reset_pw',
    method: 'post',
    data
  })
}

export function getUsers(params) {
  return request({
    url: '/api/users',
    method: 'get',
    params
  })
}

export function createUser(data) {
  return request({
    url: '/api/users',
    method: 'post',
    data
  })
}

export function updateUser(id, data) {
  return request({
    url: `/api/users/${id}`,
    method: 'put',
    data
  })
}

export function deleteUser(id) {
  return request({
    url: `/api/users/${id}`,
    method: 'delete'
  })
}

export function getUserPerms(params) {
  return request({
    url: '/api/userPerms',
    method: 'get',
    params
  })
}

export function assignUserPerm(data) {
  return request({
    url: '/api/userPerms',
    method: 'post',
    data
  })
}

export function getUserServers(params) {
  return request({
    url: '/api/userServers',
    method: 'get',
    params
  })
}

export function assignUserServer(data) {
  return request({
    url: '/api/userServers',
    method: 'post',
    data
  })
}

export function syncPerms() {
  return request({
    url: '/api/syncPerms',
    method: 'post',
  })
}
