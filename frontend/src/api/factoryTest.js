import request from '@/utils/request'

export function getFactoryTests(params) {
  return request({
    url: '/factory-tests',
    method: 'get',
    params
  })
}

export function importFactoryTests(data) {
  return request({
    url: '/factory-tests/import',
    method: 'post',
    data
  })
}

export function deleteFactoryTests(data) {
  return request({
    url: '/factory-tests/delete',
    method: 'post',
    data
  })
}

export function submitFactoryTests(data) {
  return request({
    url: '/factory-tests/submit',
    method: 'post',
    data
  })
}

export function auditFactoryTests(data) {
  return request({
    url: '/factory-tests/audit',
    method: 'post',
    data
  })
}