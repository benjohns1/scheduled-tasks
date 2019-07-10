export function withJsonAndAuth(session, fetchOptions = { headers: {} }) {
  if (fetchOptions.headers === undefined) {
    fetchOptions.headers = {}
  }
  fetchOptions.headers['content-type'] = 'application/json'
  if (session && session.auth && session.auth.token) {
    fetchOptions.headers['authorization'] = `Bearer ${session.auth.token}`
  }
  return fetchOptions
}