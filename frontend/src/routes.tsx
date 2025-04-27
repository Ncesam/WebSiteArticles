import {lazy} from "react";
import {ADD_ARTICLE_ROUTE, DASHBOARD_ROUTE, LOGIN_ROUTE, REGISTER_ROUTE, SETTINGS_ROUTE} from "./utils/consts";
import DashBoard from "./pages/DashBoard/DashBoard";
import AddArticle from "./pages/AddArticle/AddArticle";

const Registration = lazy(() => import("@/pages/Registration/Registration"));
const Login = lazy(() => import("@/pages/Login/Login"))

export const authRoutes = [
    {
        path: LOGIN_ROUTE,
        component: Login
    },
    {
        path: REGISTER_ROUTE,
        component: Registration
    },
]
export const primaryRoutes = [
    // {
    //     path: SETTINGS_ROUTE,
    //     component: Settings
    // },
    {
        path: ADD_ARTICLE_ROUTE,
        component: AddArticle
    },
    {
        path: DASHBOARD_ROUTE,
        component: DashBoard
    },
    // {
    //     path:  PROMPTS_ROUTE,
    //     component: Prompts
    // },
    // {
    //     path: QUEUE_ROUTE,
    //     component: Queue
    // },
    // {
    //     path: THEMES_ROUTE,
    //     component: Themes
    // }
]