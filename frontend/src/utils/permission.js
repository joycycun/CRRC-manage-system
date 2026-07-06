const FULL_ACCESS_ROLES = ['system_admin', 'leader']
const LEADER_ROLES = ['leader', 'system_admin']
const LIMITED_ROLES = ['software_owner', 'hardware_owner', 'project_assistant', 'production_staff', 'shipping_staff', 'shipping_auditor', 'aftersales_staff', 'quality_staff']

const PAGE_ACCESS = {
  software_owner: [
    '/dashboard',
    '/project/manage',
    '/requirement/book',
    '/requirement/change',
    '/requirement/customer-supplied',
    '/software/version',
    '/software/branch',
    '/test/case',
    '/test/issue',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ],
  hardware_owner: [
    '/dashboard',
    '/project/manage',
    '/requirement/book',
    '/requirement/change',
    '/requirement/customer-supplied',
    '/hardware/version',
    '/hardware/test',
    '/production/test-outline',
    '/test/case',
    '/test/issue',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ],
  project_assistant: [
    '/dashboard',
    '/project/manage',
    '/requirement/book',
    '/requirement/change',
    '/requirement/customer-supplied',
    '/hardware/version',
    '/hardware/test',
    '/software/version',
    '/software/branch',
    '/test/case',
    '/test/issue',
    '/aftersales/repair',
    '/aftersales/fault-analysis',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ],
  production_staff: [
    '/dashboard',
    '/project/manage',
    '/production/board-inbound',
    '/production/test-outline',
    '/production/burn',
    '/production/factory-test',
    '/production/inventory',
    '/aftersales/repair',
    '/aftersales/fault-analysis',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ],
  shipping_staff: [
    '/dashboard',
    '/project/manage',
    '/hardware/version',
    '/aftersales/repair',
    '/aftersales/fault-analysis',
    '/shipping/out',
    '/shipping/batch',
    '/inventory/out',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ],
  shipping_auditor: [
    '/dashboard',
    '/project/manage',
    '/shipping/out',
    '/shipping/batch',
    '/inventory/out',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ],
  aftersales_staff: [
    '/dashboard',
    '/project/manage',
    '/requirement/book',
    '/requirement/change',
    '/requirement/customer-supplied',
    '/hardware/version',
    '/hardware/test',
    '/software/version',
    '/software/branch',
    '/test/case',
    '/test/issue',
    '/aftersales/repair',
    '/aftersales/fault-analysis',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ],
  quality_staff: [
    '/dashboard',
    '/project/manage',
    '/production/test-outline',
    '/production/factory-test',
    '/report/project-progress',
    '/report/version-matrix',
    '/report/issue-statistics',
    '/project/progress-report',
    '/version/matrix'
  ]
}

const DEFAULT_PAGE = {
  software_owner: '/software/version',
  hardware_owner: '/hardware/version',
  project_assistant: '/project/manage',
  production_staff: '/dashboard',
  shipping_staff: '/dashboard',
  shipping_auditor: '/dashboard',
  aftersales_staff: '/dashboard',
  quality_staff: '/production/factory-test'
}

const ACTION_ACCESS = {
  software_owner: [
    'project:view',
    'requirement:view',
    'requirement:download',
    'customer:view',
    'customer:download',
    'software:*',
    'testcase:view',
    'testcase:download',
    'issue:view',
    'issue:reply',
    'issue:export',
    'report:*'
  ],
  hardware_owner: [
    'project:view',
    'requirement:view',
    'requirement:download',
    'customer:view',
    'customer:download',
    'hardware:*',
    'production:outline:view',
    'production:outline:upload',
    'testcase:view',
    'testcase:download',
    'issue:view',
    'issue:reply',
    'issue:export',
    'report:*'
  ],
  project_assistant: [
    'project:*',
    'requirement:*',
    'customer:*',
    'hardware:view',
    'hardware:download',
    'software:view',
    'software:download',
    'branch:view',
    'branch:create',
    'branch:download',
    'testcase:*',
    'issue:*',
    'aftersales:view',
    'report:*'
  ],
  production_staff: [
    'project:view',
    'production:*',
    'aftersales:*',
    'report:*'
  ],
  shipping_staff: [
    'project:view',
    'hardware:view',
    'hardware:download',
    'aftersales:view',
    'aftersales:download',
    'shipping:view',
    'shipping:manage',
    'shipping:create',
    'shipping:update',
    'shipping:delete',
    'shipping:submit',
    'report:*'
  ],
  shipping_auditor: [
    'project:view',
    'shipping:view',
    'shipping:audit',
    'report:*'
  ],
  aftersales_staff: [
    'project:view',
    'requirement:view',
    'requirement:download',
    'customer:view',
    'customer:download',
    'hardware:view',
    'hardware:download',
    'software:view',
    'software:download',
    'branch:view',
    'branch:download',
    'testcase:view',
    'testcase:download',
    'issue:view',
    'aftersales:*',
    'report:*'
  ],
  quality_staff: [
    'project:view',
    'production:view',
    'production:outline:view',
    'production:audit',
    'report:*'
  ]
}

export function getStoredRoles() {
  try {
    const text = localStorage.getItem('roles')
    const roles = text ? JSON.parse(text) : []
    return roles
      .map(role => role.roleCode || role.role_code || role.code || role)
      .filter(Boolean)
  } catch (err) {
    console.warn('读取角色权限失败：', err)
    return []
  }
}

export function hasRole(roleCode) {
  return getStoredRoles().includes(roleCode)
}

export function hasLeaderRole() {
  const roles = getStoredRoles()
  return roles.some(role => LEADER_ROLES.includes(role))
}

export function hasFullAccessRole() {
  const roles = getStoredRoles()
  return roles.some(role => FULL_ACCESS_ROLES.includes(role))
}

function shouldRestrictCurrentUser() {
  const roles = getStoredRoles()
  if (roles.length === 0 || hasFullAccessRole()) return false
  return roles.some(role => LIMITED_ROLES.includes(role))
}

export function canAccessPage(path) {
  if (!shouldRestrictCurrentUser()) return true

  const roles = getStoredRoles()
  return roles.some(role => (PAGE_ACCESS[role] || []).includes(path))
}

export function getDefaultAccessiblePage() {
  if (!shouldRestrictCurrentUser()) return '/dashboard'

  const roles = getStoredRoles()
  for (const role of roles) {
    if (DEFAULT_PAGE[role]) return DEFAULT_PAGE[role]
  }

  for (const role of roles) {
    const pages = PAGE_ACCESS[role] || []
    if (pages.length > 0) return pages[0]
  }

  return '/login'
}

function actionMatches(grantedAction, action) {
  if (grantedAction === action) return true
  if (grantedAction.endsWith(':*')) {
    return action.startsWith(grantedAction.replace('*', ''))
  }
  return false
}

export function canUseAction(action) {
  if (hasRole('system_admin')) return true

  if (action === 'production:outline:upload') {
    return hasRole('hardware_owner') || hasRole('system_admin')
  }

  if (action === 'board-inbound:import') {
    return hasRole('production_staff') || hasRole('system_admin')
  }

  if (action === 'production:outline:view') {
    return hasRole('hardware_owner') || hasRole('production_staff') || hasRole('quality_staff') || hasRole('system_admin') || hasRole('leader')
  }

  if (action === 'customer:delete') {
    return hasRole('project_assistant') || hasRole('system_admin')
  }

  if (action === 'issue:update') {
    return hasRole('project_assistant')
  }

  if (action === 'requirement:close') {
    return hasRole('project_assistant')
  }

  if (action === 'project:reopen') {
    return hasRole('project_assistant') || hasRole('system_admin')
  }

  if (action === 'burn:deleteBatch') {
    return hasRole('production_staff') || hasRole('system_admin')
  }

  if (action === 'production:audit') {
    return hasRole('quality_staff') || hasRole('system_admin')
  }

  if (action === 'shipping:audit') {
    return hasLeaderRole() || hasRole('shipping_auditor')
  }

  if (action.endsWith(':audit')) {
    return hasLeaderRole()
  }

  if (!shouldRestrictCurrentUser()) return true

  const roles = getStoredRoles()
  return roles.some(role => {
    return (ACTION_ACCESS[role] || []).some(grantedAction => actionMatches(grantedAction, action))
  })
}
