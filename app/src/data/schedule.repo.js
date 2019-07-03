import fetch from 'node-fetch'

const baseUrl = `http://${process.env.APPLICATION_HOST || 'localhost'}:${process.env.APPLICATION_PORT || '8080'}/api/v1/schedule`

export async function getAll() {
  return fetch(`${baseUrl}`)
}

export async function add(data) {
  return fetch(`${baseUrl}/`, { method: "POST", body: JSON.stringify(data) })
}

export async function addRecurringTask(id, data) {
  return fetch(`${baseUrl}/${id}/task`, { method: "POST", body: JSON.stringify(data) })
}

export async function pause(id) {
  return fetch(`${baseUrl}/${id}/pause`, { method: "PUT" })
}

export async function unpause(id) {
  return fetch(`${baseUrl}/${id}/unpause`, { method: "PUT" })
}

export async function remove(id) {
  return fetch(`${baseUrl}/${id}`, { method: "DELETE" })
}