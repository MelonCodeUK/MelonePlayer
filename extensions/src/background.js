


const ws = new WebSocket("ws://127.0.0.1:2308/ws");
// Создаем пункт контекстного меню
chrome.contextMenus.create({
  id: "add-music",
  title: "add music",
  contexts: ["all"],
  documentUrlPatterns: ["*://www.youtube.com/watch*"]
});

// Обрабатываем нажатие на пункт меню
chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (info.menuItemId === "add-music") {
    console.log(`Download.music.(${info.pageUrl})`)
    ws.send(`Download.music.(${info.pageUrl})`)
    ws.send(`Download.preview.(${info.pageUrl})`)


  }
});

ws.onopen = function(event) {
  console.log("Connected to WebSocket server");
};

ws.onerror = function(error) {
  console.error("ws error: ", error);
};
ws.onmessage = function(event) {
console.log(event.data)
}