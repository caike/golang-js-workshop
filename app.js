const WSServer = require('ws').Server;
const http = require("http");
const server = http.createServer();
const util = require("util");
const web = require("./web");
const url = require("url");

server.on('request', web);

// Create two web socket servers to run
// on top of a single http server
// from: https://github.com/websockets/ws/pull/885
let piWs = new WSServer({ noServer: true });
let webWs = new WSServer({ noServer: true });

server.on("upgrade", (request, socket, head) => {
	const pathname = url.parse(request.url).pathname;
	if (pathname === "/pi") {
		piWs.handleUpgrade(request, socket, head, (ws) => {
			piWs.emit("connection", ws);
		});
	} else if (pathname.startsWith("/web")) {
		webWs.handleUpgrade(request, socket, head, (ws) => {
			webWs.emit("connection", ws);
		});
	} else {
		socket.destroy();
	}
});

const webClients = new Map();
let raspberrySocket;

piWs.on("connection", (socket) => {
  console.log("connection from websocket");
  // currently supporting a SINGLE connection
  // from a raspberryPi
  raspberrySocket = socket;
  // receives result of running command on the PI
  raspberrySocket.on("message", (raspResponse) => {
    console.log("receipved message from raspberry");
    let [clientName, output] = raspResponse.split("::");
    let websocketClient = webClients.get(clientName);
    websocketClient.send(output);
  });

  raspberrySocket.on("close", () => console.log("Connection closed"));
});

webWs.on("connection", (browserSocket) => {
  let clientName = browserSocket.upgradeReq.url.split("-")[1]
  console.log("Connection from web client: ", clientName);
  webClients.set(clientName, browserSocket);

  browserSocket.on("message", (commandFromWeb) => {
    raspberrySocket.send(`${clientName}::${commandFromWeb}`);
  });
  browserSocket.on("close", () => {
    console.log("closing connection from web");
    webClients.splice(webClients.indexOf(browserSocket), 1);
  });
});

server.listen(3000, () => console.log("Listening web and ws"));
