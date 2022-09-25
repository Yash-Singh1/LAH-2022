export default function log(...message: any[]) {
    console.log("%c[EmailSEC]", "color: blue", ...message)
}