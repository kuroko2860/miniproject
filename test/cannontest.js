const autocannon = require("autocannon");
const { PassThrough } = require("stream");
const dumpRequestData = require("./dumpdata");

function run(url) {
  const buf = [];
  const outputStream = new PassThrough();

  const instance = autocannon({
    url,
    method: "POST",
    connections: 1000, // Number of concurrent connections
    pipelining: 100, // Number of requests per connection
    duration: 10, // Duration of test in seconds
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dumpRequestData),
  });

  autocannon.track(instance, { outputStream });

  outputStream.on("data", (data) => buf.push(data));
  instance.on("done", () => {
    process.stdout.write(Buffer.concat(buf));
  });
}

run("http://localhost:9000/objects");
