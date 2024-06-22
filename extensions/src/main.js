$("yt-button-view-model")
  .find("button")
  .each(function () {
    // Для каждой кнопки добавляем обработчик клика
    $(this).on("click", function () {
      // Логика, которая будет выполняться при клике на кнопку
      console.error("Кнопка была нажата");
    });

    // Выводим внешний HTML кнопки, если это необходимо
    console.error($(this).prop("outerHTML"));
  });
