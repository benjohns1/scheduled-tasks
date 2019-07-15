// Load env vars from file, if ENV_FILEPATH was set
if (process.env.ENV_FILEPATH) {
	require('dotenv').config({ path: process.env.ENV_FILEPATH })
}