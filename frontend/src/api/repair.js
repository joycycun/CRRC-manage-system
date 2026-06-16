import request from '@/utils/request'

export function getRepairRecords(params) {
  return request({
    url: '/repair-records',
    method: 'get',
    params
  })
}

export function createRepairRecord(data) {
  return request({
    url: '/repair-records',
    method: 'post',
    data
  })
}

export function updateRepairRecord(id, data) {
  return request({
    url: `/repair-records/${id}`,
    method: 'put',
    data
  })
}

export function deleteRepairRecord(id) {
  return request({
    url: `/repair-records/${id}`,
    method: 'delete'
  })
}

export function getFaultAnalysis(params) {
  return request({
    url: '/fault-analysis',
    method: 'get',
    params
  })
}

export function createFaultAnalysis(data) {
  return request({
    url: '/fault-analysis',
    method: 'post',
    data
  })
}

export function auditFaultAnalysis(id, data) {
  return request({
    url: `/fault-analysis/${id}/audit`,
    method: 'post',
    data
  })
}

export function deleteFaultAnalysis(id) {
  return request({
    url: `/fault-analysis/${id}`,
    method: 'delete'
  })
}
