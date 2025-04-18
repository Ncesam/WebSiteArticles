import axios from "axios";
import * as process from "node:process";

export const $api = axios.create({
    baseURL: process.env.REACT_APP_API,
    withCredentials: true,
})