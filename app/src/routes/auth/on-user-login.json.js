import { proxy, error } from '../../api/proxy'

export function post(req, res) {
  const provider = req.body && req.body.iss ? req.body.iss : null
  const userId = req.body && req.body.sub ? req.body.sub : null
  if (!provider || !userId) {
    error(res, 400, `provider and user ID could not be parsed from iss and sub fields: ${JSON.stringify(req.body)}`)
    return
  }
  const user = {
    displayname: req.body.displayname
  }
  proxy(req, res, { method: 'PUT', url: `/user/external/${provider}/${userId}/addOrUpdate`, body: user } )
}