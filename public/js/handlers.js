const proxyHandler = {
    get(target, p, receiver) {
        if(typeof target[p] === "function") {
            return target[p]
        }

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
    getOwnPropertyDescriptor(target, property) {
        if (target.hasOwnProperty(property)) {
            return {configurable: true, enumerable: true};
        }
    }
}

window.handlers = {
    optionsHandler: (value) => {
        for(let option in value) {
            switch (option) {
                case 'Scheduler':
                    $(`#day-list`).innerHTML = '';
                    $(`#time-list`).innerHTML = '';

                    value.Scheduler.Day.forEach((item, index) => {
                        let days = item.map((day, i) => {
                            if(day === "1") {
                                return __App.enums.days[i]
                            }
                            return ''
                        }).filter(d => d !== '')

                        __App.addItem(days.join(','),'day',index)
                    })

                    value.Scheduler.Time.forEach((item, index) => {
                        __App.addItem(item,'time',index)
                    })
                    break;
                case 'Templates':
                    $(`#sender-list`).innerHTML = '';
                    $(`#subject-list`).innerHTML = '';

                    value.Templates.Sender.forEach((item, index) => {
                        __App.addItem(item,'sender',index)
                    })
                    value.Templates.Subject.forEach((item, index) => {
                        __App.addItem(item,'subject',index)
                    })
                    break;
                case 'Destinations':
                    $(`#phone-list`).innerHTML = '';

                    value.Destinations.Tel_num.forEach((item, index) => {
                        __App.addItem(item,'phone',index)
                    })
                    break;
            }
        }

        toggleItems()
    },
    schedulesHandler: (value = undefined) => {
        const schedules = value.Map_id;

        $('#scheduler table tbody').html('')

        schedules.forEach((item, ind) => {
            __App.addSchedule(item, ind + 1)
        })
    },
    systemHandler: (value = undefined) => {
        $('#sys-period').val(value.System.Period);
        __App.enums.log_levels.forEach((item) => {
            $('#sys-logging').append(`<option ${item == value.System.Logging ? 'selected' : ''}>${item}</option>`)
        });

        $('#sys-server').val(value.Imap.Server);
        $('#sys-login').val(value.Imap.Login);
        $('#sys-folder-check').val(value.Imap.Folder_check);
        if (value.Imap.Trash === 1) {
            $('#sys-trash').attr('checked');
        } else {
            $('#sys-trash').removeAttr('checked');
        }

        $('#sys-ast-host').val(value.Asterisk.Host);
        $('#sys-ast-port').val(value.Asterisk.Port);
    }
}