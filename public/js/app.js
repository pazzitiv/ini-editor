const __App = new Proxy({
    addItem: (value, dataType, index, type = 'checkbox') => {
        `<div className="form-check">
                <input className="form-check-input" type="${type}" value="" id="${dataType}-${index}">
                    <label className="form-check-label" htmlFor="${dataType}-${index}">
                        ${value}
                    </label>
            </div>`
    },
    enums: {
        log_levels: [
            'INFO',
            'DEBUG',
        ],
    },
    options: {},
}, proxyHandler)

$(document).ready(() => {
    __Api.get('schedule')
        .then((data) => {
            console.log('DATA', data)
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

    $(this).parents('.modal')[0].attr('target', this.dataset.target)

    return false;
})

$('#addDay-Modal .modal-footer .btn-primary').on('click', function (e) {
    e.preventDefault()

    const modal = $(this).parents('.modal');
    const form = modal.children('form')
    //const sect = modal.dataset.target ?? null

    fetch('/api/sections', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new FormData(form[0])
    })

    return false
})