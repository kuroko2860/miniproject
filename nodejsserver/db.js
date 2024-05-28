const mongoose = require("mongoose");

mongoose
  .connect("mongodb://127.0.0.1:27017/miniproject")
  .then(() => console.log("db connected !"));

module.exports = mongoose;
