import fetch from 'node-fetch';

async function getAll() {
  return fetch(`http://localhost:8080/api/v1/task`).then(r => r.text());
}

async function get(id) {
  return fetch(`http://localhost:8080/api/v1/task/${id}`).then(r => r.text());
}

export { getAll, get }