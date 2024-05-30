const autocannon = require("autocannon");
const dumpRequestData = require("./dumpdata");

function run(url) {
  const instance = autocannon({
    url,
    method: "POST",
    connections: 100, // Number of concurrent connections
    pipelining: 100, // Number of requests per connection
    duration: 10, // Duration of test in seconds
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dumpRequestData),
  });

  autocannon.track(instance);
}

run("http://localhost:8080/objects");
