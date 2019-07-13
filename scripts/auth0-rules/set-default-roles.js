/**
 * !!! This rule was setup to run on Auth0, but does not actually set roles during user creation
 * !!! Instead, golang api is configured to set default permissions without handling them via Auth0
 * @param {*} user 
 * @param {*} context 
 * @param {*} callback 
 */
export const setDefaultRoles = function (user, context, callback) {
  // This rule is only for Auth0 databases
  if (context.connectionStrategy !== 'auth0') {
    return callback(null, user, context);
  }
  
  const roles = ['user'];
  
  // Don't apply rule if all roles already exist
  if (context &&
      context.authorization &&
      context.authorization.roles) {
    for (var r in roles) {
      if (!context.authorization.roles.includes(roles[r])) {
        return callback(null, user, context);
      }
    }
  }
  
  const request = require('request');
  
  (function chain(chain) {
    const next = chain.shift();
    next(chain);
  })([
    getManagementToken,
    getUserRoles,
    getRoleIds,
    setRoles
  ]);
    
  function getManagementToken(chain) {
    request.post({
      url: `https://${auth0.domain}/oauth/token`,
      headers: { 'content-type': 'application/x-www-form-urlencoded' },
      form: {
        grant_type: 'client_credentials',
        client_id: configuration.roleSetClientId,
        client_secret: configuration.roleSetClientSecret,
        audience: `https://${auth0.domain}/api/v2/`
      }
    }, function (err, response, body) {
      console.log('getManagementToken:', body);
      if (err) throw new Error(err);
      const accessToken = JSON.parse(body).access_token;
      
      const next = chain.shift();
      next(accessToken, chain);
    });
  }
  
  function getUserRoles(accessToken, chain) {
    const headers = {
      authorization: 'Bearer ' + accessToken,
      'content-type': 'application/json'
    };
    request.get({
      url: `${auth0.baseUrl}/users/${user.user_id}/roles`,
      headers: headers
    }, function(err, response, body) {
      console.log('getUserRoles:', body);
      if (err) throw new Error(err);
      const userRoles = JSON.parse(body).map(r => r.name);
      let updateRoles = false;
      for (r in roles) {
        if (!userRoles.includes(roles[r])) {
          updateRoles = true;
          break;
        }
      }
      if (updateRoles) {
        const next = chain.shift();
        next(headers, chain);
      } else {
        return callback(null, user, context);
      }
    });
  }
  
  function getRoleIds(headers, chain) {
    request.get({
      url: `${auth0.baseUrl}/roles`,
      qs: { name_filter: 'user' },
      headers: headers
    }, function(err, response, body) {
      console.log('getRoleIds:', body);
      if (err) throw new Error(err);
      const roleIds = JSON.parse(body).reduce((ids, r) => {
        if (roles.includes(r.name)) {
          ids.push(r.id);
        }
        return ids;
      }, []);
      const next = chain.shift();
      next(headers, roleIds, chain);
    });
  }
  
  function setRoles(headers, roleIds, chain) {
    request.post({
      url: `${auth0.baseUrl}/users/${user.user_id}/roles`,
      headers: headers,
      json: {
        "roles": roleIds
      }
    },
    function(err, response, body) {
      if (err) throw new Error(err);
      if (response.statusCode !== 204) throw new Error(`Unexpected status code ${response.statusCode}: ${JSON.stringify(body)}`);
      console.log('default user roles successfully set');
      return callback(null, user, context);
    });
  }
}