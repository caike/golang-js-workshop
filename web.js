const express = require("express");
const web = express();

web.set("env", process.env.NODE_ENV || "development");
web.use(express.static("public"));

const bodyParser = require("body-parser");
const jsonParser = bodyParser.json();
/*
 * Receives data from device
 */

// TODO: check token
web.post("/data", jsonParser, (req, res) => {
  console.log("receiving data from device: ", req.headers["x-device-name"]);
  // do something with device name
  const deviceName = req.headers["x-device-name"];
  // TODO: broadcast data to all registered web clients.
  console.log("wat: ", web.locals.wsClients);
  web.locals.wsClients.forEach((client) => {
    client.send(JSON.stringify(req.body));
  });
  res.sendStatus(201)
});

module.exports = web;
