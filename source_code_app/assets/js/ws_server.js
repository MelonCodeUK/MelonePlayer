$(document).ready(function () {
  const ws = new WebSocket("ws://127.0.0.1:9001/ws");

  ws.onopen = function (event) {
    console.log("Connected to WebSocket server");
    $("#messages").append("Connected server!");
  };

  ws.onmessage = function (event) {
    console.log("Message from server: ", event.data);
    var message = replaceColor(`<p> ${event.data} </p>`);
    $("#messages").append(message);
    console.log(event.data);
  };

  ws.onclose = function (event) {
    $("#messages").append("Disconnected from server");

    console.log("Disconnected from server");
  };

  ws.onerror = function (error) {
    console.log("ws error: ", error);
    $("#messages").append(error);
  };
  // Обработчик нажатия кнопки "Send Message"
  $("#sendMessage").click(async function () {
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
