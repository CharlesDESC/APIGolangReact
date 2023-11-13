import React from "react";
import ReactDOM from "react-dom/client";
import { HomePage } from "./components/homePage.jsx";
import "./styles/index.css";
import App from "./components/App.jsx";

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
