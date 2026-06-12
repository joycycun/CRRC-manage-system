import request from '@/utils/request'

export function getOutboundRecords(params) {
  return request({
    url: '/outbound-records',
    method: 'get',
    params
  })
}

export function returnOutboundRecord(id) {
  return request({
    url: `/outbound-records/${id}/return`,
    method: 'post'
  })
}