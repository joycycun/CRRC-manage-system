import request from '@/utils/request'

export function getHardwareDevDocuments(params) {
  return request({
    url: '/hardware-dev-documents',
    method: 'get',
    params
  })
}

export function createHardwareDevDocument(data) {
  return request({
    url: '/hardware-dev-documents',
    method: 'post',
    data,
    timeout: 120000
  })
}

export function deleteHardwareDevDocument(id) {
  return request({
    url: `/hardware-dev-documents/${id}`,
    method: 'delete'
  })
}
