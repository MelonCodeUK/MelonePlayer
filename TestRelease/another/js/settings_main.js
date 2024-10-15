// Подключаемся к сокету

function replaceColor(element) {
    // Регулярное выражение для поиска конструкций {:color:}
    var regex = /{:(\w+):}/g;
    // Получаем HTML содержимое элемента

    // Заменяем все совпадения с регулярным выражением
    var replacedHtml = element.replace(regex, function(match, colorWord) {
        // Возвращаем замененную строку с восклицательными знаками
        return `<span style="color: ${colorWord}">`;
    });
    // Заменяем HTML содержимое элемента на измененное
    return replacedHtml;
}

async function getI(parameter) {
    try {
        // Выполняем асинхронный AJAX-запрос с помощью await
        const response = await $.ajax({
            url: `./send?command=${parameter}`,
            type: "GET",
        });

        // Попытка преобразовать ответ в JSON
        try {
            const jsonResponse = JSON.parse(response); // Преобразуем ответ в объект JSON
            return jsonResponse; // Возвращаем данные
        } catch (error) {
            return response;
        }
    } catch (error) {
        console.error("Ошибка AJAX-запроса:", error);
    }
}
async function setSettings() {
    // завантаження доступних тем
    let result = await getI("get.Data.app_settings.theme"); //дістати.інформацію.налаштування_програми.тема_за_замовчуванням
    $("#thema").append(`<option value=${result}>${result}</option>`);
    result = await getI("get.Data.themes"); //дістати.інформацію.можливі_теми
    for (key in result) {
        $("#thema").append(
            `<option value="${result[key]["name"]}">${result[key]["name"]}</option>`
        );
    }
}

$(document).ready(function() {
    $("#sendTESTMSG").on("click", function(e) {
        $("#messageInput").val('/Download.video.(https://youtu.be/HMdY9CYBrIU)')

    })
    $("#tabs a").on("click", function(e) {
        e.preventDefault();
        $(this).tab("show");
    });
});