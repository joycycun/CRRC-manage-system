import request from '@/utils/request'

export function getDashboardSummary(params) {
  return request({
    url: '/dashboard/summary',
    method: 'get',
    params
  })
}

export function getProjectProgressReport(params) {
  return request({
    url: '/reports/project-progress',
    method: 'get',
    params
  })
}

export function getVersionMatrixReport(params) {
  return request({
    url: '/reports/version-matrix',
    method: 'get',
    params
  })
}

export function getIssueStatisticsReport(params) {
  return request({
    url: '/reports/issue-statistics',
    method: 'get',
    params
  })
}
