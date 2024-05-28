const mongoose = require("../db");
const { Schema } = mongoose;

const requestSchema = new Schema({
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
  dump: { type: String },
  timestamp: { type: Date, default: Date.now },
});

// Tạo index cho trường location để hỗ trợ các truy vấn địa lý
// requestSchema.index({ location: "2dsphere" });

// Tạo index cho trường timestamp để hỗ trợ các truy vấn theo khoảng thời gian
requestSchema.index({ timestamp: 1 });

const Request = mongoose.model("Request", requestSchema);

module.exports = Request;
