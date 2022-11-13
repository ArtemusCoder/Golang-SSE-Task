var source = new EventSource("/listen");
source.onmessage = function (event) {
    console.log(event.data);
    var data = JSON.parse(event.data)
    document.getElementById("message").innerText = data["word"];
};