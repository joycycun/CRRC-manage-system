import request from '@/utils/request'

export function getProjects(params) {
  return request({
    url: '/projects',
    method: 'get',
    params
  })
}

export function getProjectDetail(id) {
  return request({
    url: `/projects/${id}`,
    method: 'get'
  })
}

export function createProject(data) {
  return request({
    url: '/projects',
    method: 'post',
    data
  })
}

export function updateProject(id, data) {
  return request({
    url: `/projects/${id}`,
    method: 'put',
    data
  })
}

export function submitProject(id) {
  return request({
    url: `/projects/${id}/submit`,
    method: 'post'
  })
}

export function auditProject(id, data) {
  return request({
    url: `/projects/${id}/audit`,
    method: 'post',
    data
  })
}

export function archiveProject(id) {
  return request({
    url: `/projects/${id}/archive`,
    method: 'post'
  })
}

export function closeProject(id) {
  return request({
    url: `/projects/${id}/close`,
    method: 'post'
  })
}

export function deleteProject(id) {
  return request({
    url: `/projects/${id}`,
    method: 'delete'
  })
}