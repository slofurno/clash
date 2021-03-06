function request (options)
{ 
  if (!options) return;
  var body = options.body || "";
  var headers = options.headers || {};

  return new Promise((resolve, reject) => {
    var client = new XMLHttpRequest();
    client.onload = function (e) {
      if (client.status >= 300) {
        reject(client.statusText);
      } else {
        resolve(client.response);
      }
    };
    client.open (options.method, options.url);

    Object.keys(headers).forEach(x => {
      client.setRequestHeader(x, headers[x]);
    });        

          /*
    if (loggedin) {
      var account = JSON.parse(loggedin);
      var u = account.id + "." + account.token;
      client.setRequestHeader("Authorization", u);
    }
          */

    client.send (body);
  });
}

module.exports = request
