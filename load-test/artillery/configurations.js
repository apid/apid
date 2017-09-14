// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
//
module.exports = {
  captureConfigurations: captureConfigurations,
  //fakeStatuses: fakeStatuses,
  randomUserKey: randomUserKey
}

function captureConfigurations(requestParams, response, context, ee, next) {
    console.log("--------------------------");
    var d = JSON.parse(response.body);
    var statusArray = [];
    for (var i = 0; i< d.length; i++) {
        var status = {
            "kind": d[i].kind,
			//"self": d[i].self,
			"contents": d[i].contents,
        };
        statusArray.push(status);
    }
    console.log(statusArray);

	return next(); // MUST be called for the scenario to continue
}
/*
 requestParam are the parameters from this list
 https://github.com/request/request#requestoptions-callback
 */
/*
function fakeStatuses(requestParams, context, ee, next) {
    var d = context.vars.captureConfigurations;
requestParams.body = statusArray;
    requestParams.json = true;
    //console.log(requestParams.body);
    return next();// MUST be called for the scenario to continue
}
*/
function randomUserKey(requestParams, context, ee, next) {
    numDevs = 50000
    key = Math.floor((Math.random() * numDevs) + 1)

    requestParams.body = requestParams.body + "&key=" + key
    return next()
}
