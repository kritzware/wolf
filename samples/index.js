const Promise = require('bluebird')

const x = { y: 2, id: 3, }
console.log(x)

async function init() {
    for(let i = 0; i < 5; i++) {
        console.log('streamed message', i)
        await Promise.delay(1000)
    }
}

// init()