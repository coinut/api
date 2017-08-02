const WebSocket = require('ws');

const ws = new WebSocket('wss://bitcoinindex.io')

ws.on('open', ()=> {
  // if you have a valid access token, you can send below message to get real-time index tick updates.
  // However, if you don't have a token, you can still get a index tick update every second.
  if (process.env.EMAIL && process.env.ACCESS_TOKEN) {
    ws.send(JSON.stringify({
      request: 'auth',
      email: process.env.EMAIL,
      token: process.env.ACCESS_TOKEN
    }))
  }
})

ws.on('message', data => {
  // payload format
  // {
  //   "status": ["OK"],
  //   "asset":"XBTCUSD",
  //   "index":"2712.43874451",
  //   "timestamp":1501727879179
  // }
  console.log(data)
})

ws.on('error', err=>{
  console.error(err)
})
