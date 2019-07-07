export function withJsonAndAuth(token, fetchOptions = { headers: {} }) {
  if (fetchOptions.headers === undefined) {
    fetchOptions.headers = {}
  }
  fetchOptions.headers['content-type'] = 'application/json'
  if (token) {
    fetchOptions.headers['authorization'] = `Bearer ${token}`
  }
  return fetchOptions
}