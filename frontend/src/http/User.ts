import { $api } from "@/http/api";

export class UserService {
    static async login(nickname: string, password: string): Promise<[boolean, any]> {
        try {
            const { status, data } = await $api.post("/auth/login", {
                nickname,
                password,
            });
            console.log("231321");
            return [status === 200, data];
        } catch (error: any) {
            console.error("Login error:", error);
            return [false, error.response?.data || { message: "Unknown error" }];
        }
    }

    static async register(email: string, nickname: string, password: string): Promise<[boolean, any]> {
        try {
            const { status, data } = await $api.post("/auth/register", {
                email,
                nickname,
                password,
            });
            return [status === 200, data];
        } catch (error: any) {
            console.error("Register error:", error);
            return [false, error.response?.data || { message: "Unknown error" }];
        }
    }

    static async me(): Promise<[boolean, any]> {
        try {
            const { status, data } = await $api.put("/auth/me");
            return [status === 200, data];
        } catch (error: any) {
            console.error("Me error:", error);
            return [false, error.response?.data || { message: "Unknown error" }];
        }
    }

    static async logout(): Promise<[boolean, any]> {
        try {
            const { status, data } = await $api.post("/auth/logout");
            return [status === 200, data];
        } catch (error: any) {
            console.error("Logout error:", error);
            return [false, error.response?.data || { message: "Unknown error" }];
        }
    }

    static async refreshToken(): Promise<[boolean, any]> {
        try {
            const { status, data } = await $api.post("/auth/refresh");
            return [status === 200, data];
        } catch (error: any) {
            console.error("Refresh token error:", error);
            return [false, error.response?.data || { message: "Unknown error" }];
        }
    }
}
