import { message } from 'ant-design-vue';

export const URL = {
    Login: '/api/login.php',
    Forward: '/api/forward.php',
    Forward_DEL: '/api/forward_del.php'
}

export const POST = (url, body = {}) => {
    const loading = message.loading('正在操作中', 0);
    return new Promise((resolve) => {
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body)
        }).then(data => data.json()).then(res => {
            loading()
            if (res.status == 0) {
                message.error(res.message)
            }
            resolve(res)
        }).catch((err) => {
            resolve({ status: 0, msg: err, data: null })
        })
    })
}
export const GET = (url, query = {}) => {
    const params = new URLSearchParams(Object.assign({ page_index: 1, page_size: 10 }, query));
    return new Promise((resolve) => {
        fetch(`${url}?${params.toString()}`, {
            method: 'GET'
        }).then(data => data.json()).then(res => {
            if (res.status == 0) {
                message.error(res.message)
            }
            resolve(res)
        }).catch((err) => {
            resolve({ status: 0, msg: err, data: null })
        })
    })
}

export function formatBytes(bytes) {
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let i = 0;
    while (bytes >= 1024 && i < units.length - 1) {
        bytes /= 1024;
        i++;
    }
    return `${bytes.toFixed(2)} ${units[i]}`;
}