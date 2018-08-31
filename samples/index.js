const Promise = require('bluebird')

const x = { y: 2, id: 3, }
console.log(x)

async function init() {
    for(let i = 0; i < 5; i++) {
        console.log('streamed message', i)
        await Promise.delay(2000)
    }
}

// init()

// for(let i = 0; i < 10; i++) {
//     console.log(i)
// }
console.log("hello world!")