import fetch from 'node-fetch'
import { withAuth } from './api.auth'

const baseUrl = `http://${process.env.APPLICATION_HOST || 'localhost'}:${process.env.APPLICATION_PORT || '8080'}/api/v1/task`

export async function getAll() {
  return fetch(`${baseUrl}`, await withAuth())
}

export async function add(data) {
  return fetch(`${baseUrl}/`, await withAuth({ method: "POST", body: JSON.stringify(data)}))
}

export async function complete(id) {
  return fetch(`${baseUrl}/${id}/complete`, await withAuth({ method: "PUT"}))
}

export async function clear() {
  return fetch(`${baseUrl}/clear`, await withAuth({ method: "POST"}))
}