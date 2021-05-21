const __App = new Proxy({
    addItem: (value, dataType, index, type = 'checkbox') => {
        $(`#${dataType}-list`).append(`<div className="form-check">
                <input className="form-check-input" name="${dataType}[]" type="${type}" value="" id="${dataType}-${index}">
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
                        cell.innerHTML = __App.options.Scheduler.Time[parseInt(item) - 1]
                    }
                    break
                case 2:
                    if (__App.has(__App.options, 'Templates')) {
                        cell.innerHTML = __App.options.Templates.Sender[parseInt(item) - 1]
                    }
                    break
                case 3:
                    if (__App.has(__App.options, 'Templates')) {
                        cell.innerHTML = __App.options.Templates.Subject[parseInt(item) - 1]
                    }
                    break
                case 4:
                    if (__App.has(__App.options, 'Destinations')) {
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
        tb.append(row)
    },
    has: (object, prop) => {
        const hOP = Object.prototype.hasOwnProperty
        return hOP.call(object, prop)
    },
    enums: {
        days: ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"],
        log_levels: [
            'INFO',
            'DEBUG',
        ],
    },
    options: {},
    schedules: {},
}, proxyHandler)

$(document).ready(() => {
    __Api.get('dictionaries')
        .then((response) => {
            __App.options = response.data

            __Api.get('schedules')
                .then((response) => {
                    __App.schedules = response.data.Map
                })
        })
})

$('[data-action="add"]').on('click', function (e) {
    e.preventDefault();

    switch (this.dataset.type) {
        case 'string':
            $('#addString-Modal .modal-title > span').text("Строка")
            $('#addString-Modal').modal('show');
            break;
        case 'time':
            $('#addTime-Modal .modal-title > span').text("Время")
            $('#addTime-Modal').modal('show');
            break;
        case 'day':
            $('#addDay-Modal .modal-title > span').text("День")
            $('#addDay-Modal').modal('show');
            break;
    }

    $('.modal').attr('target', this.dataset.target)

    return false;
})

$('#addDay-Modal .modal-footer .btn-primary').on('click', function (e) {
    e.preventDefault()

    const modal = $(this).parents('.modal');
    const form = modal.find('form')

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

                    __Api.get('schedules')
                        .then((response) => {
                            __App.schedules = response.data.Map
                            $('.modal').modal('hide');
                        })
                })
        })

    return false
})