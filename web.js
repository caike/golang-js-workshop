const express = require("express");
const web = express();

web.get("/", (req, res) => {
  res.sendFile(__dirname + "/simple-index.html");
});

module.exports = web
