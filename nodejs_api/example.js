const {request, client, constants} = require('./index.js')

client.on(constants.connected, ()=> {
  request(constants.instrumentsList, {sec_type: 'SPOT'})
})

client.on(constants.instrumentsList, (msg)=>{
  console.log(msg)
})