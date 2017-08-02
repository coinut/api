const request = require('request')

const options = process.env.EMAIL && process.env.ACCESS_TOKEN ? {
  // http basic authentication
  auth: {
    user: process.env.EMAIL,
    pass: process.env.ACCESS_TOKEN
  }
} : {}

// convert 3200 usd to xxx btc
let amount = '3200'
let currency = 'usd'
request.get(`https://bitcoinindex.io/btcusd/convert/${amount}/${currency}`, options, (err, response, body) => {
  console.log(JSON.parse(body))
})

// convert 1.239 btc to xxx usd
amount = '1.239'
currency = 'btc'
request.get(`https://bitcoinindex.io/btcusd/convert/${amount}/${currency}`, options, (err, response, body) => {
  console.log(JSON.parse(body))
})
