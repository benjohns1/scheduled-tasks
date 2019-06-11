import tasks from './_tasks.js';

export function get(req, res) {
  const { task } = req.params;
  
  if (task in tasks) {
    res.writeHead(200, {
			'Content-Type': 'application/json'
    });

    res.end(JSON.stringify(tasks[task]));
  } else {
    res.writeHead(404, {
			'Content-Type': 'application/json'
    });

    res.end(JSON.stringify({
      message: `Not found`
    }));
  }
}