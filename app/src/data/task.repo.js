import fetch from 'node-fetch';

const baseUrl = `http://localhost:8080/api/v1`;

export async function getAll() {
  return fetch(`${baseUrl}/task`).then(r => r.text());
}

export async function get(id) {
  return fetch(`${baseUrl}/task/${id}`).then(r => r.text());
}

export async function add(taskData) {
  return fetch(`${baseUrl}/task/`, { method: "POST", body: JSON.stringify(taskData)}).then(r => r.text());
}