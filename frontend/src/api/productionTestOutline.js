import request from '@/utils/request'

export function getProductionTestOutlines(params) {
  return request({
    url: '/production-test-outlines',
    method: 'get',
    params
  })
}

export function createProductionTestOutline(data) {
  return request({
    url: '/production-test-outlines',
    method: 'post',
    data
  })
}
