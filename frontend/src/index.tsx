import React, {createContext} from "react";
import ReactDOM from "react-dom/client";
import {Provider} from "mobx-react";
import App from "@/App";

import rootStore from "./stores/rootStore";
const Context = createContext<rootStore | undefined>();
const rootElement = document.getElementById("root");
const root = ReactDOM.createRoot(rootElement);

root.render(
  <React.StrictMode>
    <Provider {...rootStore}>
      <App />
    </Provider>
  </React.StrictMode>
);

