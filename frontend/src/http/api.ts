import axios from "axios";
import * as process from "process";

export const $api = axios.create({
    baseURL: process.env.REACT_APP_API,
    withCredentials: true,
})