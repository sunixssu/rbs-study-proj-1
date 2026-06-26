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
    /*
    while(true) {
        await sleep(1000)
        let response = await fetch("http://localhost:8080/telemetry")
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`)
        }

        const data = await response.json();

        let ram = data.ram
        let cpu_freq = data.cpu_freq

        let cpu_freq_elem = document.getElementById("cpu_freq")
        let ram_elem = document.getElementById("ram")

        cpu_freq_elem.innerHTML = cpu_freq
        ram_elem.innerHTML = ram
    }
    */
}

async function sendEcho() {
    let jsonBody = {
        msg: "test",
    }
    let inptTextEchoBlock = document.getElementById("inptTextEcho")
    inptText = inptTextEchoBlock.value
    jsonBody.msg = inptText

    const socket = new WebSocket('ws://localhost:8080/echo')

    socket.onopen = (event) => {
        console.log("WebSocket connected!")
        const timestamp = Date.now()
        const timeString = new Date(timestamp).toLocaleTimeString(); 
        let textResponse = jsonBody.msg + "<br>Timestamp:" + timeString

        outputTextEchoBlock.style.display = "block"
        console.log(jsonBody)
        outputTextEchoBlock.innerHTML = textResponse
        socket.send(textResponse)
    }

    //let response = await fetch("http://localhost:8080/echo", {
    //    method: 'GET',
    //    headers: {
    //        'Content-Type': 'text/plain',
    //        'Accept': 'text/plain'
    //    },
    //    body: jsonBody.msg
    //})
    //let data = await response.json()
}

let outputTextEchoBlock = document.getElementById("outputTextEcho")
outputTextEchoBlock.style.display = "none";
getData()