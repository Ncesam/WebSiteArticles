import React, {createContext} from "react";
import ReactDOM from "react-dom/client";
import App from "@/App";
import RootStore from "@/stores/rootStore";

const store = new RootStore()

export const Context = createContext<RootStore>(store);

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);


root.render(
  <Context.Provider value={store}>
      <App />
  </Context.Provider>
);

