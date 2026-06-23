export function getCurrentUser() {
  try {
    const userText = localStorage.getItem('user')
    return userText ? JSON.parse(userText) : {}
  } catch (err) {
    console.warn('读取当前用户失败：', err)
    return {}
  }
}

export function getCurrentUserParams() {
  const user = getCurrentUser()

  return {
    userId: user.id || '',
    username: user.username || '',
    realName: user.realName || '',
    department: user.department || ''
  }
}

export function getCurrentUserName(fallback = '当前用户') {
  const user = getCurrentUser()
  return user.realName || user.username || localStorage.getItem('realName') || fallback
}

export function getAuditUserPayload() {
  const user = getCurrentUser()
  const name = getCurrentUserName('审核人')
  return {
    auditorId: user.id || 0,
    auditorName: name,
    auditUserId: user.id || 0,
    auditUserName: name
  }
}
