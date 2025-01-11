// JS Example for subscribing to a channel
/* eslint-disable */
const WebSocket = require("ws");
const { sign } = require("jsonwebtoken");
const crypto = require("crypto");
const fs = require("fs");
const { exit } = require("process");
require("dotenv").config();

// Derived from your Coinbase CDP API Key
//  SIGNING_KEY: the signing key provided as a part of your API key. Also called the "SECRET KEY"
//  API_KEY: the api key provided as a part of your API key. also called the "API KEY NAME"
const API_KEY = process.env.API_KEY;
const SIGNING_KEY = process.env.SIGNING_KEY;

const algorithm = "ES256";

const CHANNEL_NAMES = {
  level2: "level2",
  user: "user",
  tickers: "ticker",
  ticker_batch: "ticker_batch",
  status: "status",
  market_trades: "market_trades",
  candles: "candles",
};

// The base URL of the API
const WS_API_URL = "wss://advanced-trade-ws.coinbase.com";

function signWithJWT(message, channel, products = []) {
  const jwt = sign(
    {
      iss: "cdp",
      nbf: Math.floor(Date.now() / 1000),
      exp: Math.floor(Date.now() / 1000) + 120,
      sub: API_KEY,
    },
    SIGNING_KEY,
    {
      algorithm,
      header: {
        kid: API_KEY,
        nonce: crypto.randomBytes(16).toString("hex"),
      },
    }
  );

  return { ...message, jwt: jwt };
}

function subscribeToProducts(products, channelName, ws) {
  const message = {
    type: "subscribe",
    channel: channelName,
    product_ids: products,
  };
  //   const subscribeMsg = signWithJWT(message, channelName, products);
  ws.send(JSON.stringify(message));
}

function unsubscribeToProducts(products, channelName, ws) {
  const message = {
    type: "unsubscribe",
    channel: channelName,
    product_ids: products,
  };
  ws.send(JSON.stringify(message));
}

function onMessage(data) {
  const parsedData = JSON.parse(data);
  fs.appendFile("Output1.txt", data, (err) => {
    // In case of a error throw err.
    if (err) throw err;
  });
}

const connections = [];
let sentUnsub = false;
for (let i = 0; i < 1; i++) {
  const date1 = new Date(new Date().toUTCString());
  const ws = new WebSocket(WS_API_URL);

  ws.on("message", function (data) {
    const date2 = new Date(new Date().toUTCString());
    const diffTime = Math.abs(date2 - date1);
    if (diffTime > 5000 && !sentUnsub) {
      unsubscribeToProducts(["BTC-USD"], CHANNEL_NAMES.level2, ws);
      sentUnsub = true;
      console.log("Unsubscribed");
      exit();
    }

    const parsedData = JSON.parse(data);
    fs.appendFile("Output1.txt", data, (err) => {
      // In case of a error throw err.
      if (err) throw err;
    });
    fs.appendFile("Output1.txt", "\n", (err) => {
      // In case of a error throw err.
      if (err) throw err;
    });
  });

  ws.on("open", function () {
    const products = ["BTC-USD", "ETH-USD", "LTC-USD", "BCH-USD"];
    subscribeToProducts(products, CHANNEL_NAMES.level2, ws);
  });

  connections.push(ws);
}
