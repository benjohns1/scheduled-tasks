import sirv from 'sirv'
import polka from 'polka'
import compression from 'compression'
import { json } from 'body-parser'
import cookieParser from 'cookie-parser'
import * as sapper from '@sapper/server'

const { PORT, NODE_ENV } = process.env
const dev = NODE_ENV === 'development'

polka() // You can also use Express
	.use(
		compression({ threshold: 0 }),
		sirv('static', { dev }),
		json(),
		cookieParser(),
		sapper.middleware({
			session: (req) => {
				if (req.cookies.token) {
					return ({
						auth: {
							token: req.cookies.token,
							devLogin: req.cookies.devLogin
						}
					})
				}
			}
		})
	)
	.listen(PORT, err => {
		if (err) console.log('error', err);
	});
