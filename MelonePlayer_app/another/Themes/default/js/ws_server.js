$(document).ready(function() {

    const ws = new WebSocket("ws://127.0.0.1:2308/ws");
    $("#info_window").hide();


    ws.onopen = function(event) {
        console.log("Connected to WebSocket server");
        $("#messages").append("Connected server!");
    };

    ws.onmessage = function(event) {
        data = event.data
        console.log("Message from server: ", data);
        $("#messages").append(data);
        if (data.includes("[download]")) {
            $("#info_window").show();
            $("#info_window").css("width", `-1%`)
            download_data = data.split(" ")
            download_data.shift();
            console.log(download_data);
            $(".progress").css("width", `${download_data[0]}`)
            if (`${download_data[0]}` == "100%" || download_data.lenght == 0) {
                $("#info_window").hide();
            }


        }


        // ["[download]", "1.0%", "24.46MiB", "273.36KiB/s", "01:30"]
    };
    $("#Menu_Button_Close").click(async function() {
        await ws.send("Player.Hide");
    })
    $("#Menu_Buttons_Settings").click(async function() {
        await ws.send("SettingsW.Show");
    })
    ws.onclose = function(event) {
        $("#messages").append("Disconnected from server");

        console.log("Disconnected from server");
    };

    ws.onerror = function(error) {
        console.log("ws error: ", error);
        $("#messages").append(error);
    };
    // Обработчик нажатия кнопки "Send Message"
    $("#sendMessage").click(async function() {
        var message = $("#messageInput").val();
        if (message[0] == "/" || message[0] == "\\") {
            await ws.send(message.slice(1));
        } else {
            // Отправляем сообщение на сервер
            var event = await getI(message);
            // Событие, когда сообщение получено от сервера
            var message = replaceColor(`<p> ${event} </p>`);
            $("#messages").append(message);
        }
    });
});