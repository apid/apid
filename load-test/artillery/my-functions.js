//
// my-functions.js
//
module.exports = {
  captureDeployments: captureDeployments,
  fakeStatuses: fakeStatuses,
  randomUserKey: randomUserKey
}

function captureDeployments(requestParams, response, context, ee, next) {

  var array = JSON.parse(response.body);
  context.vars.captureDeployments = array;
  return next(); // MUST be called for the scenario to continue
}
/*
    requestParam are the parameters from this list
    https://github.com/request/request#requestoptions-callback
*/
function fakeStatuses(requestParams, context, ee, next) {
  //console.log("--------------------------");
  //console.log(context.vars.captureDeployments)
  var d = context.vars.captureDeployments;
  var statusArray = [];

  for (var i = 0; i< d.length; i++) {
        var status = {
            "id":d[i].id,
            "status": (i>(d.length*0.1))?"FAIL":"SUCCESS",
            "message": "Some random message long.Some random message long.Some random message long.Some random message long.Some random message long.Some random message long.Some random message long.",
            "errorCode":1
        };
        statusArray.push(status);
  }
  requestParams.body = statusArray;
  requestParams.json = true;
  //console.log(requestParams.body);
  return next();// MUST be called for the scenario to continue
}

function randomUserKey(requestParams, context, ee, next) {
    numDevs = 50000
    key = Math.floor((Math.random() * numDevs) + 1)

    requestParams.body = requestParams.body + "&key=" + key
    return next()
}