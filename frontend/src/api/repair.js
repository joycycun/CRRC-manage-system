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
