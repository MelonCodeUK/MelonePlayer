let progressBar;
let music;
let preview;
let title;
let nextButton;
let previousButton;
let playPause;
let playList;
let settings;
let musics = [];
let musicIndex = 0;

function init() {
    progressBar = $("#TwoProgressBar");
    music = $("#music");
    preview = $(".PreviewImg");
    title = $("#title");
    nextButton = $("#next");
    previousButton = $("#previous");
    playPause = $(".play_pause");
    playList = $("#Menu_Buttons_PlayList");
    settings = $("#Menu_Buttons_Settings");
}



async function getI(parameter) {
    try {
        // Выполняем асинхронный AJAX-запрос с помощью await
        const response = await $.ajax({
            url: `./send?command=${parameter}`,
            type: "GET",
        });

        // // Преобразуем ответ в объект JSON
        // const jsonResponse = JSON.parse(response); // response уже строка, так что просто парсим

        return response; // Возвращаем данные
    } catch (error) {
        console.error("Ошибка при выполнении запроса или парсинга JSON:", error);
        throw error; // Пробрасываем ошибку вверх для обработки
    }
}

function musicLegular(wo) {
    if (wo == "next") {
        musicIndex++;
        if (musicIndex >= musics[0].length) {
            musicIndex = 0;
        }
        $("#music").attr("src", musics[0][musicIndex]);
        $(".PreviewImg").attr("src", musics[1][musicIndex]);
        $(".PreviewImg").attr("alt", musics[1][musicIndex]);
        $("#music")[0].play();
        playPause.attr("id", "pause");
        playPause.attr("src", "./icons/svg/pausa.svg");
    } else if (wo == "previous") {
        musicIndex--;
        if (musicIndex < 0) {
            musicIndex = musics[0].length - 1;
        }
        console.log(musicIndex);
        $("#music").attr("src", musics[0][musicIndex]);
        $(".PreviewImg").attr("src", musics[1][musicIndex]);
        $(".PreviewImg").attr("alt", musics[1][musicIndex]);
        $("#music")[0].play();
        playPause.attr("id", "pause");
        playPause.attr("src", "./icons/svg/pausa.svg");
    } else {
        console.error(`Invalid wo==${wo}`);
    }
}

function getMusics() {
    return new Promise(async(resolve, reject) => {
        var result = await getI("get.Default.PlayList");
        console.log(result);
        console.log("Отправлено сообщение: " + "get.Default.PlayList");

        var i1 = [
            [],
            []
        ];
        var r1 = 0;
        for (i in result) {
            result[i].forEach((item, index) => {
                console.log(i1[r1]);
                i1[r1].push(item);
            });
            r1++;
            resolve(i1);
        }
    });
}

function play__pause() {
    let getId = playPause.attr("id");
    if (getId == "pause") {
        $("#music")[0].pause();
        playPause.attr("id", "play");
        playPause.attr("src", "./icons/svg/play.svg");
    } else if (getId == "play") {
        $("#music")[0].play();
        playPause.attr("id", "pause");
        playPause.attr("src", "./icons/svg/pausa.svg");
    } else {
        console.error(`Invalid argument==${getId}`);
    }
}

function progressBarUpdate(e) {
    let currentTime = e.currentTime;
    let duration = e.duration;
    let procent = (currentTime / duration) * 100;
    progressBar.css("width", `${procent}%`);
}

function setProgress(thissss, e) {
    let width = thissss.clientWidth;
    let clickX = e.offsetX;
    let duration = $("#music")[0].duration;
    $("#music")[0].currentTime = (clickX / width) * duration;
}

$(document).ready(async function() {
    init();
    try {
        musics = await getMusics(); // Дожидаемся разрешения обещания
        console.log(musics);
        music.attr("src", musics[0][musicIndex]);
        playPause.on("click", function() {
            play__pause();
        });
        nextButton.on("click", function() {
            musicLegular("next");
        });
        musicLegular("next");
        previousButton.on("click", function() {
            musicLegular("previous");
        });
        musicLegular("previous");

        // Используем jQuery для назначения обработчика события timeupdate
        $("#music").on("timeupdate", function() {
            progressBarUpdate(this);
        });
        $("#ProgressBar").on("click", function(e) {
            setProgress(this, e);
        });
        music.on("ended", function() {
            musicLegular("next");
        });
    } catch (error) {
        console.error("Ошибка получения музыки:", error);
        console.error("Що за хуйня міститься у цій змінній:" + musics);
    }
});