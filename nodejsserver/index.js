const express = require("express");
const app = express();
const bodyParser = require("body-parser");
const { Worker } = require("worker_threads");
const path = require("path");

const requestQueue = [];
const BATCH_SIZE = 1000;
const FLUSH_INTERVAL = 1000; // in ms
const worker = new Worker(path.resolve(__dirname, "./worker.js"));
const mongoWorker = new Worker(
  path.resolve(__dirname, "./redisToMongoWorker.js")
);

mongoWorker.on("message", (message) => {
  if (message.success) {
    console.log(`Transferred ${message.count} records from Redis to MongoDB`);
  } else {
    console.error(
      "Transfer from Redis to MongoDB failed:",
      message.message || message.error
    );
  }
});

worker.on("message", (message) => {
  if (!message.success) {
    console.error("Batch insert failed:", message.error);
  }
});

// Function to flush the queue
const flushQueue = () => {
  if (requestQueue.length === 0) return;

  const batch = requestQueue.splice(0, BATCH_SIZE);
  worker.postMessage(batch);
};
// Set interval to flush the queue
setInterval(flushQueue, FLUSH_INTERVAL);

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

app.post("/request", async (req, res) => {
  const requestData = req.body;
  requestQueue.push(requestData);

  if (requestQueue.length >= BATCH_SIZE) {
    flushQueue();
  }
  res.send("ok");
});

app.get("/", (req, res) => {
  res.send("Hello World!");
});

app.listen(8080, () => {
  console.log("Server is running at http://localhost:8080/");
});
