import fetch from 'node-fetch';

async function getAll() {
  const apiRes = await fetch(`http://localhost:8080/api/v1/task`);
  return apiRes.text();
}

async function get(id) {
  const apiRes = await fetch(`http://localhost:8080/api/v1/task/${id}`);
  return apiRes.text();
}

export { getAll, get }