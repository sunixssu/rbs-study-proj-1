let echoSocket = null;
let jsonBody = {
    msg: "test",
}

async function getData() {
    const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));
    const socket = new WebSocket('ws://localhost:8080/telemetry')

    socket.onopen = async (event) => {
        console.log("WebSocket connected!")
    }
    socket.onmessage = (event) => {
        const data = JSON.parse(event.data)
        let ram = data.ram
        let cpu_freq = data.cpu_freq

        let cpu_freq_elem = document.getElementById("cpu_freq")
        let ram_elem = document.getElementById("ram")

        cpu_freq_elem.innerHTML = cpu_freq
        ram_elem.innerHTML = ram
    }
    socket.onerror = (error) => {
        console.error('Ошибка WebSocket:', error);
    };
}

async function sendEcho() {

    if (echoSocket && (echoSocket.readyState == WebSocket.OPEN || echoSocket.readyState == WebSocket.CONNECTING)) {
        console.log("WebSocket already open")

        let inptTextEchoBlock = document.getElementById("inptTextEcho")
        inptText = inptTextEchoBlock.value
        jsonBody.msg = inptText

        const timestamp = Date.now()
        const timeString = new Date(timestamp).toLocaleTimeString(); 
        let textResponse = jsonBody.msg + "<br>Timestamp:" + timeString

        outputTextEchoBlock.style.display = "block"
        console.log(jsonBody)
        outputTextEchoBlock.innerHTML = textResponse
        echoSocket.send(textResponse)
    } 
    else {
        echoSocket = new WebSocket('ws://localhost:8080/echo')
    }

    echoSocket.onopen = (event) => {
        console.log("WebSocket connected!")
        const timestamp = Date.now()
        const timeString = new Date(timestamp).toLocaleTimeString(); 
        let textResponse = jsonBody.msg + "<br>Timestamp:" + timeString

        outputTextEchoBlock.style.display = "block"
        console.log(jsonBody)
        outputTextEchoBlock.innerHTML = textResponse
        echoSocket.send(textResponse)
    }
}

let outputTextEchoBlock = document.getElementById("outputTextEcho")
outputTextEchoBlock.style.display = "none";
getData()