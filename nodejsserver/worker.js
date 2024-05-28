// worker.js
const { parentPort } = require("worker_threads");
const Redis = require("ioredis");
const redis = new Redis({ port: 6379, host: "localhost" });

parentPort.on("message", async (dataBatch) => {
  const pipeline = redis.pipeline();
  dataBatch.forEach((data) => {
    pipeline.rpush("requests", JSON.stringify(data));
  });

  try {
    await pipeline.exec();
    parentPort.postMessage({ success: true });
  } catch (error) {
    parentPort.postMessage({ success: false, error: error.message });
  }
});
