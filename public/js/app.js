const __App = new Proxy({
    addItem: (value, dataType, index, type = 'checkbox') => {
        $(`#${dataType}-list`).append(`<div className="form-check">
                <input className="form-check-input" name="${dataType}[]" type="${type}" data-id="${index + 1}" value="" id="${dataType}-${index}">
                    <label className="form-check-label" htmlFor="${dataType}-${index}">
                        ${value}
                    </label>
            </div>`)
    },
    addSchedule: (schedule, num) => {
        let tb = $('#scheduler table tbody')
        let row = document.createElement('tr')
        schedule.forEach((item, index) => {
            let cell = document.createElement('td')

            switch (index) {
                case 0:
                    if (__App.has(__App.options, 'Scheduler')) {
                        $(`#day-${parseInt(item) - 1}`).attr('disabled', true)
                        let days = __App.options.Scheduler.Day[parseInt(item) - 1].map((day, i) => {
                            if (day === "1") {
                                return __App.enums.days[i]
                            }
                            return ''
                        }).filter(d => d !== '')

                        cell.innerHTML = days.join(', ')
                    }
                    break
                case 1:
                    if (__App.has(__App.options, 'Scheduler')) {
                        $(`#time-${parseInt(item) - 1}`).attr('disabled', true)
                        cell.innerHTML = __App.options.Scheduler.Time[parseInt(item) - 1]
                    }
                    break
                case 2:
                    if (__App.has(__App.options, 'Templates')) {
                        $(`#sender-${parseInt(item) - 1}`).attr('disabled', true)
                        cell.innerHTML = __App.options.Templates.Sender[parseInt(item) - 1]
                    }
                    break
                case 3:
                    if (__App.has(__App.options, 'Templates')) {
                        $(`#subject-${item}`).attr('disabled', true)
                        cell.innerHTML = __App.options.Templates.Subject[parseInt(item) - 1]
                    }
                    break
                case 4:
                    if (__App.has(__App.options, 'Destinations')) {
                        $(`#phone-${parseInt(item) - 1}`).attr('disabled', true)
                        cell.innerHTML = __App.options.Destinations.Tel_num[parseInt(item) - 1]
                    }
                    break
                case 5:
                    cell.innerHTML =
                        `<div className="form-check">
                            <input className="form-check-input" type="checkbox" value="${parseInt(item)}" ${parseInt(item) === 1 ? 'checked' : ''} id="task-${num}">
                        </div>`
                    break
            }

            row.append(cell)
        })

        let cell = document.createElement('td')
        cell.innerHTML = `<td><i class="fas fa-trash" id="delete-task-${num}"></i></td>`
        row.append(cell)

        tb.append(row)
    },
    has: (object, prop) => {
        const hOP = Object.prototype.hasOwnProperty
        return hOP.call(object, prop)
    },
    enums: {
        days: ["????", "????", "????", "????", "????", "????", "????"],
        log_levels: [
            'INFO',
            'DEBUG',
        ],
    },
    options: {},
    schedules: {},
    system: {},
}, proxyHandler)

$(document).ready(() => {
    $('.timepicker').timepicker({
        timeFormat: 'HH:mm',
        interval: 60,
        minTime: '0',
        maxTime: '23',
        defaultTime: '8',
        startTime: '0',
        dynamic: false,
        dropdown: true,
        scrollbar: true
    });

    __Api.get('dictionaries')
        .then((response) => {
            __App.options = response.data

            __Api.get('schedules')
                .then((response) => {
                    __App.schedules = response.data.Map
                })
        })
})

window.toggleItems = () => {
    $('#reference').on('change', 'input[type="checkbox"]', function (e) {
        e.preventDefault()

        const card = $(this).parents('.card');
        const checked = card.find('input[type="checkbox"]:checked')

        if (checked.length > 0) {
            card.find('[data-action="delete"]').removeAttr('disabled')
        } else {
            card.find('[data-action="delete"]').attr('disabled', true)
        }

        return false
    })
}

$('[data-action="add"]').on('click', function (e) {
    e.preventDefault();

    switch (this.dataset.type) {
        case 'string':
            $('#addString-Modal .modal-title > span').text("????????????")
            $('#addString-Modal').modal('show');
            break;
        case 'time':
            $('#addTime-Modal .modal-title > span').text("??????????")
            $('#addTime-Modal').modal('show');
            break;
        case 'day':
            $('#addDay-Modal .modal-title > span').text("????????")
            $('#addDay-Modal').modal('show');
            break;
    }

    $('.modal').attr('target', this.dataset.target)

    return false;
})

$('[data-action="delete"]').on('click', function (e) {
    e.preventDefault();

    const card = $(this).parents('.card');
    const checked = card.find('input[type="checkbox"]:checked')
    let ids = [];

    checked.each((i, item) => {
        ids.push($(item).data('id'))
    })
    __Api.delete(`dictionaries/${this.dataset.target}`, null, ids)
        .then(() => {
            __Api.get('dictionaries')
                .then((response) => {
                    $('.card-body[id*="-list"]').each((ind, item) => item.innerHTML = '')
                    __App.options = response.data

                    toggleItems()

                    __Api.get('schedules')
                        .then((response) => {
                            __App.schedules = response.data.Map
                            $('.modal').modal('hide');
                        })
                })
        })

    return false
})

$('#system-options').on('click', function (e) {
    e.preventDefault()

    const modal = $('#System-Modal')

    __Api.get('system')
        .then((response) => {
            __App.system = response.data
        })

    modal.modal('show')

    return false
})

$('#addTask-Modal .modal-footer .btn-primary').on('click', function (e) {
    e.preventDefault()

    const modal = $(this).parents('.modal');
    const form = modal.find('form')

    __Api.post('schedules', __Api.serializeObject(form), false)
        .then(() => {
            __Api.get('dictionaries')
                .then((response) => {
                    $('.card-body[id*="-list"]').each((ind, item) => item.innerHTML = '')
                    __App.options = response.data

                    toggleItems()

                    __Api.get('schedules')
                        .then((response) => {
                            __App.schedules = response.data.Map
                            $('.modal').modal('hide');
                        })
                })
        })


    return false
})

$('#System-Modal .modal-footer .btn-primary').on('click', function (e) {
    e.preventDefault()

    const modal = $(this).parents('.modal');
    const form = modal.find('form')

    __Api.post('system', __Api.serializeObject(form), false)
        .then(() => {
            __Api.get('system')
                .then((response) => {
                    __App.system = response.data
                    $('.modal').modal('hide');
                })
        })


    return false
})

$(document).on('click', '#scheduler table td input[type="checkbox"][id^="task-"]', function (e) {
    e.preventDefault()

    const id = $(this).attr('id').replace('task-', '')

    __Api.get(`schedules/toggle/${id}`, null, false)
        .then(() => {
            __Api.get('dictionaries')
                .then((response) => {
                    $('.card-body[id*="-list"]').each((ind, item) => item.innerHTML = '')
                    __App.options = response.data

                    toggleItems()

                    __Api.get('schedules')
                        .then((response) => {
                            __App.schedules = response.data.Map
                            $('.modal').modal('hide');
                        })
                })
        })

    return false
})

$(document).on('click', '#task-add', function (e) {
    e.preventDefault()

    const modal = $('#addTask-Modal');
    modal.modal('show');
    modal.find('select').html('')

    let select = modal.find('[name="day"]');
    __App.options.Scheduler.Day.forEach((item, index) => {
        let days = item.map((day, i) => {
            if (day === "1") {
                return __App.enums.days[i]
            }
            return ''
        }).filter(d => d !== '')

        select.append(`<option value="${index + 1}">${days.join(',')}</option>`)
    })

    select = modal.find('[name="time"]');
    __App.options.Scheduler.Time.forEach((item, index) => {
        select.append(`<option value="${index + 1}">${item}</option>`)
    })

    select = modal.find('[name="sender"]');
    __App.options.Templates.Sender.forEach((item, index) => {
        select.append(`<option value="${index + 1}">${item}</option>`)
    })

    select = modal.find('[name="subject"]');
    __App.options.Templates.Subject.forEach((item, index) => {
        select.append(`<option value="${index + 1}">${item}</option>`)
    })

    select = modal.find('[name="phone"]');
    __App.options.Destinations.Tel_num.forEach((item, index) => {
        select.append(`<option value="${index + 1}">${item}</option>`)
    })

    return false
})

$(document).on('click', '[id^="delete-task-"]', function (e) {
    e.preventDefault()

    const id = $(this).attr('id').replace('delete-task-', '')

    __Api.delete(`schedules`, id, null)
        .then(() => {
            __Api.get('dictionaries')
                .then((response) => {
                    $('.card-body[id*="-list"]').each((ind, item) => item.innerHTML = '')
                    __App.options = response.data

                    toggleItems()

                    __Api.get('schedules')
                        .then((response) => {
                            __App.schedules = response.data.Map
                            $('.modal').modal('hide');
                        })
                })
        })

    return false
})

$('#addDay-Modal .modal-footer .btn-primary, ' +
    '#addTime-Modal .modal-footer .btn-primary, ' +
    '#addString-Modal .modal-footer .btn-primary').on('click', function (e) {
    e.preventDefault()

    const modal = $(this).parents('.modal');
    const form = modal.find('form')
    form.find('input').attr('name', modal.attr('target'))

    fetch('/api/dictionaries', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: __Api.prepareData(__Api.serializeObject(form))
    })
        .then(response => {
            //response.json()
            __Api.get('dictionaries')
                .then((response) => {
                    $('.card-body[id*="-list"]').each((ind, item) => item.innerHTML = '')
                    __App.options = response.data

                    toggleItems()

                    __Api.get('schedules')
                        .then((response) => {
                            __App.schedules = response.data.Map
                            $('.modal').modal('hide');
                        })
                })
        })

    return false
})