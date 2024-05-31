wrk.method = "POST"
wrk.headers["content-type"] = "application/json"
wrk.body = '{"id":"1","type":"car","color":"red","location":{"type":"Point","coordinates":[106.660172,10.762622]},"status":"moving","dump":"123456"}'