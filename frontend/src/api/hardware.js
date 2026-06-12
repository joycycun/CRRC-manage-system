import request from '@/utils/request'

// 查询硬件版本
export function getHardwareVersions(params) {
  return request({
    url: '/hardware-versions',
    method: 'get',
    params
  })
}

// 新增硬件版本
export function createHardwareVersion(data) {
  return request({
    url: '/hardware-versions',
    method: 'post',
    data
  })
}

// 修改硬件版本
export function updateHardwareVersion(id, data) {
  return request({
    url: `/hardware-versions/${id}`,
    method: 'put',
    data
  })
}

// 上传硬件 ZIP
export function uploadHardwareZip(id, data) {
  return request({
    url: `/hardware-versions/${id}/upload-zip`,
    method: 'post',
    data
  })
}

// =========================
// 硬件测试记录
// =========================

export function getHardwareTests(params) {
  return request({
    url: '/hardware-tests',
    method: 'get',
    params
  })
}

export function createHardwareTest(data) {
  return request({
    url: '/hardware-tests',
    method: 'post',
    data
  })
}
export function submitHardwareTest(id) {
  return request({
    url: `/hardware-tests/${id}/submit`,
    method: 'post'
  })
}
export function auditHardwareTest(id, data) {
  return request({
    url: `/hardware-tests/${id}/audit`,
    method: 'post',
    data
  })
}

export function deleteHardwareTest(id) {
  return request({
    url: `/hardware-tests/${id}`,
    method: 'delete'
  })
}