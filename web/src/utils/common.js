export const Login_URL = "/api/login.php"

export const Api = (url, body) => {
    return new Promise((resolve) => {
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body)
        }).then(data => data.json()).then(resolve).catch((err) => {
            resolve({ status: 0, msg: err, data: null })
        })
    })
}
