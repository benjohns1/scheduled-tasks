import fetch from 'node-fetch'
import { withAuth } from './api.auth'

const baseUrl = `http://${process.env.APPLICATION_HOST || 'localhost'}:${process.env.APPLICATION_PORT || '8080'}/api/v1/schedule`

export async function getAll() {
  return fetch(`${baseUrl}`, await withAuth())
}

export async function add(data) {
  return fetch(`${baseUrl}/`, await withAuth({ method: "POST", body: JSON.stringify(data) }))
}

export async function addRecurringTask(id, data) {
  return fetch(`${baseUrl}/${id}/task`, await withAuth({ method: "POST", body: JSON.stringify(data) }))
}

export async function pause(id) {
  return fetch(`${baseUrl}/${id}/pause`, await withAuth({ method: "PUT" }))
}

export async function unpause(id) {
  return fetch(`${baseUrl}/${id}/unpause`, await withAuth({ method: "PUT" }))
}

export async function remove(id) {
  return fetch(`${baseUrl}/${id}`, await withAuth({ method: "DELETE" }))
}