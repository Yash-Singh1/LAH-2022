import { createRoot } from "react-dom/client";
import App from "./app";
import log from "@src/utils/log";
import { ENDPOINTS, POPUP_KEYS } from "@src/constants";

const loadInterval = setInterval(tryLoad, 500)

function tryLoad() {
    const elem = document.getElementById(":3")
    if (elem == undefined) return log("Cannot find elements :rage:")

    clearInterval(loadInterval)
    
    const root = document.createElement("div");
    root.id = POPUP_KEYS.INFO;
    
    elem.appendChild(root)
    
    const email = window["GLOBALS"][10]
    chrome.storage.local.get(email, async (d) => {
        createRoot(root).render(<App />);
    })
}
tryLoad()

