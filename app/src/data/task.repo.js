import fetch from 'node-fetch';

const baseUrl = `http://${process.env.APPLICATION_HOST || 'localhost'}:${process.env.APPLICATION_PORT || '8080'}/api/v1/task`;

export async function getAll() {
  return fetch(`${baseUrl}`);
}

export async function add(taskData) {
  return fetch(`${baseUrl}/`, { method: "POST", body: JSON.stringify(taskData)});
}

export async function complete(id) {
  return fetch(`${baseUrl}/${id}/complete`, { method: "PUT"});
}

export async function clear() {
  return fetch(`${baseUrl}/clear`, { method: "POST"});
}