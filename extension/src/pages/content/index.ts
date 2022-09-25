import { unwatchFile } from "fs";

console.log("Loaded extension content entrypoint")

var s = document.createElement("script");
s.setAttribute('type', 'text/javascript');
s.setAttribute('src', chrome.runtime.getURL('/inject.js'));
document.getElementsByTagName("body")[0].appendChild(s)

// chrome.storage.local.clear()
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
const API = "http://localhost:3500"

function onHashChange(){
    const oldInfo = document.getElementById("lah-info-popup")
    if (/#(inbox|spam|sent|trash|starred)\/.+$/.test(location.hash) && !document.getElementById("lah-login-popup")) {
        if(oldInfo) {
            console.log("Showing old popup")
            oldInfo.style.visibility = "visible"
        } else {
            import("./components/info-popup")
        }

        const email = window['GLOBALS'][10];
        chrome.storage.local.get(email, async (d) => {
            if(!d[email]) return
            if(d[email].optOut == true) return console.log("Email opt-out, quitting")
            if(!d[email].accessToken) return

            let emailElem = undefined
            while(emailElem == undefined) {
                emailElem = document.querySelector('[data-message-id]')
                await sleep(500)
            }

            const emailId = emailElem.getAttribute('data-legacy-message-id')
            let res = await fetch(`${API}/predict/email/${emailId}`, {
                mode: "cors",
                headers: {
                    "Authorization": d[email].accessToken
                }
            })
    
            if (res.status !== 200) {
                console.log("Something bad happened!")
            }
    
            const data = await res.json()
            window["setEmailInfo"](data)
        })    
    } else if (oldInfo) {
        console.log("Hiding old popup")
        oldInfo.style.visibility = "hidden"
    }
}

window.addEventListener('message', (e) => {
    if(!e.data['GLOBALS']) return
    window['GLOBALS'] = e.data['GLOBALS']

    const email = window['GLOBALS'][10]
    console.log(`Received globals - email is ${email}`)

    window.addEventListener('hashchange', onHashChange)
    onHashChange()

    chrome.storage.local.get(email, (d) => {
        if(d[email] && d[email].optOut == true) return console.log("Email opt-out, quitting")
        if(!d[email] || !d[email].accessToken) return import("./components/login-popup")
    })
});

