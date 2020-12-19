const Mock = require("mockjs")
const fs = require("fs")
const Random = Mock.Random

function genCatetory() {
    const CATORYIES = ["修真", "玄幻", "灵异", "穿越", "言情", "网游", "都市", "科幻", "武侠"]
    return CATORYIES[Math.floor(Math.random() * CATORYIES.length)]
}

function genBooksData(size) {
    let res = []
    for (let i = 0; i < size; i++) {
        res.push(Mock.mock({
            author: Random.cname(),
            title: Random.ctitle(5, 20),
            description: Random.csentence(30, 80),
            pub_date: Random.date("yyyy-MM-dd HH:mm:ss"),
            category: genCatetory()
        }))
    }
    return res
}

const books = genBooksData(100)
const file = "./books.json"

fs.writeFile(file, JSON.stringify(books), (err) => {
    if (err) {
        console.error("generate books data error: %s", err)
    } else {
        console.log("done")
    }
})

