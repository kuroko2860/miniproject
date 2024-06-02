const { MongoClient } = require("mongodb");
const { Pool } = require("pg");
const pool = new Pool({
  user: "kuroko",
  host: "postgres",
  database: "miniproject",
  password: "123456",
  port: 5432,
});

let batch = [];
const BATCH_SIZE = 1000;
const BATCH_INTERVAL = 1000; // 1 second

async function addToBatch(document) {
  batch.push(document);
  if (batch.length >= BATCH_SIZE) {
    await processBatch();
  }
}

async function processBatch() {
  if (batch.length === 0) return;

  const client = await pool.connect();

  try {
    await client.query("BEGIN");
    const insertPromises = batch.map((doc) => {
      const text = `
        INSERT INTO objects (object_id, type, color, location, status, timestamp)
        VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7)
      `;
      const values = [
        doc.id,
        doc.type,
        doc.color,
        doc.location.coordinates[0],
        doc.location.coordinates[1],
        doc.status,
        new Date(doc.timestamp),
      ];
      return client.query(text, values);
    });

    await Promise.all(insertPromises);
    await client.query("COMMIT");

    console.log(`Processed batch of ${batch.length} records`);
    batch = [];
  } catch (e) {
    await client.query("ROLLBACK");
    console.error("Error processing batch:", e);
  } finally {
    client.release();
  }
}

setInterval(processBatch, BATCH_INTERVAL);

async function watchCollection() {
  const client = new MongoClient("mongodb://mongodb:27017");
  await client.connect();

  const db = client.db("miniproject");
  const collection = db.collection("objects");

  const changeStream = collection.watch();

  changeStream.on("change", async (change) => {
    if (change.operationType === "insert") {
      const newDocument = change.fullDocument;
      await addToBatch(newDocument);
    }
  });

  console.log("Watching for changes...");
}

watchCollection().catch(console.error);
