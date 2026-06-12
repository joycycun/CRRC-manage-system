import request from '@/utils/request'

// 查询需求书
export function getRequirementBooks(params) {
  return request({
    url: '/requirement-books',
    method: 'get',
    params
  })
}

// 新增需求书
export function createRequirementBook(data) {
  return request({
    url: '/requirement-books',
    method: 'post',
    data
  })
}

// 提交需求书
export function submitRequirementBook(id) {
  return request({
    url: `/requirement-books/${id}/submit`,
    method: 'post'
  })
}

// 审核需求书
export function auditRequirementBook(id, data) {
  return request({
    url: `/requirement-books/${id}/audit`,
    method: 'post',
    data
  })
}

export function approveRequirementBook(id) {
  return request({
    url: `/requirement-books/${id}/approve`,
    method: 'post'
  })
}

export function rejectRequirementBook(id, data) {
  return request({
    url: `/requirement-books/${id}/reject`,
    method: 'post',
    data
  })
}

// 删除需求书
export function deleteRequirementBook(id) {
  return request({
    url: `/requirement-books/${id}`,
    method: 'delete'
  })
}

// =========================
// 需求变更
// =========================

export function getRequirementChanges(params) {
  return request({
    url: '/requirement-changes',
    method: 'get',
    params
  })
}

export function createRequirementChange(data) {
  return request({
    url: '/requirement-changes',
    method: 'post',
    data
  })
}

export function submitRequirementChange(id) {
  return request({
    url: `/requirement-changes/${id}/submit`,
    method: 'post'
  })
}

export function auditRequirementChange(id, data) {
  return request({
    url: `/requirement-changes/${id}/audit`,
    method: 'post',
    data
  })
}

export function closeRequirementChange(id, data) {
  return request({
    url: `/requirement-changes/${id}/close`,
    method: 'post',
    data
  })
}

export function deleteRequirementChange(id) {
  return request({
    url: `/requirement-changes/${id}`,
    method: 'delete'
  })
}

// =========================
// 客供资料
// =========================

export function getCustomerSuppliedFiles(params) {
  return request({
    url: '/customer-supplied-files',
    method: 'get',
    params
  })
}

export function createCustomerSuppliedFile(data) {
  return request({
    url: '/customer-supplied-files',
    method: 'post',
    data
  })
}

export function deleteCustomerSuppliedFile(id) {
  return request({
    url: `/customer-supplied-files/${id}`,
    method: 'delete'
  })
}