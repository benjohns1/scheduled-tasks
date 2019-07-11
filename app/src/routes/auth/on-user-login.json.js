import { proxy, error } from '../../api/proxy'

export function post(req, res) {
  const idData = req.body.sub.split('|', 2)
  const provider = idData[0]
  const userId = idData[1]
  if (!provider || !userId) {
    error(res, 400, `provider and user ID could not be parsed from 'sub' field: ${JSON.stringify(req.body)}`)
    return
  }
  const user = {
    displayname: req.body.displayname
  }
  proxy(req, res, { method: 'PUT', url: `/user/external/${provider}/${userId}/addOrUpdate`, body: user } )
}