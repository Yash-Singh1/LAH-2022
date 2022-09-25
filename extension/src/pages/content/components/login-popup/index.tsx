import { createRoot } from "react-dom/client";
import App from "./app";
import log from "@src/utils/log";
import { POPUP_KEYS } from "@src/constants";

const root = document.createElement("div");
root.id = POPUP_KEYS.LOGIN;

const loadInterval = setInterval(tryLoad, 500)

function tryLoad() {
    const elem = document.getElementById(":3")
    if (elem == undefined) return log("Cannot find element :rage:")
    clearInterval(loadInterval)
    elem.appendChild(root)

    createRoot(root).render(<App />);
}
tryLoad()