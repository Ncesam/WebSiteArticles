


export interface Articles {
    title: string,
    subtitle: string,
    isPublished: boolean,
    blocks: Block[],
    category: string
}

export interface Block {
    text: string,
    type: string,
    number: number
}