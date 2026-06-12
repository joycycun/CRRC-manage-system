import request from '@/utils/request'

// 查询问题列表
export function getIssues(params) {
  return request({
    url: '/issues',
    method: 'get',
    params
  })
}

// 新增问题
export function createIssue(data) {
  return request({
    url: '/issues',
    method: 'post',
    data
  })
}

// 修改问题
export function updateIssue(id, data) {
  return request({
    url: `/issues/${id}`,
    method: 'put',
    data
  })
}

// 回复问题
export function replyIssue(id, data) {
  return request({
    url: `/issues/${id}/reply`,
    method: 'post',
    data
  })
}

// 关闭问题
export function closeIssueApi(id, data) {
  return request({
    url: `/issues/${id}/close`,
    method: 'post',
    data
  })
}

// 重新打开问题
export function reopenIssueApi(id, data) {
  return request({
    url: `/issues/${id}/reopen`,
    method: 'post',
    data
  })
}

// 删除问题，当前页面没按钮，先预留
export function deleteIssue(id) {
  return request({
    url: `/issues/${id}`,
    method: 'delete'
  })
}