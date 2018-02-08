"use strict";

const http = require("http");
const server = http.createServer();
const web = require("./web");

// Dispatches web requests
// to Express.
server.on("request", web);

// WebSocket code
//
const WSServer = require("ws").Server;
const webWs = new WSServer({ noServer: true });
web.locals.wsClients = webWs.clients;
const url = require("url");

server.on("upgrade", (request, socket, head) => {
  const pathname = url.parse(request.url).pathname;
  if (pathname === "/web") {
    webWs.handleUpgrade(request, socket, head, (ws) => {
      webWs.emit("connection", ws);
    });
  } else {
    socket.destroy();
  }
});

webWs.on("connection", (socket) => {
  console.log("Received connection from web");
  console.log("Current clients connected: " + webWs.clients.size);

  socket.on("close", () => {
    //web.locals.webClients.delete(socket);
    console.log("closing connection from web");
  });

  socket.on("error", () => {
    //web.locals.webClients.delete(socket);
    console.log("closing connection from web");
  });
});


const PORT = process.env.PORT || 3000;
server.listen(PORT, () => console.log(`Listening on ${PORT}`));

