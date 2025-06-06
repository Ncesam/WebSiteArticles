
export interface BotConfig {
    Id: string,
    name: string,
    userId: number,
    prompt: string,
    delay: number,
}

export interface FormBotConfig {
    name: string,
    prompt: string,
    delay: number
    email: string,
    password: string,
}

export interface StartConfig {
    configId: string,
    userId: number,
    prompt: string,
    data: File,
}

