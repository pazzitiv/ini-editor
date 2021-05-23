class Api {
    get base() {
        return '/api/'
    }

    prepareData = (data) => {
        if (data) {
            switch (typeof data) {
                case "object":
                    if(Array.isArray(data)) {
                        return data.map((item, index) => `data[${index}]=${item}`).join('&')
                    } else {
                        let res = '';
                        for(let s of Object.entries(data)) {
                            res += `${s[0]}=${s[1]}`
                        }
                        return res
                    }
                default:
                    return `data=${data}`
            }
        }

        return data
    }

    get(method, data = null) {
        return fetch(this.base + method + (this.prepareData(data) ?? ''), {
            method: 'GET',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            referrerPolicy: 'no-referrer'
        }).then(response => response.json())
    }

    post(method, data) {
        return fetch(this.base + method, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            referrerPolicy: 'no-referrer',
            body: JSON.stringify(data)
        }).then(response => response.json())
    }

    delete(method, id = null, data = null) {
        return fetch(this.base + method + `${id ? '/' + id : ''}` + (data ? '?' + this.prepareData(data) : ''), {
            method: 'DELETE',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            referrerPolicy: 'no-referrer'
        })
    }

    serializeObject(el) {
        let obj = {};
        const arr = el.serializeArray();

        arr.forEach((item) => {
            if(obj.hasOwnProperty(item.name)) {
                obj[item.name] = [
                    obj[item.name],
                ]
            }

            if(Array.isArray(obj[item.name])) {
                obj[item.name].push(item.value)
            } else {
                obj[item.name] = item.value
            }
        })

        return obj
    }
}

window.__Api = new Api()