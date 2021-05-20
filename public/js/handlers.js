const proxyHandler = {
    get(target, p, receiver) {
        if (!target.hasOwnProperty(p)) {
            console.error(`[APP] Неизвестное свойство: ${p}`)
            return undefined
        }

        if (typeof target[p] === 'object' && target[p] !== null) {
            return new Proxy(target[p], proxyHandler)
        }

        return target[p]
    },
    set(target, p, value, receiver) {
        target[p] = value

        if (typeof handlers[`${p}Handler`] === 'function') {
            handlers[`${p}Handler`](value)
        }
    },
    deleteProperty(target, p) {
        if (!target.hasOwnProperty(p)) {
            console.error(`[APP] Неизвестное свойство: ${p}`)
            return false
        }

        delete target[p]
        return true
    },
}

window.handlers = {
    ShedulerHandler: (value = undefined) => {
        __App.addItem('123', 'day', 0)
    },
}