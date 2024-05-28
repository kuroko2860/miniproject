// redisToMongoWorker.js
const { parentPort } = require("worker_threads");
const Redis = require("ioredis");
const Request = require("./model/request"); // Import the Request model

const redis = new Redis({ port: 6379, host: "localhost" });

const transferData = async () => {
  try {
    const dataBatch = await redis.lrange("requests", 0, 999); // Get 1000 items
    if (dataBatch.length > 0) {
      await redis.ltrim("requests", dataBatch.length, -1); // Remove items from list

      const records = dataBatch.map((item) => JSON.parse(item));
      await Request.insertMany(records);
      parentPort.postMessage({ success: true, count: records.length });
    } else {
      parentPort.postMessage({
        success: false,
        message: "No data to transfer",
      });
    }
  } catch (error) {
    parentPort.postMessage({ success: false, error: error.message });
  }
};

parentPort.on("message", async () => {
  await transferData();
});

// Periodically transfer data every second
setInterval(() => {
  parentPort.postMessage("transfer");
}, 1000);
