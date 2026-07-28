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

export function updateProductionTestOutline(id, data) {
  return request({
    url: `/production-test-outlines/${id}`,
    method: 'put',
    data
  })
}

export function deleteProductionTestOutline(id) {
  return request({
    url: `/production-test-outlines/${id}`,
    method: 'delete'
  })
}
