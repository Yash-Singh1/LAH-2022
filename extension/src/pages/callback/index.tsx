import { createRoot } from "react-dom/client";
import App from "./app";

const root = document.getElementById("app-container")!;
createRoot(root).render(<App />);
