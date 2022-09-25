export const EXTENSION_NAME = "[ext name]"

export const API = "http://localhost:3500"

export const ENDPOINTS = Object.freeze({
    AUTH: (email: string) => `${API}/auth/login?email=${encodeURIComponent(email)}&callback=${encodeURIComponent(chrome.runtime.getURL('/src/pages/callback/index.html'))}`,
    PREDICT_EMAIL: (emailId: string) => `${API}/predict/email/${emailId}`
})

export const POPUP_KEYS = Object.freeze({
    LOGIN: "lah-login-popup",
    INFO: "lah-info-popup"
})