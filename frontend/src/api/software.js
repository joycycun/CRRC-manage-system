import request from '@/utils/request'

// =========================
// 软件版本
// =========================

export function getSoftwareVersions(params) {
  return request({
    url: '/software-versions',
    method: 'get',
    params
  })
}

export function createSoftwareVersion(data) {
  return request({
    url: '/software-versions',
    method: 'post',
    data
  })
}

export function updateSoftwareVersion(id, data) {
  return request({
    url: `/software-versions/${id}`,
    method: 'put',
    data
  })
}

export function releaseSoftwareVersion(id) {
  return request({
    url: `/software-versions/${id}/release`,
    method: 'post'
  })
}

export function discardSoftwareVersion(id) {
  return request({
    url: `/software-versions/${id}/discard`,
    method: 'post'
  })
}

export function deleteSoftwareVersion(id) {
  return request({
    url: `/software-versions/${id}`,
    method: 'delete'
  })
}


// =========================
// 项目分支
// =========================

export function getBranches(params) {
  return request({
    url: '/branches',
    method: 'get',
    params
  })
}

export function createBranch(data) {
  return request({
    url: '/branches',
    method: 'post',
    data
  })
}

export function updateBranch(id, data) {
  return request({
    url: `/branches/${id}`,
    method: 'put',
    data
  })
}

export function deleteBranch(id) {
  return request({
    url: `/branches/${id}`,
    method: 'delete'
  })
}