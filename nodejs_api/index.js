const Index = require('ws')
const EventEmitter = require('events')
const CryptoJS = require('crypto-js')
const client = new EventEmitter()

const constants = {
  connected: 'connected',
  hb: 'hb',
  instrumentsList: 'inst_list',
  orderBook: 'inst_order_book',
  orderBookUpdate: 'inst_order_book_update',
  login: 'login'
}

let hbUpdatedAt
let hbInterval
const currentRequests = []

const generateNonce = (offset = 0) => {
  const nonce = (Date.now() % 429496729) + offset
  if (!currentRequests[nonce]) {
    return nonce
  }
  return generateNonce(offset + 1)
}

const ws = new Index('wss://wsapi.coinut.com', 'beta')

const request = (type, payload = {}, shouldTrack = true) => {
  payload.request = type
  payload.nonce = payload.nonce || generateNonce()
  if (shouldTrack) {
    if (payload.request === constants.hb && Object.keys(currentRequests).find(nonce => currentRequests[nonce].request === constants.hb)) {
      return
    }
    currentRequests[payload.nonce] = Object.assign({requestAt: Date.now()}, payload)
  }
  ws.send(JSON.stringify(payload))
  return payload.nonce
}

ws.onopen = () => {
  clearInterval(hbInterval)
  hbInterval = setInterval(() => {
    if (Date.now() - hbUpdatedAt > 5000) {
      console.error('Server is unreachable now.')
    }
  }, 1000)
  request(constants.hb)
  client.emit(constants.connected)

  const timestamp = Math.floor(new Date() / 1000)
  const nonce = generateNonce()
  request(constants.login, {
    nonce,
    username: '{username}',
    timestamp,
    hmac_sha256: CryptoJS.HmacSHA256('{username}' + '|' + timestamp + '|' + nonce, '{api key}').toString(CryptoJS.enc.Hex)
  })
}

ws.onmessage = e => {
  try {
    const msg = JSON.parse(e.data)
    const fullMsg = Object.assign(msg, currentRequests[msg.nonce] || {})
    delete currentRequests[msg.nonce]

    switch (fullMsg.reply) {
      case constants.hb:
        // keep rounding heart beats
        hbUpdatedAt = Date.now()
        setTimeout(() => request(constants.hb), 1000)
        break
      case constants.instrumentsList:
        client.emit(constants.instrumentsList, msg)
        break
      case constants.orderBook:
        client.emit(constants.orderBook, msg)
        break
      case constants.orderBookUpdate:
        client.emit(constants.orderBookUpdate, msg)
        break
      case constants.login:
        client.emit(constants.login)
        break
      default:
        if (msg.msg !== 'user_id field is not allowed') {
          console.log(JSON.stringify(msg, null, 3))
        }
    }
  } catch (err) {
    console.error(err)
  }
}

ws.onclose = () => {
  console.warn('disconnected')
}

module.exports = {
  request,
  constants,
  client,
}
