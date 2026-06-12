import request from '@/utils/request'

// =========================
// 测试用例
// =========================

export function getTestCases(params) {
  return request({
    url: '/test-cases',
    method: 'get',
    params
  })
}

export function createTestCase(data) {
  return request({
    url: '/test-cases',
    method: 'post',
    data
  })
}

export function submitTestCaseApi(id) {
  return request({
    url: `/test-cases/${id}/submit`,
    method: 'post'
  })
}

export function auditTestCaseApi(id, data) {
  return request({
    url: `/test-cases/${id}/audit`,
    method: 'post',
    data
  })
}

export function deleteTestCaseApi(id) {
  return request({
    url: `/test-cases/${id}`,
    method: 'delete'
  })
}

// =========================
// 测试报告
// =========================

export function uploadTestReportApi(testCaseId, data) {
  return request({
    url: `/test-cases/${testCaseId}/report`,
    method: 'post',
    data
  })
}