const mongoose = require("mongoose");
const Redis = require("ioredis");
const { Schema } = mongoose;

// Define the Object schema
const objectSchema = new Schema({
  id: { type: String, required: true },
  type: {
    type: String,
    enum: ["car", "bike", "pedestrian", "truck"],
    required: true,
  },
  color: {
    type: String,
    enum: ["red", "blue", "yellow", "black", "white"],
    required: true,
  },
  location: {
    type: { type: String, default: "Point" },
    coordinates: { type: [Number], required: true },
  },
  status: { type: String, enum: ["stationary", "moving"], required: true },
  timestamp: { type: Date, default: Date.now },
  dump: { type: String },
});

// Create the Mongoose model
const ObjectModel = mongoose.model("Object", objectSchema);

// Connect to MongoDB
mongoose.connect(
  process.env.MONGO_URI || "mongodb://mongodb:27017/miniproject"
);

// Initialize Redis client
const redis = new Redis({
  host: process.env.REDIS_HOST || "localhost",
  port: process.env.REDIS_PORT || 6379,
});

// Function to process a batch of records from Redis and insert them into MongoDB
async function processBatch() {
  const batchSize = 1000; // Number of records to process per batch
  try {
    // Get a batch of records from Redis
    const records = await redis.lrange("objects", 0, batchSize - 1);
    if (records.length === 0) return; // No records to process

    // Parse the records in parallel using Promise.all
    const objects = await Promise.all(
      records.map((record) => JSON.parse(record))
    );

    // Insert the records into MongoDB
    await ObjectModel.insertMany(objects, { ordered: false });

    // Remove the processed records from Redis
    await redis.ltrim("objects", batchSize, -1);

    console.log(`Processed ${objects.length} records`);
  } catch (error) {
    console.error("Error processing batch:", error);
  }
}

// Run the worker in an infinite loop, processing a batch every second
(async function () {
  while (true) {
    await processBatch();
    await new Promise((resolve) => setTimeout(resolve, 1000)); // Sleep for 1 second
  }
})();
