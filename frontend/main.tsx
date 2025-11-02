import React from "react";
import {createRoot} from "react-dom/client";

import * as Sentry from "@sentry/react";

import App from "./App";
import "./css/index.css";

Sentry.init({
    dsn: "https://863fd3bc9b0e4524b442365aa7b55f38@bugs.dimhost.ru/2",
    integrations: [],
    tracesSampleRate: 0.1,
});


const container = document.getElementById("root") as HTMLElement;
const root = createRoot(container, {
    // Callback called when an error is thrown and not caught by an ErrorBoundary.
    onUncaughtError: Sentry.reactErrorHandler((error, errorInfo) => {
        console.warn('Uncaught error', error, errorInfo.componentStack);
    }),
    // Callback called when React catches an error in an ErrorBoundary.
    onCaughtError: Sentry.reactErrorHandler(),
    // Callback called when React automatically recovers from errors.
    onRecoverableError: Sentry.reactErrorHandler(),
});
root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
);
