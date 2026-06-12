import request from '@/utils/request'

// 查询烧录记录
export function getBurnRecords(params) {
  return request({
    url: '/burn-records',
    method: 'get',
    params
  })
}

// 批量导入烧录记录
export function importBurnRecords(data) {
  return request({
    url: '/burn-records/import',
    method: 'post',
    data
  })
}

// 删除单条烧录记录
export function deleteBurnRecordApi(id) {
  return request({
    url: `/burn-records/${id}`,
    method: 'delete'
  })
}

// 删除某个生产批次下全部烧录记录
export function deleteBurnBatchApi(batchNo) {
  return request({
    url: `/burn-records/batch/${encodeURIComponent(batchNo)}`,
    method: 'delete'
  })
}

export function getBurnRecordOptions() {
  return request({
    url: '/burn-records/options',
    method: 'get'
  })
}