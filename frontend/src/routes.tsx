import { LOGIN_ROUTE, REGISTER_ROUTE, SETTINGS_ROUTE, QUEUE_ROUTE, THEMES_ROUTE, PROMPTS_ROUTE, DASHBOARD_ROUTE, ADD_ARTICLE_ROUTE} from "./utils/consts";


export const authRoutes = [
    {
        path: LOGIN_ROUTE,
        component: Login
    }, 
    {
        path: REGISTER_ROUTE,
        component: Register
    },
]
export const primaryRoutes = [
    {
        path: SETTINGS_ROUTE,
        component: Settings
    },
    {
        path: ADD_ARTICLE_ROUTE,
        component: Add_Article
    },
    {
        path: DASHBOARD_ROUTE,
        component: DashBoard
    },
    {
        path:  PROMPTS_ROUTE,
        component: Prompts
    },
    {
        path: QUEUE_ROUTE,
        component: Queue
    },
    {
        path: THEMES_ROUTE,
        component: Themes
    }
]