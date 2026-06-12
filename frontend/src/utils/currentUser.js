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
