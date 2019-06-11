import * as taskRepo from '../../data/task.repo';

export async function get(req, res) {
  const id = req.params.task;
  const taskRes = taskRepo.get(id);
  
  const contentType = {
    'Content-Type': 'application/json'
  };

  const task = await taskRes;
  
  if (task !== undefined) {
    res.writeHead(200, contentType);

    res.end(task);
  } else {
    res.writeHead(404, contentType);

    res.end(JSON.stringify({
      message: `Not found`
    }));
  }
}